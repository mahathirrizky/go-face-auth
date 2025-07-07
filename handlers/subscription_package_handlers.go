package handlers

import (
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
	id, err := strconv.ParseUint(packageID, 10, 64)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid package ID")
		return
	}

	var updatedPackage models.SubscriptionPackageTable
	if err := c.ShouldBindJSON(&updatedPackage); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var existingPackage models.SubscriptionPackageTable
	if err := database.DB.First(&existingPackage, id).Error; err != nil {
		helper.SendError(c, http.StatusNotFound, "Subscription package not found")
		return
	}

	// Update fields
	existingPackage.PackageName = updatedPackage.PackageName
	existingPackage.PriceMonthly = updatedPackage.PriceMonthly
	existingPackage.PriceYearly = updatedPackage.PriceYearly
	existingPackage.MaxEmployees = updatedPackage.MaxEmployees
	existingPackage.Features = updatedPackage.Features
	existingPackage.IsActive = updatedPackage.IsActive

	if err := database.DB.Save(&existingPackage).Error; err != nil {
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
