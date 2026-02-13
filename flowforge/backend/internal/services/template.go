package services

import (
	"errors"
	"flowforge/internal/models"
	"flowforge/internal/repository"
)

type TemplateService struct {
	templateRepo *repository.TemplateRepository
}

func NewTemplateService() *TemplateService {
	return &TemplateService{
		templateRepo: repository.NewTemplateRepository(),
	}
}

func (s *TemplateService) GetByOrgID(orgID int64) ([]models.WorkflowTemplate, error) {
	return s.templateRepo.GetByOrgID(orgID)
}

func (s *TemplateService) GetActiveByOrgID(orgID int64) ([]models.WorkflowTemplate, error) {
	return s.templateRepo.GetActiveByOrgID(orgID)
}

func (s *TemplateService) GetByID(id int64) (*models.WorkflowTemplate, error) {
	return s.templateRepo.GetByID(id)
}

func (s *TemplateService) GetVersions(templateID int64) ([]models.WorkflowTemplate, error) {
	template, err := s.templateRepo.GetByID(templateID)
	if err != nil {
		return nil, err
	}

	allTemplates, err := s.templateRepo.GetByOrgID(template.OrganizationID)
	if err != nil {
		return nil, err
	}

	var versions []models.WorkflowTemplate
	for _, t := range allTemplates {
		if t.Name == template.Name {
			versions = append(versions, t)
		}
	}
	return versions, nil
}

func (s *TemplateService) Create(template *models.WorkflowTemplate) (*models.WorkflowTemplate, error) {
	if len(template.Stages) == 0 {
		return nil, errors.New("at least one stage is required")
	}

	for i := range template.Stages {
		template.Stages[i].OrderIndex = i
	}

	return s.templateRepo.Create(template)
}

func (s *TemplateService) Update(template *models.WorkflowTemplate) error {
	if len(template.Stages) == 0 {
		return errors.New("at least one stage is required")
	}

	for i := range template.Stages {
		template.Stages[i].OrderIndex = i
	}

	return s.templateRepo.Update(template)
}

func (s *TemplateService) GetStageByID(stageID int64) (*models.TemplateStage, error) {
	template, err := s.templateRepo.GetStagesByTemplateID(stageID)
	if err != nil {
		return nil, err
	}
	if len(template) == 0 {
		return nil, errors.New("stage not found")
	}
	return &template[0], nil
}
