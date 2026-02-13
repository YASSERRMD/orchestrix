package handlers

import (
	"net/http"
	"strconv"

	"orchestrix/internal/middleware"
	"orchestrix/internal/models"
	"orchestrix/internal/services"

	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	workflowService *services.WorkflowService
}

func NewWorkflowHandler(workflowService *services.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{workflowService: workflowService}
}

func (h *WorkflowHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	status := c.Query("status")
	var assignedTo *int64
	if role == "viewer" || role == "operator" {
		assignedTo = &userID
	}

	workflows, err := h.workflowService.GetByOrgID(orgID, status, assignedTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workflows)
}

func (h *WorkflowHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	workflow, err := h.workflowService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}
	c.JSON(http.StatusOK, workflow)
}

type CreateWorkflowRequest struct {
	TemplateID int64  `json:"template_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	AssignedTo *int64 `json:"assigned_to"`
}

func (h *WorkflowHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)

	var req CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instance := &models.WorkflowInstance{
		OrganizationID: orgID,
		TemplateID:     req.TemplateID,
		Title:          req.Title,
		CreatedBy:      userID,
		AssignedTo:     req.AssignedTo,
	}

	created, err := h.workflowService.Create(instance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

type WorkflowActionRequest struct {
	Reason string `json:"reason"`
}

func (h *WorkflowHandler) Advance(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	var req WorkflowActionRequest
	c.ShouldBindJSON(&req)

	workflow, err := h.workflowService.Advance(id, userID, role, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) Reject(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	var req WorkflowActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := h.workflowService.Reject(id, userID, role, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) Rollback(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	var req WorkflowActionRequest
	c.ShouldBindJSON(&req)

	workflow, err := h.workflowService.Rollback(id, userID, role, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) Cancel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	var req WorkflowActionRequest
	c.ShouldBindJSON(&req)

	workflow, err := h.workflowService.Cancel(id, userID, role, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) GetAudit(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	logs, err := h.workflowService.GetAuditLogs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
