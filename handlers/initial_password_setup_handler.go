package handlers

import (
	"go-face-auth/services"
	"go-face-auth/helper"

	
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitialPasswordSetupRequest defines the structure for an initial password setup request.
type InitialPasswordSetupRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// InitialPasswordSetup handles the request to set the initial password for an employee.
func InitialPasswordSetup(c *gin.Context) {
	var req InitialPasswordSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := services.SetupInitialPassword(req.Token, req.Password); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Initial password has been set successfully.", nil)
}