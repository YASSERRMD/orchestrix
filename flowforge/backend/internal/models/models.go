package models

import "time"

type Organization struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organization_id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WorkflowTemplate struct {
	ID             int64           `json:"id"`
	OrganizationID int64           `json:"organization_id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Version        int             `json:"version"`
	IsActive       bool            `json:"is_active"`
	CreatedBy      int64           `json:"created_by"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Stages         []TemplateStage `json:"stages,omitempty"`
}

type TemplateStage struct {
	ID             int64  `json:"id"`
	TemplateID     int64  `json:"template_id"`
	Name           string `json:"name"`
	OrderIndex     int    `json:"order_index"`
	RequiredRole   string `json:"required_role"`
	ApprovalType   string `json:"approval_type"`
	TimeoutMinutes int    `json:"timeout_minutes"`
	EscalationRole string `json:"escalation_role,omitempty"`
}

type WorkflowInstance struct {
	ID              int64             `json:"id"`
	OrganizationID  int64             `json:"organization_id"`
	TemplateID      int64             `json:"template_id"`
	TemplateVersion int               `json:"template_version"`
	Title           string            `json:"title"`
	CurrentStageID  *int64            `json:"current_stage_id"`
	Status          string            `json:"status"`
	CreatedBy       int64             `json:"created_by"`
	AssignedTo      *int64            `json:"assigned_to"`
	EscalatedTo     *int64            `json:"escalated_to"`
	StartedAt       time.Time         `json:"started_at"`
	CompletedAt     *time.Time        `json:"completed_at"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	CurrentStage    *TemplateStage    `json:"current_stage,omitempty"`
	Template        *WorkflowTemplate `json:"template,omitempty"`
	CreatedByUser   *User             `json:"created_by_user,omitempty"`
	AssignedToUser  *User             `json:"assigned_to_user,omitempty"`
}

type AuditLog struct {
	ID          int64     `json:"id"`
	WorkflowID  int64     `json:"workflow_id"`
	ActorID     int64     `json:"actor_id"`
	ActorRole   string    `json:"actor_role"`
	FromStageID *int64    `json:"from_stage_id"`
	ToStageID   *int64    `json:"to_stage_id"`
	Action      string    `json:"action"`
	Reason      string    `json:"reason"`
	Timestamp   time.Time `json:"timestamp"`
	Actor       *User     `json:"actor,omitempty"`
}

type Lock struct {
	Key       string    `json:"key"`
	Owner     string    `json:"owner"`
	ExpiresAt time.Time `json:"expires_at"`
}
