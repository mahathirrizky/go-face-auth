package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"sort"
)

type AdminCompanyService interface {
	CreateAdminCompany(adminCompany *models.AdminCompaniesTable) error
	GetAdminCompanyByCompanyID(companyID int) (*models.AdminCompaniesTable, error)
	GetAdminCompanyByEmployeeID(employeeID int) (*models.AdminCompaniesTable, error)
	ChangeAdminPassword(adminID int, oldPassword, newPassword string) error
	CheckAndNotifySubscriptions() error
	GetDashboardSummaryData(companyID int) (map[string]interface{}, error)
}

type adminCompanyService struct {
	adminCompanyRepo repository.AdminCompanyRepository
	companyRepo      repository.CompanyRepository
	employeeRepo     repository.EmployeeRepository
	attendanceRepo   repository.AttendanceRepository
	leaveRepo        repository.LeaveRequestRepository
	db               *gorm.DB
}

func NewAdminCompanyService(adminCompanyRepo repository.AdminCompanyRepository, companyRepo repository.CompanyRepository, employeeRepo repository.EmployeeRepository, attendanceRepo repository.AttendanceRepository, leaveRepo repository.LeaveRequestRepository, db *gorm.DB) AdminCompanyService {
	return &adminCompanyService{
		adminCompanyRepo: adminCompanyRepo,
		companyRepo:      companyRepo,
		employeeRepo:     employeeRepo,
		attendanceRepo:   attendanceRepo,
		leaveRepo:        leaveRepo,
		db:               db,
	}
}

func (s *adminCompanyService) CreateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	return s.adminCompanyRepo.CreateAdminCompany(adminCompany)
}

func (s *adminCompanyService) GetAdminCompanyByCompanyID(companyID int) (*models.AdminCompaniesTable, error) {
	return s.adminCompanyRepo.GetAdminCompanyByCompanyID(companyID)
}

func (s *adminCompanyService) GetAdminCompanyByEmployeeID(employeeID int) (*models.AdminCompaniesTable, error) {
	return s.adminCompanyRepo.GetAdminCompanyByEmployeeID(employeeID)
}


func (s *adminCompanyService) ChangeAdminPassword(adminID int, oldPassword, newPassword string) error {
	// 1. Fetch the current admin user from the database
	admin, err := s.adminCompanyRepo.GetAdminCompanyByID(adminID)
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
	if err := s.adminCompanyRepo.ChangeAdminPassword(adminID, string(newPasswordHash)); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

func (s *adminCompanyService) CheckAndNotifySubscriptions() error {
	companies, err := s.companyRepo.GetAllActiveCompanies()
	if err != nil {
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
			if err := s.companyRepo.UpdateCompany(&company); err != nil {
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
type Activity struct {
	Type        string    `json:"type"` // e.g., "attendance", "leave_request"
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}
func (s *adminCompanyService) GetDashboardSummaryData(companyID int) (map[string]interface{}, error) {
	totalEmployees, err := s.employeeRepo.GetTotalEmployeesByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	presentToday, err := s.attendanceRepo.GetPresentEmployeesCountToday(companyID)
	if err != nil {
		return nil, err
	}

	absentToday, err := s.employeeRepo.GetAbsentEmployeesCountToday(companyID, presentToday) // Logic might need adjustment
	if err != nil {
		return nil, err
	}

	onLeaveToday, err := s.leaveRepo.GetOnLeaveEmployeesCountToday(companyID)
	if err != nil {
		return nil, err
	}

	limit := 10
	activities := []Activity{}

	// Fetch recent attendances
	attendances, err := s.attendanceRepo.GetRecentAttendancesByCompanyID(companyID, limit)
	if err == nil {
		for _, att := range attendances {
			if att.Employee.Name != "" {
				description := ""
				if att.CheckOutTime == nil {
					description = att.Employee.Name + " absen masuk pada " + att.CheckInTime.Format("15:04")
				} else {
					description = att.Employee.Name + " absen keluar pada " + att.CheckOutTime.Format("15:04")
				}
				activities = append(activities, Activity{
					Type: "attendance", Description: description, Timestamp: att.CheckInTime,
				})
			}
		}
	}

	// Fetch recent leave requests
	leaveRequests, err := s.leaveRepo.GetRecentLeaveRequestsByCompanyID(companyID, limit)
	if err == nil {
		for _, lr := range leaveRequests {
			if lr.Employee.Name != "" {
				description := lr.Employee.Name + " mengajukan " + lr.Type + " (" + lr.Status + ")"
				activities = append(activities, Activity{
					Type: "leave_request", Description: description, Timestamp: lr.CreatedAt,
				})
			}
		}
	}
	
	// Fetch recent overtime attendances
	overtimeAttendances, err := s.attendanceRepo.GetRecentOvertimeAttendancesByCompanyID(companyID, limit)
	if err == nil {
		for _, att := range overtimeAttendances {
			if att.Employee.Name != "" {
				description := ""
				if att.Status == "overtime_in" {
					description = att.Employee.Name + " mulai lembur pada " + att.CheckInTime.Format("15:04")
				} else if att.Status == "overtime_out" && att.CheckOutTime != nil {
					description = att.Employee.Name + " selesai lembur pada " + att.CheckOutTime.Format("15:04")
				}
				if description != "" {
					activities = append(activities, Activity{
						Type: "overtime", Description: description, Timestamp: att.CheckInTime,
					})
				}
			}
		}
	}

	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Timestamp.After(activities[j].Timestamp)
	})

	if len(activities) > limit {
		activities = activities[:limit]
	}

	summary := map[string]interface{}{
		"total_employees":   totalEmployees,
		"present_today":     presentToday,
		"absent_today":      absentToday,
		"on_leave_today":    onLeaveToday,
		"recent_activities": activities,
	}
	return summary, nil
}
