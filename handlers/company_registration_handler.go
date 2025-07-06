package handlers

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"go-face-auth/helper"
	"net/http"
	"time"


	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterCompanyRequest defines the structure for the company registration request body.
type RegisterCompanyRequest struct {
	CompanyName         string `json:"company_name" binding:"required"`
	CompanyAddress      string `json:"company_address"`
	AdminEmail          string `json:"admin_email" binding:"required,email"`
	AdminPassword       string `json:"admin_password" binding:"required,min=6"`
	SubscriptionPackageID int   `json:"subscription_package_id" binding:"required"`
}

// RegisterCompany handles the registration of a new company and its admin user.
func RegisterCompany(c *gin.Context) {
	var req RegisterCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the subscription package exists
	var subPackage models.SubscriptionPackageTable
	if err := database.DB.First(&subPackage, req.SubscriptionPackageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.SendError(c, http.StatusBadRequest, "Subscription package not found")
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve subscription package")
		}
		return
	}

	// Hash the admin password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Start a database transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to start database transaction")
		return
	}

	// Create the new company
	now := time.Now()
	trialEndDate := now.AddDate(0, 0, 14)
	company := models.CompaniesTable{
		Name:                req.CompanyName,
		Address:             req.CompanyAddress,
		SubscriptionPackageID: req.SubscriptionPackageID,
		SubscriptionStatus:  "trial", // Set status to trial
		TrialStartDate:      &now,
		TrialEndDate:        &trialEndDate,
	}
	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		helper.SendError(c, http.StatusInternalServerError, "Failed to create company")
		return
	}

	// Create the admin company user
	adminCompany := models.AdminCompaniesTable{
		CompanyID: company.ID,
		Email:     req.AdminEmail,
		Password:  string(hashedPassword),
		Role:      "admin", // Default role for company admin
	}
	if err := tx.Create(&adminCompany).Error; err != nil {
		tx.Rollback()
		helper.SendError(c, http.StatusInternalServerError, "Failed to create company admin")
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to commit database transaction")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Company registered successfully. Your 14-day free trial has started.", gin.H{
		"company_id": company.ID,
		"admin_email": adminCompany.Email,
	})
}
