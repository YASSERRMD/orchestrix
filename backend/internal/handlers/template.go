package handlers

import (
	"net/http"
	"strconv"

	"orchestrix/internal/middleware"
	"orchestrix/internal/models"
	"orchestrix/internal/services"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	templateService *services.TemplateService
}

func NewTemplateHandler(templateService *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{templateService: templateService}
}

func (h *TemplateHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	templates, err := h.templateService.GetActiveByOrgID(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *TemplateHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	template, err := h.templateService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

func (h *TemplateHandler) GetVersions(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	versions, err := h.templateService.GetVersions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

type CreateTemplateRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Stages      []models.TemplateStage `json:"stages" binding:"required,min=1"`
}

func (h *TemplateHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := &models.WorkflowTemplate{
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		Version:        1,
		IsActive:       true,
		CreatedBy:      userID,
		Stages:         req.Stages,
	}

	created, err := h.templateService.Create(template)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *TemplateHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.templateService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	if existing.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	template := &models.WorkflowTemplate{
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		Version:        existing.Version,
		IsActive:       true,
		CreatedBy:      userID,
		Stages:         req.Stages,
	}

	if err := h.templateService.Update(template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template updated"})
}
