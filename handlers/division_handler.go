package handlers

import (
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DivisionHandler defines the interface for division related handlers.
type DivisionHandler interface {
	CreateDivision(c *gin.Context)
	GetDivisions(c *gin.Context)
	GetDivisionByID(c *gin.Context)
	UpdateDivision(c *gin.Context)
	DeleteDivision(c *gin.Context)
}

// divisionHandler is the concrete implementation of DivisionHandler.
type divisionHandler struct {
	divisionService services.DivisionService
}

// NewDivisionHandler creates a new instance of DivisionHandler.
func NewDivisionHandler(divisionService services.DivisionService) DivisionHandler {
	return &divisionHandler{
		divisionService: divisionService,
	}
}

func (h *divisionHandler) CreateDivision(c *gin.Context) {
	var input models.DivisionTable
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}
	input.CompanyID = companyID.(uint)

	createdDivision, err := h.divisionService.CreateDivision(&input)
	if err != nil {
		helper.SendError(c, http.StatusConflict, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Division created successfully", createdDivision)
}

func (h *divisionHandler) GetDivisions(c *gin.Context) {
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}

	divisions, err := h.divisionService.GetDivisionsByCompanyID(companyID.(uint))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve divisions")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Divisions fetched successfully", divisions)
}

func (h *divisionHandler) GetDivisionByID(c *gin.Context) {
	divisionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid division ID format")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}

	division, err := h.divisionService.GetDivisionByID(uint(divisionID), companyID.(uint))
	if err != nil {
		helper.SendError(c, http.StatusNotFound, "Division not found or does not belong to the company")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Division fetched successfully", division)
}

func (h *divisionHandler) UpdateDivision(c *gin.Context) {
	divisionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid division ID format")
		return
	}

	var input models.DivisionTable
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}

	input.ID = uint(divisionID)
	updatedDivision, err := h.divisionService.UpdateDivision(&input, companyID.(uint))
	if err != nil {
		helper.SendError(c, http.StatusConflict, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Division updated successfully", updatedDivision)
}

func (h *divisionHandler) DeleteDivision(c *gin.Context) {
	divisionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid division ID format")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}

	if err := h.divisionService.DeleteDivision(uint(divisionID), companyID.(uint)); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete division")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Division deleted successfully", nil)
}