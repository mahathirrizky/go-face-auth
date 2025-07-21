package handlers

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/models"
	"go-face-auth/helper"
	"strconv"

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

// CreateSubscriptionPackage creates a new subscription package.
func CreateSubscriptionPackage(c *gin.Context) {
	var newPackage models.SubscriptionPackageTable
	if err := c.ShouldBindJSON(&newPackage); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := database.DB.Create(&newPackage).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Subscription package created successfully", newPackage)
}

// UpdateSubscriptionPackage updates an existing subscription package.
func UpdateSubscriptionPackage(c *gin.Context) {
	packageID := c.Param("id")

	var existingPackage models.SubscriptionPackageTable
	if err := database.DB.First(&existingPackage, packageID).Error; err != nil {
		helper.SendError(c, http.StatusNotFound, "Subscription package not found")
		return
	}
	
	// Bind the incoming JSON to the existing package struct
	if err := c.ShouldBindJSON(&existingPackage); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Log the received data
	
	fmt.Printf("Received data for update: %+v\n", existingPackage)

	// Use Updates to save the changes
	if err := database.DB.Model(&existingPackage).Updates(existingPackage).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription package updated successfully", existingPackage)
}

// DeleteSubscriptionPackage deletes a subscription package.
func DeleteSubscriptionPackage(c *gin.Context) {
	packageID := c.Param("id")
	id, err := strconv.ParseUint(packageID, 10, 64)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid package ID")
		return
	}

	if err := database.DB.Delete(&models.SubscriptionPackageTable{}, id).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete subscription package")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Subscription package deleted successfully", nil)
}
