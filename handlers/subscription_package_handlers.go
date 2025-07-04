package handlers

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"go-face-auth/helper"

	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSubscriptionPackages retrieves all available subscription packages.
func GetSubscriptionPackages(c *gin.Context) {
	var packages []models.SubscriptionPackageTable
	if err := database.DB.Find(&packages).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve subscription packages")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription packages retrieved successfully", packages)
}
