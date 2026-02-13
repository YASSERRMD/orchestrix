package main

import (
	"flowforge/internal/config"
	"flowforge/internal/database"
	"flowforge/internal/handlers"
	"flowforge/internal/middleware"
	"flowforge/internal/services"
	"flowforge/internal/websocket"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.Load()

	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	runMigrations()

	authService := services.NewAuthService(cfg)
	userService := services.NewUserService()
	templateService := services.NewTemplateService()
	workflowService := services.NewWorkflowService()

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	templateHandler := handlers.NewTemplateHandler(templateService)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)

	hub := websocket.NewHub()
	go hub.Run()

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		api.GET("/auth/me", authHandler.Me)

		api.GET("/users", middleware.RoleMiddleware("admin", "manager"), userHandler.List)
		api.POST("/users", middleware.RoleMiddleware("admin"), userHandler.Create)
		api.GET("/users/:id", userHandler.Get)
		api.PUT("/users/:id", middleware.RoleMiddleware("admin"), userHandler.Update)
		api.DELETE("/users/:id", middleware.RoleMiddleware("admin"), userHandler.Delete)

		api.GET("/templates", templateHandler.List)
		api.POST("/templates", middleware.RoleMiddleware("admin", "manager"), templateHandler.Create)
		api.GET("/templates/:id", templateHandler.Get)
		api.GET("/templates/:id/versions", templateHandler.GetVersions)
		api.PUT("/templates/:id", middleware.RoleMiddleware("admin", "manager"), templateHandler.Update)

		api.GET("/workflows", workflowHandler.List)
		api.POST("/workflows", middleware.RoleMiddleware("admin", "manager", "operator"), workflowHandler.Create)
		api.GET("/workflows/:id", workflowHandler.Get)
		api.POST("/workflows/:id/advance", workflowHandler.Advance)
		api.POST("/workflows/:id/reject", workflowHandler.Reject)
		api.POST("/workflows/:id/rollback", middleware.RoleMiddleware("admin", "manager"), workflowHandler.Rollback)
		api.POST("/workflows/:id/cancel", workflowHandler.Cancel)
		api.GET("/workflows/:id/audit", workflowHandler.GetAudit)

		api.GET("/ws", func(c *gin.Context) {
			websocket.ServeWs(hub, c.Writer, c.Request)
		})
	}

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}

func runMigrations() {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			name TEXT NOT NULL,
			role TEXT NOT NULL CHECK(role IN ('admin', 'manager', 'operator', 'viewer')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS workflow_templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			version INTEGER NOT NULL DEFAULT 1,
			is_active INTEGER DEFAULT 1,
			created_by INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS template_stages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			template_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			order_index INTEGER NOT NULL,
			required_role TEXT NOT NULL CHECK(required_role IN ('admin', 'manager', 'operator', 'viewer')),
			approval_type TEXT NOT NULL CHECK(approval_type IN ('manual', 'auto')),
			timeout_minutes INTEGER DEFAULT 0,
			escalation_role TEXT CHECK(escalation_role IN ('admin', 'manager', 'operator', 'viewer')),
			FOREIGN KEY (template_id) REFERENCES workflow_templates(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS workflow_instances (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			template_id INTEGER NOT NULL,
			template_version INTEGER NOT NULL,
			title TEXT NOT NULL,
			current_stage_id INTEGER,
			status TEXT NOT NULL CHECK(status IN ('pending', 'approved', 'rejected', 'escalated', 'canceled', 'completed')),
			created_by INTEGER NOT NULL,
			assigned_to INTEGER,
			escalated_to INTEGER,
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			completed_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (template_id) REFERENCES workflow_templates(id) ON DELETE CASCADE,
			FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (current_stage_id) REFERENCES template_stages(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			workflow_id INTEGER NOT NULL,
			actor_id INTEGER NOT NULL,
			actor_role TEXT NOT NULL,
			from_stage_id INTEGER,
			to_stage_id INTEGER,
			action TEXT NOT NULL,
			reason TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (workflow_id) REFERENCES workflow_instances(id) ON DELETE CASCADE,
			FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS locks (
			key TEXT PRIMARY KEY,
			owner TEXT NOT NULL,
			expires_at DATETIME NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_org ON users(organization_id)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_templates_org ON workflow_templates(organization_id)`,
		`CREATE INDEX IF NOT EXISTS idx_instances_org ON workflow_instances(organization_id)`,
		`CREATE INDEX IF NOT EXISTS idx_instances_status ON workflow_instances(status)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_workflow ON audit_logs(workflow_id)`,
	}

	for _, m := range migrations {
		if _, err := database.DB.Exec(m); err != nil {
			log.Printf("Migration error: %v\nSQL: %s", err, m)
		}
	}
}
