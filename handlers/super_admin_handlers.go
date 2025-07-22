package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"time"

	"go-face-auth/database"
	"go-face-auth/database/repository"
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

	// Register the client with the hub
	hub.Register <- client

	// Start the read and write pumps in separate goroutines
	go client.WritePump()
	go client.ReadPump(hub)

	// Send the initial dashboard data to the newly connected client
	go hub.BroadcastSuperAdminDashboardUpdate()
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
	Month        string  `json:"month"`
	Year         string  `json:"year"`
	TotalRevenue float64 `json:"total_revenue"`
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

// GetCustomPackageRequests handles fetching all custom package requests for superadmin.
func GetCustomPackageRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	requests, totalRecords, err := repository.GetCustomPackageRequestsPaginated(page, pageSize, search)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve custom package requests.")
		return
	}

	paginatedData := gin.H{
		"items":         requests,
		"total_records": totalRecords,
	}

	helper.SendSuccess(c, http.StatusOK, "Custom package requests retrieved successfully.", paginatedData)
}

// UpdateCustomPackageRequestStatus handles updating the status of a custom package request.
func UpdateCustomPackageRequestStatus(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request ID.")
		return
	}

	newStatus := c.Param("status") // 'contacted' or 'resolved'

	if newStatus != "contacted" && newStatus != "resolved" {
		helper.SendError(c, http.StatusBadRequest, "Invalid status provided.")
		return	}

	request, err := repository.GetCustomPackageRequestByID(uint(requestID))
	if err != nil || request == nil {
		helper.SendError(c, http.StatusNotFound, "Custom package request not found.")
		return
	}

	request.Status = newStatus

	if err := repository.UpdateCustomPackageRequest(request); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update request status.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Request status updated successfully.", nil)
}