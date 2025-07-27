package handlers

import (
	"go-face-auth/services"
	"go-face-auth/helper"
	"go-face-auth/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

// SubscriptionPackageHandler defines the interface for subscription package related handlers.
type SubscriptionPackageHandler interface {
	GetSubscriptionPackages(c *gin.Context)
	CreateSubscriptionPackage(c *gin.Context)
	UpdateSubscriptionPackage(c *gin.Context)
	DeleteSubscriptionPackage(c *gin.Context)
}

// subscriptionPackageHandler is the concrete implementation of SubscriptionPackageHandler.
type subscriptionPackageHandler struct {
	subscriptionPackageService services.SubscriptionPackageService
}

// NewSubscriptionPackageHandler creates a new instance of SubscriptionPackageHandler.
func NewSubscriptionPackageHandler(subscriptionPackageService services.SubscriptionPackageService) SubscriptionPackageHandler {
	return &subscriptionPackageHandler{
		subscriptionPackageService: subscriptionPackageService,
	}
}

// GetSubscriptionPackages retrieves all available subscription packages.
func (h *subscriptionPackageHandler) GetSubscriptionPackages(c *gin.Context) {
	packages, err := h.subscriptionPackageService.GetSubscriptionPackages()
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve subscription packages")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription packages retrieved successfully", packages)
}

// CreateSubscriptionPackage creates a new subscription package.
func (h *subscriptionPackageHandler) CreateSubscriptionPackage(c *gin.Context) {
	var newPackage models.SubscriptionPackageTable
	if err := c.ShouldBindJSON(&newPackage); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.subscriptionPackageService.CreateSubscriptionPackage(&newPackage); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Subscription package created successfully", newPackage)
}

// UpdateSubscriptionPackage updates an existing subscription package.
func (h *subscriptionPackageHandler) UpdateSubscriptionPackage(c *gin.Context) {
	packageID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Prevent updating ID or other sensitive fields if necessary
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "updated_at")

	if err := h.subscriptionPackageService.UpdateSubscriptionPackage(packageID, updates); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription package updated successfully", nil)
}

// DeleteSubscriptionPackage deletes a subscription package.
func (h *subscriptionPackageHandler) DeleteSubscriptionPackage(c *gin.Context) {
	packageID := c.Param("id")

	if err := h.subscriptionPackageService.DeleteSubscriptionPackage(packageID); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription package deleted successfully", nil)
}