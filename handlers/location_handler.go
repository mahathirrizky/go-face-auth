package handlers

import (
	"fmt"
	"go-face-auth/services"
	"go-face-auth/helper"
	"go-face-auth/models"
	"net/http"
	"strconv"
	"errors"

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

	var location models.AttendanceLocation
	if err := c.ShouldBindJSON(&location); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	createdLocation, err := services.CreateAttendanceLocation(companyID, &location)
	if err != nil {
		if errors.Is(err, services.ErrLocationLimitReached) {
			helper.SendError(c, http.StatusForbidden, err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
		}
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

	locations, err := services.GetAttendanceLocationsByCompanyID(companyID)
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

	var locationUpdates models.AttendanceLocation
	if err := c.ShouldBindJSON(&locationUpdates); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedLocation, err := services.UpdateAttendanceLocation(companyID, uint(locationID), &locationUpdates)
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

	if err := services.DeleteAttendanceLocation(companyID, uint(locationID)); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete attendance location.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance location deleted successfully"})
}