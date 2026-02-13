package models

import (
	"database/sql"
	"time"
)

type Organization struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	Slug      string       `json:"slug"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type User struct {
	ID             int64        `json:"id"`
	OrganizationID int64        `json:"organization_id"`
	Email          string       `json:"email"`
	PasswordHash   string       `json:"-"`
	FullName       string       `json:"full_name"`
	Role           string       `json:"role"`
	IsActive       bool         `json:"is_active"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at"`
}

type WorkflowTemplate struct {
	ID             int64        `json:"id"`
	OrganizationID int64        `json:"organization_id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	IsActive       bool         `json:"is_active"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at"`
}

type TemplateStage struct {
	ID          int64     `json:"id"`
	TemplateID  int64     `json:"template_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OrderIndex  int       `json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WorkflowInstance struct {
	ID             int64         `json:"id"`
	TemplateID     int64         `json:"template_id"`
	OrganizationID int64         `json:"organization_id"`
	Name           string        `json:"name"`
	Status         string        `json:"status"`
	CurrentStageID sql.NullInt64 `json:"current_stage_id"`
	CreatedBy      int64         `json:"created_by"`
	StartedAt      time.Time     `json:"started_at"`
	CompletedAt    sql.NullTime  `json:"completed_at"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	DeletedAt      sql.NullTime  `json:"deleted_at"`
}

type AuditLog struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organization_id"`
	UserID         int64     `json:"user_id"`
	Action         string    `json:"action"`
	EntityType     string    `json:"entity_type"`
	EntityID       int64     `json:"entity_id"`
	Details        string    `json:"details"`
	CreatedAt      time.Time `json:"created_at"`
}

type Lock struct {
	ID         int64     `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   int64     `json:"entity_id"`
	UserID     int64     `json:"user_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}
