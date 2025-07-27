package handlers

import (
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ShiftHandler defines the interface for shift related handlers.
type ShiftHandler interface {
	CreateShift(c *gin.Context)
	GetShiftsByCompany(c *gin.Context)
	UpdateShift(c *gin.Context)
	DeleteShift(c *gin.Context)
	SetDefaultShift(c *gin.Context)
}

// shiftHandler is the concrete implementation of ShiftHandler.
type shiftHandler struct {
	shiftService services.ShiftService
}

// NewShiftHandler creates a new instance of ShiftHandler.
func NewShiftHandler(shiftService services.ShiftService) ShiftHandler {
	return &shiftHandler{
		shiftService: shiftService,
	}
}

// CreateShift handles the creation of a new shift.
func (h *shiftHandler) CreateShift(c *gin.Context) {
	var shift models.ShiftsTable
	if err := c.ShouldBindJSON(&shift); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// It's safer to get company_id from the JWT middleware context
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}
	shift.CompanyID = int(companyID.(float64))

	createdShift, err := h.shiftService.CreateShift(&shift)
	if err != nil {
		// Check for specific business logic errors
		if err.Error() == "shift limit reached for your current plan" {
			helper.SendError(c, http.StatusForbidden, err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Shift created successfully", createdShift)
}

// GetShiftsByCompany handles retrieving all shifts for a company.
func (h *shiftHandler) GetShiftsByCompany(c *gin.Context) {
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}

	shifts, err := h.shiftService.GetShiftsByCompanyID(int(companyID.(float64)))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve shifts")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shifts retrieved successfully", shifts)
}

// UpdateShift handles updating an existing shift.
func (h *shiftHandler) UpdateShift(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid shift ID")
		return
	}

	var shift models.ShiftsTable
	if err := c.ShouldBindJSON(&shift); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}

	shift.ID = id
	shift.CompanyID = int(companyID.(float64))

	updatedShift, err := h.shiftService.UpdateShift(&shift)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update shift")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shift updated successfully", updatedShift)
}

// DeleteShift handles deleting a shift.
func (h *shiftHandler) DeleteShift(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid shift ID")
		return
	}

	// Optional: Add logic to verify the shift belongs to the company before deleting
	if err := h.shiftService.DeleteShift(id); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete shift")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shift deleted successfully", nil)
}

// SetDefaultShift handles setting a shift as the default for the company.
func (h *shiftHandler) SetDefaultShift(c *gin.Context) {
	type SetDefaultRequest struct {
		ShiftID int `json:"shift_id" binding:"required"`
	}
	var req SetDefaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request: shift_id is required")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}

	if err := h.shiftService.SetDefaultShift(int(companyID.(float64)), req.ShiftID); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to set default shift")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Default shift set successfully", nil)
}
