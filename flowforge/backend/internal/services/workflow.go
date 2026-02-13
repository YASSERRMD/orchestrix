package services

import (
	"errors"
	"flowforge/internal/models"
	"flowforge/internal/repository"
	"time"
)

type WorkflowService struct {
	workflowRepo *repository.WorkflowRepository
	templateRepo *repository.TemplateRepository
	auditRepo    *repository.AuditRepository
	lockRepo     *repository.LockRepository
}

func NewWorkflowService() *WorkflowService {
	return &WorkflowService{
		workflowRepo: repository.NewWorkflowRepository(),
		templateRepo: repository.NewTemplateRepository(),
		auditRepo:    repository.NewAuditRepository(),
		lockRepo:     repository.NewLockRepository(),
	}
}

func (s *WorkflowService) GetByOrgID(orgID int64, status string, assignedTo *int64) ([]models.WorkflowInstance, error) {
	return s.workflowRepo.GetByOrgID(orgID, status, assignedTo)
}

func (s *WorkflowService) GetByID(id int64) (*models.WorkflowInstance, error) {
	return s.workflowRepo.GetByID(id)
}

func (s *WorkflowService) Create(instance *models.WorkflowInstance) (*models.WorkflowInstance, error) {
	template, err := s.templateRepo.GetByID(instance.TemplateID)
	if err != nil {
		return nil, errors.New("template not found")
	}

	if len(template.Stages) == 0 {
		return nil, errors.New("template has no stages")
	}

	instance.TemplateVersion = template.Version
	instance.Status = "pending"
	instance.CurrentStageID = &template.Stages[0].ID

	created, err := s.workflowRepo.Create(instance)
	if err != nil {
		return nil, err
	}

	s.auditRepo.Create(&models.AuditLog{
		WorkflowID: created.ID,
		ActorID:    instance.CreatedBy,
		ActorRole:  "creator",
		ToStageID:  created.CurrentStageID,
		Action:     "created",
		Reason:     "Workflow instance created",
	})

	return created, nil
}

func (s *WorkflowService) Advance(instanceID int64, actorID int64, actorRole string, reason string) (*models.WorkflowInstance, error) {
	instance, err := s.workflowRepo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}

	if instance.Status != "pending" && instance.Status != "escalated" {
		return nil, errors.New("workflow cannot be advanced in current state")
	}

	template, err := s.templateRepo.GetByID(instance.TemplateID)
	if err != nil {
		return nil, err
	}

	currentStageIdx := -1
	for i, stage := range template.Stages {
		if instance.CurrentStageID != nil && stage.ID == *instance.CurrentStageID {
			currentStageIdx = i
			break
		}
	}

	if currentStageIdx == -1 || currentStageIdx >= len(template.Stages)-1 {
		instance.Status = "completed"
		instance.CompletedAt = &instance.UpdatedAt
		s.workflowRepo.Update(instance)

		s.auditRepo.Create(&models.AuditLog{
			WorkflowID:  instance.ID,
			ActorID:     actorID,
			ActorRole:   actorRole,
			FromStageID: instance.CurrentStageID,
			Action:      "completed",
			Reason:      reason,
		})
		return instance, nil
	}

	nextStage := template.Stages[currentStageIdx+1]
	oldStageID := instance.CurrentStageID
	instance.CurrentStageID = &nextStage.ID
	instance.Status = "pending"

	if instance.Status == "escalated" {
		instance.EscalatedTo = nil
	}

	s.workflowRepo.Update(instance)

	s.auditRepo.Create(&models.AuditLog{
		WorkflowID:  instance.ID,
		ActorID:     actorID,
		ActorRole:   actorRole,
		FromStageID: oldStageID,
		ToStageID:   &nextStage.ID,
		Action:      "advanced",
		Reason:      reason,
	})

	return instance, nil
}

func (s *WorkflowService) Reject(instanceID int64, actorID int64, actorRole string, reason string) (*models.WorkflowInstance, error) {
	instance, err := s.workflowRepo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}

	if instance.Status == "completed" || instance.Status == "rejected" || instance.Status == "canceled" {
		return nil, errors.New("workflow cannot be rejected in current state")
	}

	instance.Status = "rejected"
	s.workflowRepo.Update(instance)

	s.auditRepo.Create(&models.AuditLog{
		WorkflowID:  instance.ID,
		ActorID:     actorID,
		ActorRole:   actorRole,
		FromStageID: instance.CurrentStageID,
		Action:      "rejected",
		Reason:      reason,
	})

	return instance, nil
}

func (s *WorkflowService) Rollback(instanceID int64, actorID int64, actorRole string, reason string) (*models.WorkflowInstance, error) {
	instance, err := s.workflowRepo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}

	if instance.Status != "pending" && instance.Status != "rejected" {
		return nil, errors.New("workflow cannot be rolled back in current state")
	}

	template, err := s.templateRepo.GetByID(instance.TemplateID)
	if err != nil {
		return nil, err
	}

	currentStageIdx := -1
	for i, stage := range template.Stages {
		if instance.CurrentStageID != nil && stage.ID == *instance.CurrentStageID {
			currentStageIdx = i
			break
		}
	}

	if currentStageIdx <= 0 {
		return nil, errors.New("workflow is at the first stage")
	}

	prevStage := template.Stages[currentStageIdx-1]
	oldStageID := instance.CurrentStageID
	instance.CurrentStageID = &prevStage.ID
	instance.Status = "pending"

	s.workflowRepo.Update(instance)

	s.auditRepo.Create(&models.AuditLog{
		WorkflowID:  instance.ID,
		ActorID:     actorID,
		ActorRole:   actorRole,
		FromStageID: oldStageID,
		ToStageID:   &prevStage.ID,
		Action:      "rolled_back",
		Reason:      reason,
	})

	return instance, nil
}

func (s *WorkflowService) Cancel(instanceID int64, actorID int64, actorRole string, reason string) (*models.WorkflowInstance, error) {
	instance, err := s.workflowRepo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}

	if instance.Status == "completed" || instance.Status == "canceled" {
		return nil, errors.New("workflow cannot be canceled in current state")
	}

	instance.Status = "canceled"
	s.workflowRepo.Update(instance)

	s.auditRepo.Create(&models.AuditLog{
		WorkflowID: instance.ID,
		ActorID:    actorID,
		ActorRole:  actorRole,
		Action:     "canceled",
		Reason:     reason,
	})

	return instance, nil
}

func (s *WorkflowService) GetAuditLogs(workflowID int64) ([]models.AuditLog, error) {
	return s.auditRepo.GetByWorkflowID(workflowID)
}

func (s *WorkflowService) GetTimedOutWorkflows(threshold time.Time) ([]models.WorkflowInstance, error) {
	return s.workflowRepo.GetTimedOutWorkflows(threshold)
}

func (s *WorkflowService) Escalate(instanceID int64) (*models.WorkflowInstance, error) {
	instance, err := s.workflowRepo.GetByID(instanceID)
	if err != nil {
		return nil, err
	}

	template, err := s.templateRepo.GetByID(instance.TemplateID)
	if err != nil {
		return nil, err
	}

	currentStageIdx := -1
	var currentStage *models.TemplateStage
	for i, stage := range template.Stages {
		if instance.CurrentStageID != nil && stage.ID == *instance.CurrentStageID {
			currentStageIdx = i
			currentStage = &template.Stages[i]
			break
		}
	}

	_ = currentStageIdx

	if currentStage == nil || currentStage.EscalationRole == "" {
		return nil, errors.New("no escalation configured for this stage")
	}

	oldStatus := instance.Status
	instance.Status = "escalated"
	s.workflowRepo.Update(instance)

	s.auditRepo.Create(&models.AuditLog{
		WorkflowID:  instance.ID,
		ActorID:     0,
		ActorRole:   "system",
		FromStageID: instance.CurrentStageID,
		Action:      "escalated",
		Reason:      "Stage timeout exceeded",
	})

	_ = oldStatus
	_ = currentStage

	return instance, nil
}
