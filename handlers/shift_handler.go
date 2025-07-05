package handlers

import (
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateShift handles the creation of a new shift.
func CreateShift(c *gin.Context) {
	var shift models.ShiftsTable
	if err := c.ShouldBindJSON(&shift); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	companyID, _ := c.Get("companyID")
	shift.CompanyID = int(companyID.(float64))

	if err := repository.CreateShift(&shift); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create shift.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Shift created successfully.", shift)
}

// GetShiftsByCompany handles retrieving all shifts for a company.
func GetShiftsByCompany(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	shifts, err := repository.GetShiftsByCompanyID(int(companyID.(float64)))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve shifts.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shifts retrieved successfully.", shifts)
}

// UpdateShift handles updating an existing shift.
func UpdateShift(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid shift ID.")
		return
	}

	var shift models.ShiftsTable
	if err := c.ShouldBindJSON(&shift); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	shift.ID = id
	companyID, _ := c.Get("companyID")
	shift.CompanyID = int(companyID.(float64))

	if err := repository.UpdateShift(&shift); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update shift.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shift updated successfully.", shift)
}

// DeleteShift handles deleting a shift.
func DeleteShift(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid shift ID.")
		return
	}

	if err := repository.DeleteShift(id); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete shift.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shift deleted successfully.", nil)
}
