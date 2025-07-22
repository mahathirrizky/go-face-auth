package services

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	return repository.CreateAdminCompany(adminCompany)
}

func GetAdminCompanyByCompanyID(companyID int) (*models.AdminCompaniesTable, error) {
	return repository.GetAdminCompanyByCompanyID(companyID)
}

func GetAdminCompanyByEmployeeID(employeeID int) (*models.AdminCompaniesTable, error) {
	return repository.GetAdminCompanyByEmployeeID(employeeID)
}

func ChangeAdminPassword(adminID int, oldPassword, newPassword string) error {
	// 1. Fetch the current admin user from the database
	admin, err := repository.GetAdminCompanyByID(adminID)
	if err != nil {
		return fmt.Errorf("failed to retrieve admin details: %w", err)
	}
	if admin == nil {
		return fmt.Errorf("admin user not found")
	}

	// 2. Compare the provided old password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("incorrect old password")
	}

	// 3. Hash the new password
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// 4. Update the password in the database
	if err := repository.ChangeAdminPassword(adminID, string(newPasswordHash)); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

func CheckAndNotifySubscriptions() error {
	var companies []models.CompaniesTable
	// Fetch companies with active or trial subscriptions
	if err := database.DB.Preload("AdminCompaniesTable").Where("subscription_status = ? OR subscription_status = ?", "active", "trial").Find(&companies).Error; err != nil {
		return fmt.Errorf("error fetching companies for subscription check: %w", err)
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

	return nil
}
