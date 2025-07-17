package handlers

import (
	"encoding/json"
	"fmt" // Added fmt import
	"log"
	"net/http"
	"time"

	"go-face-auth/database"
	"go-face-auth/helper"
	"go-face-auth/middleware"
	"go-face-auth/models"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
)

// GetSuperAdminDashboardSummary handles fetching a summary for the superadmin dashboard.
func GetSuperAdminDashboardSummary(c *gin.Context) {
	var totalCompanies int64
	if err := database.DB.Model(&models.CompaniesTable{}).Count(&totalCompanies).Error; err != nil {
		log.Printf("Error counting total companies: %v", err)
		// helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve total companies.") // Removed helper.SendError
		return
	}

	var activeSubscriptions int64
	if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ?", "active").Count(&activeSubscriptions).Error; err != nil {
		log.Printf("Error counting active subscriptions: %v", err)
		// helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve active subscriptions.") // Removed helper.SendError
		return
	}

	var expiredSubscriptions int64
	if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ? OR subscription_status = ?", "expired", "expired_trial").Count(&expiredSubscriptions).Error; err != nil {
		log.Printf("Error counting expired subscriptions: %v", err)
		// helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve expired subscriptions.") // Removed helper.SendError
		return
	}

	var trialSubscriptions int64
	if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ?", "trial").Count(&trialSubscriptions).Error; err != nil {
		log.Printf("Error counting trial subscriptions: %v", err)
		return
	}

	// Fetch recent company registrations for recent activities
	var recentCompanies []models.CompaniesTable
	if err := database.DB.Order("created_at DESC").Limit(5).Find(&recentCompanies).Error; err != nil {
		log.Printf("Error fetching recent companies: %v", err)
		// Continue without recent activities if there's an error
	}

	recentActivities := make([]map[string]interface{}, len(recentCompanies))
	for i, company := range recentCompanies {
		recentActivities[i] = map[string]interface{}{
			"id":          company.ID,
			"description": fmt.Sprintf("Company %s registered", company.Name),
			"timestamp":   company.CreatedAt.UnixMilli(),
		}
	}

	// helper.SendSuccess(c, http.StatusOK, "SuperAdmin dashboard summary retrieved successfully.", gin.H{ // Removed helper.SendSuccess
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "SuperAdmin dashboard summary retrieved successfully.",
		"data": gin.H{
			"total_companies":       totalCompanies,
			"active_subscriptions":  activeSubscriptions,
			"expired_subscriptions": expiredSubscriptions,
			"trial_subscriptions":   trialSubscriptions,
			"recent_activities":     recentActivities,
		},
	})
}

// SuperAdminDashboardWebSocketHandler handles WebSocket connections for superadmin dashboard updates.
func SuperAdminDashboardWebSocketHandler(hub *websocket.Hub, c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("SuperAdmin WebSocket: Missing token")
		return
	}

	claims, err := middleware.ValidateToken(tokenString)
	if err != nil {
		log.Println("SuperAdmin WebSocket: Invalid token:", err)
		return
	}

	if claims["role"] != "super_admin" {
		log.Println("SuperAdmin WebSocket: Unauthorized role:", claims["role"])
		return
	}

	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("SuperAdmin WebSocket upgrade failed:", err)
		return
	}

	client := &websocket.Client{Conn: conn, Send: make(chan []byte, 256), Done: make(chan struct{})}
	client.CompanyID = 0 // Superadmin is not associated with a company

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump(hub)

	// Send initial data and then periodically update
	go func() {
		defer func() {
			hub.Unregister <- client
		}()
		for {
			select {
			case <-client.Done:
				log.Println("SuperAdmin WebSocket client done, stopping periodic updates.")
				return
			case <-time.After(5 * time.Second):
				var totalCompanies int64
				if err := database.DB.Model(&models.CompaniesTable{}).Count(&totalCompanies).Error; err != nil {
					log.Printf("Error counting total companies for WS: %v", err)
					totalCompanies = 0
				}

				var activeSubscriptions int64
				if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ?", "active").Count(&activeSubscriptions).Error; err != nil {
					log.Printf("Error counting active subscriptions for WS: %v", err)
					activeSubscriptions = 0
				}

				var expiredSubscriptions int64
				if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ? OR subscription_status = ?", "expired", "expired_trial").Count(&expiredSubscriptions).Error; err != nil {
					log.Printf("Error counting expired subscriptions for WS: %v", err)
					expiredSubscriptions = 0
				}

				var trialSubscriptions int64
				if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ?", "trial").Count(&trialSubscriptions).Error; err != nil {
					log.Printf("Error counting trial subscriptions for WS: %v", err)
					trialSubscriptions = 0
				}

				var wsRecentCompanies []models.CompaniesTable
				if err := database.DB.Order("created_at DESC").Limit(5).Find(&wsRecentCompanies).Error; err != nil {
					log.Printf("Error fetching recent companies for WebSocket: %v", err)
				}

				wsRecentActivities := make([]map[string]interface{}, len(wsRecentCompanies))
				for i, company := range wsRecentCompanies {
					wsRecentActivities[i] = map[string]interface{}{
						"id":          company.ID,
						"description": fmt.Sprintf("Company %s registered (WS)", company.Name),
						"timestamp":   company.CreatedAt.UnixMilli(),
					}
				}

				var wsMonthlyRevenue []MonthlyRevenue
				if err := database.DB.Model(&models.InvoiceTable{}).Select(
					"DATE_FORMAT(created_at, '%Y-%m') AS month, DATE_FORMAT(created_at, '%Y') AS year, SUM(amount) AS total_revenue").Where("status = ?", "paid").Group("month, year").Order("year DESC, month DESC").Scan(&wsMonthlyRevenue).Error; err != nil {
					log.Printf("Error fetching revenue summary for WebSocket: %v", err)
				}

				summary := gin.H{
					"total_companies":       totalCompanies,
					"active_subscriptions":  activeSubscriptions,
					"expired_subscriptions": expiredSubscriptions,
					"trial_subscriptions":   trialSubscriptions,
					"recent_activities":     wsRecentActivities,
					"monthly_revenue":       wsMonthlyRevenue,
				}

				jsonSummary, _ := json.Marshal(summary)

				select {
				case client.Send <- jsonSummary:
				default:
					log.Printf("Client send channel closed or full, unregistering client.")
					hub.Unregister <- client
					return
				}
			}
		}
	}()
}

// GetCompanies handles fetching a list of all companies.
func GetCompanies(c *gin.Context) {
	var companies []models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").Find(&companies).Error; err != nil {
		log.Printf("Error fetching companies: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve companies.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Companies retrieved successfully.", companies)
}

// GetSubscriptions handles fetching a list of all subscriptions.
func GetSubscriptions(c *gin.Context) {
	var companies []models.CompaniesTable
	// Fetch all companies with their subscription details
	if err := database.DB.Preload("SubscriptionPackage").Find(&companies).Error; err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve subscriptions.")
		return
	}

	// You might want to format this data more specifically for subscriptions
	// For now, returning company data which includes subscription info
	helper.SendSuccess(c, http.StatusOK, "Subscriptions retrieved successfully.", companies)
}

// MonthlyRevenue represents the structure for monthly revenue data.
type MonthlyRevenue struct {
	Month string  `json:"month"`
	Year  string  `json:"year"`
	Total float64 `json:"total_revenue"`
}

// GetRevenueSummary handles fetching a summary of revenue within a specified date range.
func GetRevenueSummary(c *gin.Context) {
	var monthlyRevenue []MonthlyRevenue

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	query := database.DB.Model(&models.InvoiceTable{}).Where("status = ?", "paid")

	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD.")
			return
		}
		query = query.Where("created_at >= ?", startDate)
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD.")
			return
		}
		// Add 23 hours, 59 minutes, 59 seconds to include the entire end day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("created_at <= ?", endDate)
	}

	// Query to get monthly revenue from paid invoices
	// Using DATE_FORMAT for MySQL
	if err := query.Select(
		"DATE_FORMAT(created_at, '%Y-%m') AS month, DATE_FORMAT(created_at, '%Y') AS year, SUM(amount) AS total_revenue").Group("month, year").Order("year DESC, month DESC").Scan(&monthlyRevenue).Error; err != nil {
		log.Printf("Error fetching revenue summary: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve revenue summary.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Revenue summary retrieved successfully.", monthlyRevenue)
}