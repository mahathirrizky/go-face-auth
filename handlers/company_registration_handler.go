package handlers

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	confirmationToken := uuid.New().String()
	adminCompany := models.AdminCompaniesTable{
		CompanyID:         company.ID,
		Email:             req.AdminEmail,
		Password:          string(hashedPassword),
		Role:              "admin", // Default role for company admin
		ConfirmationToken: confirmationToken,
		IsConfirmed:       false,
	}
	if err := tx.Create(&adminCompany).Error; err != nil {
		tx.Rollback()
		helper.SendError(c, http.StatusInternalServerError, "Failed to create company admin")
		return
	}

	// Create a default shift for the new company
	defaultShift := models.ShiftsTable{
		CompanyID:          company.ID,
		Name:               "Shift Pagi",
		StartTime:          "09:00:00",
		EndTime:            "17:00:00",
		GracePeriodMinutes: 15,
		IsDefault:          true,
	}
	if err := tx.Create(&defaultShift).Error; err != nil {
		tx.Rollback()
		helper.SendError(c, http.StatusInternalServerError, "Failed to create default shift")
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to commit database transaction")
		return
	}

	// Send confirmation email
	frontendAdminBaseURL := os.Getenv("FRONTEND_ADMIN_BASE_URL")
	if frontendAdminBaseURL == "" {
		log.Println("WARNING: FRONTEND_ADMIN_BASE_URL environment variable is not set. Email confirmation link may not work correctly.")
	}
	confirmationLink := fmt.Sprintf("%s/confirm-email?token=%s", frontendAdminBaseURL, confirmationToken)

	go func() {
		err := helper.SendConfirmationEmail(req.AdminEmail, req.CompanyName, confirmationLink)
		if err != nil {
			log.Printf("Error sending confirmation email to %s: %v", req.AdminEmail, err)
		}
	}()

	helper.SendSuccess(c, http.StatusCreated, "Company registered successfully. Please check your email for a confirmation link.", gin.H{
		"company_id": company.ID,
		"admin_email": adminCompany.Email,
	})
}

// ConfirmEmail handles the email confirmation process for admin companies.
func ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		helper.SendError(c, http.StatusBadRequest, "Confirmation token is missing.")
		return
	}

	var adminCompany models.AdminCompaniesTable
	result := database.DB.Where("confirmation_token = ?", token).First(&adminCompany)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			helper.SendError(c, http.StatusNotFound, "Invalid or expired confirmation token.")
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve admin company.")
		}
		return
	}

	if adminCompany.IsConfirmed {
		helper.SendSuccess(c, http.StatusOK, "Email already confirmed.", nil)
		return
	}

	// Update confirmation status and clear token
	adminCompany.IsConfirmed = true
	adminCompany.ConfirmationToken = nil // Clear the token after use

	if err := database.DB.Save(&adminCompany).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to confirm email.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Email confirmed successfully. You can now log in.", nil)
}
