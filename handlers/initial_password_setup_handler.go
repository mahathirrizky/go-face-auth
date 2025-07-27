package handlers

import (
	"go-face-auth/helper"
	"go-face-auth/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

// InitialPasswordSetupHandler defines the interface for initial password setup handlers.
type InitialPasswordSetupHandler interface {
	InitialPasswordSetup(c *gin.Context)
}

// initialPasswordSetupHandler is the concrete implementation of InitialPasswordSetupHandler.
type initialPasswordSetupHandler struct {
	initialPasswordSetupService services.InitialPasswordSetupService
}

// NewInitialPasswordSetupHandler creates a new instance of InitialPasswordSetupHandler.
func NewInitialPasswordSetupHandler(initialPasswordSetupService services.InitialPasswordSetupService) InitialPasswordSetupHandler {
	return &initialPasswordSetupHandler{
		initialPasswordSetupService: initialPasswordSetupService,
	}
}

// InitialPasswordSetupRequest defines the structure for an initial password setup request.
type InitialPasswordSetupRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// InitialPasswordSetup handles the request to set the initial password for an employee.
func (h *initialPasswordSetupHandler) InitialPasswordSetup(c *gin.Context) {
	var req InitialPasswordSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.initialPasswordSetupService.SetupInitialPassword(req.Token, req.Password); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Initial password has been set successfully.", nil)
}
