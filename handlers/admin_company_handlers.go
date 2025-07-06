package handlers

import (
	"net/http"
	"strconv"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
)

// AdminCompanyHandler handles HTTP requests related to admin companies.
type AdminCompanyHandler struct {
}

// NewAdminCompanyHandler creates a new AdminCompanyHandler.
func NewAdminCompanyHandler() *AdminCompanyHandler {
	return &AdminCompanyHandler{}
}

// CreateAdminCompany handles the creation of a new admin company.
func (h *AdminCompanyHandler) CreateAdminCompany(c *gin.Context) {
	var adminCompany models.AdminCompaniesTable
	if err := c.BindJSON(&adminCompany); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := repository.CreateAdminCompany(&adminCompany); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Admin company created successfully", adminCompany)
}

// GetAdminCompanyByCompanyID handles fetching an admin company by its CompanyID.
func (h *AdminCompanyHandler) GetAdminCompanyByCompanyID(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID")
		return
	}

	adminCompany, err := repository.GetAdminCompanyByCompanyID(companyID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if adminCompany == nil {
		helper.SendError(c, http.StatusNotFound, "Admin company not found for this company ID")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Admin company fetched successfully", adminCompany)
}

// GetAdminCompanyByEmployeeID handles fetching an admin company by its EmployeeID.
func (h *AdminCompanyHandler) GetAdminCompanyByEmployeeID(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	adminCompany, err := repository.GetAdminCompanyByEmployeeID(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if adminCompany == nil {
		helper.SendError(c, http.StatusNotFound, "Admin company not found for this employee ID")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Admin company fetched successfully", adminCompany)
}

// CheckAndNotifySubscriptions checks subscription statuses and sends notifications.
func (h *AdminCompanyHandler) CheckAndNotifySubscriptions(c *gin.Context) {
	var companies []models.CompaniesTable
	// Fetch companies with active or trial subscriptions
	if err := database.DB.Preload("AdminCompaniesTable").Where("subscription_status = ? OR subscription_status = ?", "active", "trial").Find(&companies).Error; err != nil {
		log.Printf("Error fetching companies for subscription check: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to fetch companies for subscription check")
		return
	}

	now := time.Now()
	adminFrontendURL := helper.GetFrontendAdminBaseURL()

	for _, company := range companies {
		// Determine the relevant end date (TrialEndDate for trial, SubscriptionEndDate for active)
		var endDate *time.Time
		var statusToUpdate string

		if company.SubscriptionStatus == "trial" && company.TrialEndDate != nil {
			endDate = company.TrialEndDate
			statusToUpdate = "expired_trial"
		} else if company.SubscriptionStatus == "active" && company.SubscriptionEndDate != nil {
			endDate = company.SubscriptionEndDate
			statusToUpdate = "expired"
		} else {
			continue // Skip if no valid end date or status is not active/trial
		}

		if endDate == nil {
			continue // Should not happen if logic above is correct, but for safety
		}

		daysRemaining := int(endDate.Sub(now).Hours() / 24)

		// Ensure there's at least one admin email to send to
		var adminEmail string
		if len(company.AdminCompaniesTable) > 0 {
			adminEmail = company.AdminCompaniesTable[0].Email
		} else {
			log.Printf("No admin email found for company %d (%s). Skipping notification.", company.ID, company.Name)
			continue
		}

		// Send reminders
		if daysRemaining <= 7 && daysRemaining > 0 {
			log.Printf("Sending subscription reminder to %s for company %s. %d days remaining.", adminEmail, company.Name, daysRemaining)
			if err := helper.SendSubscriptionReminderEmail(adminEmail, company.Name, daysRemaining, adminFrontendURL+"/dashboard/subscribe"); err != nil {
				log.Printf("Failed to send reminder email to %s: %v", adminEmail, err)
			}
		}

		// Handle expired subscriptions
		if daysRemaining <= 0 {
			log.Printf("Subscription for company %s has expired. Updating status to %s.", company.Name, statusToUpdate)
			company.SubscriptionStatus = statusToUpdate
			if err := database.DB.Save(&company).Error; err != nil {
				log.Printf("Failed to update subscription status for company %s: %v", company.Name, err)
			} else {
				log.Printf("Subscription status for company %s updated to %s.", company.Name, statusToUpdate)
				// Send expired notification email
				if err := helper.SendSubscriptionExpiredEmail(adminEmail, company.Name, adminFrontendURL+"/dashboard/subscribe"); err != nil {
					log.Printf("Failed to send expired email to %s: %v", adminEmail, err)
				}
			}
		}
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription check and notifications processed.", nil)
}
