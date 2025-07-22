package services

import (
	
	"go-face-auth/database/repository"
	
	"go-face-auth/models"
	"log"
	"sort"
	"time"
)

type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func CreateCompany(req CreateCompanyRequest) (*models.CompaniesTable, error) {
	company := &models.CompaniesTable{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := repository.CreateCompany(company); err != nil {
		return nil, err
	}

	return company, nil
}

func GetCompanyByID(companyID int) (*models.CompaniesTable, error) {
	return repository.GetCompanyByID(companyID)
}

// Activity represents a single recent activity for the dashboard.
type Activity struct {
	Type        string    `json:"type"` // e.g., "attendance", "leave_request"
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

// GetDashboardSummaryData fetches the raw summary data for a given company ID.
// This function is reusable by both HTTP handler and WebSocket push.
func GetDashboardSummaryData(companyID int) (map[string]interface{}, error) {
	// Fetch total employees
	totalEmployees, err := repository.GetTotalEmployeesByCompanyID(companyID)
	if err != nil {
		log.Printf("Error getting total employees for company %d: %v", companyID, err)
		return nil, err
	}

	// Fetch today's attendance (present, absent, leave)
	presentToday, err := repository.GetPresentEmployeesCountToday(companyID)
	if err != nil {
		log.Printf("Error getting present employees today for company %d: %v", companyID, err)
		return nil, err
	}

	absentToday, err := repository.GetAbsentEmployeesCountToday(companyID)
	if err != nil {
		log.Printf("Error getting absent employees today for company %d: %v", companyID, err)
		return nil, err
	}

	onLeaveToday, err := repository.GetOnLeaveEmployeesCountToday(companyID)
	if err != nil {
		log.Printf("Error getting on leave employees today for company %d: %v", companyID, err)
		return nil, err
	}

	// Fetch recent activities
	limit := 10 // Number of recent activities to fetch
	activities := []Activity{} // Initialize as empty slice

	// Fetch recent attendances
	attendances, err := repository.GetRecentAttendancesByCompanyID(companyID, limit)
	if err != nil {
		log.Printf("Error getting recent attendances: %v", err)
		// Continue even if there's an error, just log it
	} else {
		for _, att := range attendances {
			if att.Employee.Name == "" { // Check if Employee name is empty (e.g., if preload failed)
				continue // Skip this activity if employee data is missing
			}
			description := ""
			if att.CheckOutTime == nil {
				description = att.Employee.Name + " absen masuk pada " + att.CheckInTime.Format("15:04")
			} else {
				description = att.Employee.Name + " absen keluar pada " + att.CheckOutTime.Format("15:04")
			}
			activities = append(activities, Activity{
				Type:        "attendance",
				Description: description,
				Timestamp:   att.CheckInTime,
			})
		}
	}

	// Fetch recent leave requests
	leaveRequests, err := repository.GetRecentLeaveRequestsByCompanyID(companyID, limit)
	if err != nil {
		log.Printf("Error getting recent leave requests: %v", err)
		// Continue even if there's an error, just log it
	} else {
		for _, lr := range leaveRequests {
			if lr.Employee.Name == "" { // Check if Employee name is empty
				continue // Skip this activity if employee data is missing
			}
			description := lr.Employee.Name + " mengajukan " + lr.Type + " (" + lr.Status + ")"
			activities = append(activities, Activity{
				Type:        "leave_request",
				Description: description,
				Timestamp:   lr.CreatedAt,
			})
		}
	}

	// Fetch recent overtime attendances
	overtimeAttendances, err := repository.GetRecentOvertimeAttendancesByCompanyID(companyID, limit)
	if err != nil {
		log.Printf("Error getting recent overtime attendances: %v", err)
		// Continue even if there's an error, just log it
	} else {
		for _, att := range overtimeAttendances {
			if att.Employee.Name == "" {
				continue // Skip this activity if employee data is missing
			}
			description := ""
			if att.Status == "overtime_in" {
				description = att.Employee.Name + " mulai lembur pada " + att.CheckInTime.Format("15:04")
			} else if att.Status == "overtime_out" && att.CheckOutTime != nil {
				description = att.Employee.Name + " selesai lembur pada " + att.CheckOutTime.Format("15:04")
			}

			activities = append(activities, Activity{
				Type:        "overtime",
				Description: description,
				Timestamp:   att.CheckInTime, // Or CheckOutTime depending on the status
			})
		}
	}

	// Sort activities by timestamp (newest first)
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Timestamp.After(activities[j].Timestamp)
	})

	// Limit to the desired number of activities
	if len(activities) > limit {
		activities = activities[:limit]
	}

	summary := map[string]interface{}{
		"total_employees":  totalEmployees,
		"present_today":    presentToday,
		"absent_today":     absentToday,
		"on_leave_today":   onLeaveToday,
		"recent_activities": activities, // Include recent activities
	}
	return summary, nil
}
