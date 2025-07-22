package handlers

import (
	"fmt"
	"go-face-auth/services"
	"go-face-auth/helper"

	"go-face-auth/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCustomPackageRequest handles the submission of a custom package request.
type CreateCustomPackageRequest struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func HandleCustomPackageRequest(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get companyID and adminID from context
		companyIDClaim, exists := c.Get("companyID")
		if !exists {
			helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token.")
			return
		}
		compID := uint(companyIDClaim.(float64))

		adminIDClaim, exists := c.Get("id") // Assuming "id" is the admin's ID in the token
		if !exists {
			helper.SendError(c, http.StatusUnauthorized, "Admin ID not found in token.")
			return
		}
		adminID := uint(adminIDClaim.(float64))

		var req CreateCustomPackageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.SendError(c, http.StatusBadRequest, helper.GetValidationError(err))
			return
		}

		customRequest, err := services.CreateCustomPackageRequest(compID, adminID, req.Phone, req.Message)
		if err != nil {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}

		// 7. Trigger WebSocket notification for superadmins
		go hub.SendSuperAdminNotification(websocket.SuperAdminNotificationPayload{
			Type:       "new_custom_package_request",
			Message:    fmt.Sprintf("Permintaan paket kustom baru dari %s.", customRequest.CompanyName),
			CompanyID:  compID,
			CompanyName: customRequest.CompanyName,
		})

		helper.SendSuccess(c, http.StatusCreated, "Custom package request submitted successfully.", nil)
	}
}
