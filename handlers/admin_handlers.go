package handlers

import (
	"net/http"
	"strconv"
	"time"

	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
	"log"
	"sort"
)

// Activity represents a single recent activity for the dashboard.
type Activity struct {
	Type        string    `json:"type"` // e.g., "attendance", "leave_request"
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

// --- Company Handlers ---

type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func CreateCompany(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	company := &models.CompaniesTable{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := repository.CreateCompany(company); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create company.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Company created successfully.", company)
}

func GetCompanyByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID.")
		return
	}

	company, err := repository.GetCompanyByID(id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company.")
		return
	}

	if company == nil {
		helper.SendError(c, http.StatusNotFound, "Company not found.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Company retrieved successfully.", company)
}

// GetDashboardSummary handles fetching summary data for the admin dashboard.
func GetDashboardSummary(hub *websocket.Hub, c *gin.Context) {
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}

	compID, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	summaryData, err := GetDashboardSummaryData(int(compID))
	if err != nil {
		log.Printf("Error getting dashboard summary data: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve dashboard summary.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Dashboard summary fetched successfully.", summaryData)

	// Send update to WebSocket clients
	hub.SendDashboardUpdate(int(compID), summaryData)
}

// GetDashboardSummaryData fetches the raw summary data for a given company ID.
// This function is reusable by both HTTP handler and WebSocket push.
func GetDashboardSummaryData(companyID int) (gin.H, error) {
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

	summary := gin.H{
		"total_employees":  totalEmployees,
		"present_today":    presentToday,
		"absent_today":     absentToday,
		"on_leave_today":   onLeaveToday,
		"recent_activities": activities, // Include recent activities
	}
	return summary, nil
}