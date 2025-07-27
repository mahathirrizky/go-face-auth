package handlers

import (
	"log"
	"net/http"
	"strconv"

	"go-face-auth/helper"
	"go-face-auth/middleware"
	"go-face-auth/services"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
)

// SuperAdminHandler defines the interface for super admin related handlers.
type SuperAdminHandler interface {
	GetSuperAdminDashboardSummary(c *gin.Context)
	SuperAdminDashboardWebSocketHandler(hub *websocket.Hub, c *gin.Context)
	GetCompanies(c *gin.Context)
	GetSubscriptions(c *gin.Context)
	GetRevenueSummary(c *gin.Context)
	GetCustomPackageRequests(c *gin.Context)
	UpdateCustomPackageRequestStatus(c *gin.Context)
}

// superAdminHandler is the concrete implementation of SuperAdminHandler.
type superAdminHandler struct {
	superAdminService services.SuperAdminService
}

// NewSuperAdminHandler creates a new instance of SuperAdminHandler.
func NewSuperAdminHandler(superAdminService services.SuperAdminService) SuperAdminHandler {
	return &superAdminHandler{
		superAdminService: superAdminService,
	}
}

// GetSuperAdminDashboardSummary handles fetching a summary for the superadmin dashboard.
func (h *superAdminHandler) GetSuperAdminDashboardSummary(c *gin.Context) {
	summary, err := h.superAdminService.GetSuperAdminDashboardSummary()
	if err != nil {
		log.Printf("Error fetching super admin dashboard summary: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve super admin dashboard summary.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "SuperAdmin dashboard summary retrieved successfully.",
		"data":    summary,
	})
}

// SuperAdminDashboardWebSocketHandler handles WebSocket connections for superadmin dashboard updates.
func (h *superAdminHandler) SuperAdminDashboardWebSocketHandler(hub *websocket.Hub, c *gin.Context) {
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
func (h *superAdminHandler) GetCompanies(c *gin.Context) {
	companies, err := h.superAdminService.GetCompanies()
	if err != nil {
		log.Printf("Error fetching companies: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve companies.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Companies retrieved successfully.", companies)
}

// GetSubscriptions handles fetching a list of all subscriptions.
func (h *superAdminHandler) GetSubscriptions(c *gin.Context) {
	companies, err := h.superAdminService.GetSubscriptions()
	if err != nil {
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
func (h *superAdminHandler) GetRevenueSummary(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	monthlyRevenue, err := h.superAdminService.GetRevenueSummary(startDateStr, endDateStr)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve revenue summary.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Revenue summary retrieved successfully.", monthlyRevenue)
}

// GetCustomPackageRequests handles fetching all custom package requests for superadmin.
func (h *superAdminHandler) GetCustomPackageRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	requests, totalRecords, err := h.superAdminService.GetCustomPackageRequests(page, pageSize, search)
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
func (h *superAdminHandler) UpdateCustomPackageRequestStatus(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request ID.")
		return
	}

	newStatus := c.Param("status") // 'contacted' or 'resolved'

	if newStatus != "contacted" && newStatus != "resolved" {
		helper.SendError(c, http.StatusBadRequest, "Invalid status provided.")
		return	}

	if err := h.superAdminService.UpdateCustomPackageRequestStatus(uint(requestID), newStatus); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Request status updated successfully.", nil)
}
