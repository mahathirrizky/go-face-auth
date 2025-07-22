package handlers

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
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

		// 2. Fetch Company details
		company, err := repository.GetCompanyByID(int(compID))
		if err != nil || company == nil {
			helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company information.")
			return
		}

		// 3. Fetch AdminCompany details
		admin, err := repository.GetAdminCompanyByID(int(adminID))
		if err != nil || admin == nil {
			helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve admin information.")
			return
		}

		// 4. Bind message and phone from request body
		var req CreateCustomPackageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.SendError(c, http.StatusBadRequest, helper.GetValidationError(err))
			return
		}

		// 5. Populate CustomPackageRequest model
		customRequest := &models.CustomPackageRequest{
			CompanyID:   compID,
			Email:       admin.Email,
			Phone:       req.Phone,
			CompanyName: company.Name,
			Message:     req.Message,
			Status:      "pending",
		}

		// 6. Create the request in DB
		if err := repository.CreateCustomPackageRequest(customRequest); err != nil {
			helper.SendError(c, http.StatusInternalServerError, "Failed to submit custom package request.")
			return
		}

		// 7. Trigger WebSocket notification for superadmins
		go hub.SendSuperAdminNotification(websocket.SuperAdminNotificationPayload{
			Type:       "new_custom_package_request",
			Message:    fmt.Sprintf("Permintaan paket kustom baru dari %s.", company.Name),
			CompanyID:  compID,
			CompanyName: company.Name,
		})

		helper.SendSuccess(c, http.StatusCreated, "Custom package request submitted successfully.", nil)
	}
}
