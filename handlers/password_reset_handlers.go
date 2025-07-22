package handlers

import (

	"go-face-auth/services"
	"go-face-auth/helper"

	"net/http"


	"github.com/gin-gonic/gin"

)

// ForgotPasswordRequest defines the structure for a forgot password request.
type ForgotPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
}

// ForgotPassword handles the request to initiate a password reset.
func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := services.ForgotPassword(req.Email, "admin"); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SendSuccess(c, http.StatusOK, "If an account with that email exists, a password reset link has been sent.", nil)
}

// ForgotPasswordEmployee handles the request to initiate a password reset for an employee.
func ForgotPasswordEmployee(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := services.ForgotPassword(req.Email, "employee"); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SendSuccess(c, http.StatusOK, "If an account with that email exists, a password reset link has been sent.", nil)
}

// ResetPasswordRequest defines the structure for a reset password request.
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPassword handles the request to reset the password using a token.
func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := services.ResetPassword(req.Token, req.NewPassword); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Password has been reset successfully.", nil)
}

