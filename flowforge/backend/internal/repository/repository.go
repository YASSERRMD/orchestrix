package repository

import (
	"flowforge/internal/database"
	"flowforge/internal/models"
	"time"
)

type OrganizationRepository struct{}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{}
}

func (r *OrganizationRepository) Create(name string) (*models.Organization, error) {
	result, err := database.DB.Exec(
		"INSERT INTO organizations (name) VALUES (?)",
		name,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetByID(id)
}

func (r *OrganizationRepository) GetByID(id int64) (*models.Organization, error) {
	org := &models.Organization{}
	err := database.DB.QueryRow(
		"SELECT id, name, created_at, updated_at FROM organizations WHERE id = ?",
		id,
	).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return org, nil
}

func (r *OrganizationRepository) GetAll() ([]models.Organization, error) {
	rows, err := database.DB.Query("SELECT id, name, created_at, updated_at FROM organizations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var org models.Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(orgID int64, email, passwordHash, name, role string) (*models.User, error) {
	result, err := database.DB.Exec(
		"INSERT INTO users (organization_id, email, password_hash, name, role) VALUES (?, ?, ?, ?, ?)",
		orgID, email, passwordHash, name, role,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetByID(id)
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(
		"SELECT id, organization_id, email, password_hash, name, role, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(
		"SELECT id, organization_id, email, password_hash, name, role, created_at, updated_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByOrgID(orgID int64) ([]models.User, error) {
	rows, err := database.DB.Query(
		"SELECT id, organization_id, email, password_hash, name, role, created_at, updated_at FROM users WHERE organization_id = ?",
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := database.DB.Exec(
		"UPDATE users SET name = ?, role = ?, updated_at = ? WHERE id = ?",
		user.Name, user.Role, time.Now(), user.ID,
	)
	return err
}

func (r *UserRepository) Delete(id int64) error {
	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

type TemplateRepository struct{}

func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{}
}

func (r *TemplateRepository) Create(template *models.WorkflowTemplate) (*models.WorkflowTemplate, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO workflow_templates (organization_id, name, description, version, is_active, created_by) VALUES (?, ?, ?, ?, ?, ?)",
		template.OrganizationID, template.Name, template.Description, template.Version, template.IsActive, template.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	template.ID = id

	for i := range template.Stages {
		stage := &template.Stages[i]
		stageResult, err := tx.Exec(
			"INSERT INTO template_stages (template_id, name, order_index, required_role, approval_type, timeout_minutes, escalation_role) VALUES (?, ?, ?, ?, ?, ?, ?)",
			id, stage.Name, stage.OrderIndex, stage.RequiredRole, stage.ApprovalType, stage.TimeoutMinutes, stage.EscalationRole,
		)
		if err != nil {
			return nil, err
		}
		stageID, _ := stageResult.LastInsertId()
		stage.ID = stageID
		stage.TemplateID = id
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *TemplateRepository) GetByID(id int64) (*models.WorkflowTemplate, error) {
	template := &models.WorkflowTemplate{}
	err := database.DB.QueryRow(
		"SELECT id, organization_id, name, description, version, is_active, created_by, created_at, updated_at FROM workflow_templates WHERE id = ?",
		id,
	).Scan(&template.ID, &template.OrganizationID, &template.Name, &template.Description, &template.Version, &template.IsActive, &template.CreatedBy, &template.CreatedAt, &template.UpdatedAt)

	if err != nil {
		return nil, err
	}

	stages, err := r.GetStagesByTemplateID(id)
	if err != nil {
		return nil, err
	}
	template.Stages = stages

	return template, nil
}

func (r *TemplateRepository) GetByOrgID(orgID int64) ([]models.WorkflowTemplate, error) {
	rows, err := database.DB.Query(
		"SELECT id, organization_id, name, description, version, is_active, created_by, created_at, updated_at FROM workflow_templates WHERE organization_id = ? ORDER BY created_at DESC",
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []models.WorkflowTemplate
	for rows.Next() {
		var template models.WorkflowTemplate
		if err := rows.Scan(&template.ID, &template.OrganizationID, &template.Name, &template.Description, &template.Version, &template.IsActive, &template.CreatedBy, &template.CreatedAt, &template.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func (r *TemplateRepository) GetActiveByOrgID(orgID int64) ([]models.WorkflowTemplate, error) {
	rows, err := database.DB.Query(
		"SELECT id, organization_id, name, description, version, is_active, created_by, created_at, updated_at FROM workflow_templates WHERE organization_id = ? AND is_active = 1 ORDER BY created_at DESC",
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []models.WorkflowTemplate
	for rows.Next() {
		var template models.WorkflowTemplate
		if err := rows.Scan(&template.ID, &template.OrganizationID, &template.Name, &template.Description, &template.Version, &template.IsActive, &template.CreatedBy, &template.CreatedAt, &template.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func (r *TemplateRepository) GetStagesByTemplateID(templateID int64) ([]models.TemplateStage, error) {
	rows, err := database.DB.Query(
		"SELECT id, template_id, name, order_index, required_role, approval_type, timeout_minutes, escalation_role FROM template_stages WHERE template_id = ? ORDER BY order_index",
		templateID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []models.TemplateStage
	for rows.Next() {
		var stage models.TemplateStage
		if err := rows.Scan(&stage.ID, &stage.TemplateID, &stage.Name, &stage.OrderIndex, &stage.RequiredRole, &stage.ApprovalType, &stage.TimeoutMinutes, &stage.EscalationRole); err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}
	return stages, nil
}

func (r *TemplateRepository) Update(template *models.WorkflowTemplate) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE workflow_templates SET is_active = 0 WHERE id = ?",
		template.ID,
	)
	if err != nil {
		return err
	}

	var newVersion int
	err = tx.QueryRow(
		"SELECT COALESCE(MAX(version), 0) + 1 FROM workflow_templates WHERE organization_id = ? AND name = ?",
		template.OrganizationID, template.Name,
	).Scan(&newVersion)
	if err != nil {
		return err
	}

	result, err := tx.Exec(
		"INSERT INTO workflow_templates (organization_id, name, description, version, is_active, created_by) VALUES (?, ?, ?, ?, ?, ?)",
		template.OrganizationID, template.Name, template.Description, newVersion, true, template.CreatedBy,
	)
	if err != nil {
		return err
	}

	newID, _ := result.LastInsertId()

	for i := range template.Stages {
		stage := &template.Stages[i]
		_, err = tx.Exec(
			"INSERT INTO template_stages (template_id, name, order_index, required_role, approval_type, timeout_minutes, escalation_role) VALUES (?, ?, ?, ?, ?, ?, ?)",
			newID, stage.Name, stage.OrderIndex, stage.RequiredRole, stage.ApprovalType, stage.TimeoutMinutes, stage.EscalationRole,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

type WorkflowRepository struct{}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{}
}

func (r *WorkflowRepository) Create(instance *models.WorkflowInstance) (*models.WorkflowInstance, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO workflow_instances (organization_id, template_id, template_version, title, current_stage_id, status, created_by, assigned_to) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		instance.OrganizationID, instance.TemplateID, instance.TemplateVersion, instance.Title, instance.CurrentStageID, instance.Status, instance.CreatedBy, instance.AssignedTo,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	instance.ID = id

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *WorkflowRepository) GetByID(id int64) (*models.WorkflowInstance, error) {
	instance := &models.WorkflowInstance{}
	err := database.DB.QueryRow(
		`SELECT id, organization_id, template_id, template_version, title, current_stage_id, status, 
		 created_by, assigned_to, escalated_to, started_at, completed_at, created_at, updated_at 
		 FROM workflow_instances WHERE id = ?`,
		id,
	).Scan(&instance.ID, &instance.OrganizationID, &instance.TemplateID, &instance.TemplateVersion,
		&instance.Title, &instance.CurrentStageID, &instance.Status, &instance.CreatedBy,
		&instance.AssignedTo, &instance.EscalatedTo, &instance.StartedAt, &instance.CompletedAt,
		&instance.CreatedAt, &instance.UpdatedAt)

	if err != nil {
		return nil, err
	}

	templateRepo := NewTemplateRepository()
	template, err := templateRepo.GetByID(instance.TemplateID)
	if err != nil {
		return nil, err
	}
	instance.Template = template

	if instance.CurrentStageID != nil {
		for _, stage := range template.Stages {
			if stage.ID == *instance.CurrentStageID {
				instance.CurrentStage = &stage
				break
			}
		}
	}

	userRepo := NewUserRepository()
	if instance.CreatedBy > 0 {
		creator, _ := userRepo.GetByID(instance.CreatedBy)
		instance.CreatedByUser = creator
	}
	if instance.AssignedTo != nil && *instance.AssignedTo > 0 {
		assignee, _ := userRepo.GetByID(*instance.AssignedTo)
		instance.AssignedToUser = assignee
	}

	return instance, nil
}

func (r *WorkflowRepository) GetByOrgID(orgID int64, status string, assignedTo *int64) ([]models.WorkflowInstance, error) {
	query := `SELECT id, organization_id, template_id, template_version, title, current_stage_id, status, 
			  created_by, assigned_to, escalated_to, started_at, completed_at, created_at, updated_at 
			  FROM workflow_instances WHERE organization_id = ?`
	args := []interface{}{orgID}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if assignedTo != nil {
		query += " AND (assigned_to = ? OR escalated_to = ?)"
		args = append(args, *assignedTo, *assignedTo)
	}

	query += " ORDER BY created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []models.WorkflowInstance
	for rows.Next() {
		var instance models.WorkflowInstance
		if err := rows.Scan(&instance.ID, &instance.OrganizationID, &instance.TemplateID, &instance.TemplateVersion,
			&instance.Title, &instance.CurrentStageID, &instance.Status, &instance.CreatedBy,
			&instance.AssignedTo, &instance.EscalatedTo, &instance.StartedAt, &instance.CompletedAt,
			&instance.CreatedAt, &instance.UpdatedAt); err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (r *WorkflowRepository) Update(instance *models.WorkflowInstance) error {
	_, err := database.DB.Exec(
		`UPDATE workflow_instances SET current_stage_id = ?, status = ?, assigned_to = ?, escalated_to = ?, 
		 completed_at = ?, updated_at = ? WHERE id = ?`,
		instance.CurrentStageID, instance.Status, instance.AssignedTo, instance.EscalatedTo,
		instance.CompletedAt, time.Now(), instance.ID,
	)
	return err
}

func (r *WorkflowRepository) GetTimedOutWorkflows(threshold time.Time) ([]models.WorkflowInstance, error) {
	query := `
		SELECT wi.id, wi.organization_id, wi.template_id, wi.template_version, wi.title, wi.current_stage_id, wi.status,
			   wi.created_by, wi.assigned_to, wi.escalated_to, wi.started_at, wi.completed_at, wi.created_at, wi.updated_at
		FROM workflow_instances wi
		JOIN template_stages ts ON wi.current_stage_id = ts.id
		WHERE wi.status IN ('pending', 'escalated')
		  AND ts.timeout_minutes > 0
		  AND datetime(wi.updated_at, '+' || ts.timeout_minutes || ' minutes') <= datetime(?)
		ORDER BY wi.updated_at ASC
	`

	rows, err := database.DB.Query(query, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []models.WorkflowInstance
	for rows.Next() {
		var instance models.WorkflowInstance
		if err := rows.Scan(&instance.ID, &instance.OrganizationID, &instance.TemplateID, &instance.TemplateVersion,
			&instance.Title, &instance.CurrentStageID, &instance.Status, &instance.CreatedBy,
			&instance.AssignedTo, &instance.EscalatedTo, &instance.StartedAt, &instance.CompletedAt,
			&instance.CreatedAt, &instance.UpdatedAt); err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

type AuditRepository struct{}

func NewAuditRepository() *AuditRepository {
	return &AuditRepository{}
}

func (r *AuditRepository) Create(log *models.AuditLog) (*models.AuditLog, error) {
	result, err := database.DB.Exec(
		"INSERT INTO audit_logs (workflow_id, actor_id, actor_role, from_stage_id, to_stage_id, action, reason) VALUES (?, ?, ?, ?, ?, ?, ?)",
		log.WorkflowID, log.ActorID, log.ActorRole, log.FromStageID, log.ToStageID, log.Action, log.Reason,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	log.ID = id
	log.Timestamp = time.Now()
	return log, nil
}

func (r *AuditRepository) GetByWorkflowID(workflowID int64) ([]models.AuditLog, error) {
	rows, err := database.DB.Query(
		`SELECT al.id, al.workflow_id, al.actor_id, al.actor_role, al.from_stage_id, al.to_stage_id, 
		 al.action, al.reason, al.timestamp
		 FROM audit_logs al WHERE al.workflow_id = ? ORDER BY al.timestamp ASC`,
		workflowID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.WorkflowID, &log.ActorID, &log.ActorRole, &log.FromStageID,
			&log.ToStageID, &log.Action, &log.Reason, &log.Timestamp); err != nil {
			return nil, err
		}
		userRepo := NewUserRepository()
		actor, _ := userRepo.GetByID(log.ActorID)
		log.Actor = actor
		logs = append(logs, log)
	}
	return logs, nil
}

type LockRepository struct{}

func NewLockRepository() *LockRepository {
	return &LockRepository{}
}

func (r *LockRepository) Acquire(key, owner string, expiry time.Time) (bool, error) {
	_, err := database.DB.Exec(
		"INSERT OR IGNORE INTO locks (key, owner, expires_at) VALUES (?, ?, ?)",
		key, owner, expiry,
	)
	if err != nil {
		return false, err
	}

	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM locks WHERE key = ? AND owner = ? AND expires_at > datetime('now')", key, owner).Scan(&count)
	return count > 0, err
}

func (r *LockRepository) Release(key, owner string) error {
	_, err := database.DB.Exec("DELETE FROM locks WHERE key = ? AND owner = ?", key, owner)
	return err
}

func (r *LockRepository) Cleanup() error {
	_, err := database.DB.Exec("DELETE FROM locks WHERE expires_at <= datetime('now')")
	return err
}
