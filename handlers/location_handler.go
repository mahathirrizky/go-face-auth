package handlers

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Utility function to get companyID from context
func getCompanyIDFromContext(c *gin.Context) (uint, error) {
	companyIDClaim, exists := c.Get("companyID")
	if !exists {
		return 0, fmt.Errorf("companyID not found in context")
	}

	companyID, ok := companyIDClaim.(float64) // JWT claims are float64
	if !ok {
		return 0, fmt.Errorf("invalid format for companyID in context")
	}

	return uint(companyID), nil
}

// CreateAttendanceLocation handles the creation of a new attendance location
func CreateAttendanceLocation(c *gin.Context) {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Invalid company information.")
		return
	}

	// Retrieve company and its subscription package
	var company models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").First(&company, companyID).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company information")
		return
	}

	// Check current location count
	var locationCount int64
	if err := database.DB.Model(&models.AttendanceLocation{}).Where("company_id = ?", companyID).Count(&locationCount).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to count existing locations")
		return
	}

	// Check if adding a new location would exceed the package limit
	if locationCount >= int64(company.SubscriptionPackage.MaxLocations) {
		helper.SendError(c, http.StatusForbidden, "Location limit reached for your subscription package")
		return
	}

	var location models.AttendanceLocation
	if err := c.ShouldBindJSON(&location); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set the company ID from the authenticated user, ignoring any value from the request body
	location.CompanyID = companyID

	createdLocation, err := repository.CreateAttendanceLocation(&location)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create attendance location.")
		return
	}

	c.JSON(http.StatusCreated, createdLocation)
}

// GetAttendanceLocations handles fetching all attendance locations for a company
func GetAttendanceLocations(c *gin.Context) {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Invalid company information.")
		return
	}

	locations, err := repository.GetAttendanceLocationsByCompanyID(companyID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to fetch attendance locations.")
		return
	}

	c.JSON(http.StatusOK, locations)
}

// UpdateAttendanceLocation handles updating an existing attendance location
func UpdateAttendanceLocation(c *gin.Context) {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Invalid company information.")
		return
	}

	locationID, err := strconv.ParseUint(c.Param("location_id"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid location ID.")
		return
	}

	// Check if the location to be updated actually belongs to the company
	existingLocation, err := repository.GetAttendanceLocationByID(uint(locationID))
	if err != nil {
		helper.SendError(c, http.StatusNotFound, "Location not found.")
		return
	}
	if existingLocation.CompanyID != companyID {
		helper.SendError(c, http.StatusForbidden, "Forbidden: You can only update locations for your own company.")
		return
	}

	var locationUpdates models.AttendanceLocation
	if err := c.ShouldBindJSON(&locationUpdates); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update fields
	existingLocation.Name = locationUpdates.Name
	existingLocation.Latitude = locationUpdates.Latitude
	existingLocation.Longitude = locationUpdates.Longitude
	existingLocation.Radius = locationUpdates.Radius

	updatedLocation, err := repository.UpdateAttendanceLocation(existingLocation)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update attendance location.")
		return
	}

	c.JSON(http.StatusOK, updatedLocation)
}

// DeleteAttendanceLocation handles deleting an attendance location
func DeleteAttendanceLocation(c *gin.Context) {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Invalid company information.")
		return
	}

	locationID, err := strconv.ParseUint(c.Param("location_id"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid location ID.")
		return
	}

	// Check if the location to be deleted actually belongs to the company
	existingLocation, err := repository.GetAttendanceLocationByID(uint(locationID))
	if err != nil {
		helper.SendError(c, http.StatusNotFound, "Location not found.")
		return
	}
	if existingLocation.CompanyID != companyID {
		helper.SendError(c, http.StatusForbidden, "Forbidden: You can only delete locations for your own company.")
		return
	}

	if err := repository.DeleteAttendanceLocation(uint(locationID)); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete attendance location.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance location deleted successfully"})
}