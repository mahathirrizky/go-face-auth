package handlers

import (
	
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"


	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSubscriptionPackages retrieves all available subscription packages.
func GetSubscriptionPackages(c *gin.Context) {
	packages, err := repository.GetSubscriptionPackages()
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve subscription packages")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription packages retrieved successfully", packages)
}

// CreateSubscriptionPackage creates a new subscription package.
func CreateSubscriptionPackage(c *gin.Context) {
	var newPackage models.SubscriptionPackageTable
	if err := c.ShouldBindJSON(&newPackage); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := repository.CreateSubscriptionPackage(&newPackage); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Subscription package created successfully", newPackage)
}

// UpdateSubscriptionPackage updates an existing subscription package.
func UpdateSubscriptionPackage(c *gin.Context) {
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

	if err := repository.UpdateSubscriptionPackageFields(packageID, updates); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription package updated successfully", nil)
}

// DeleteSubscriptionPackage deletes a subscription package.
func DeleteSubscriptionPackage(c *gin.Context) {
	packageID := c.Param("id")

	if err := repository.DeleteSubscriptionPackage(packageID); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription package deleted successfully", nil)
}
