package handlers

import (
	"go-face-auth/helper"

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
	var req services.CreateDivisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}
	compID, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}
	req.CompanyID = uint(compID)

	createdDivision, err := h.divisionService.CreateDivision(req)
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

	compID, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	divisions, err := h.divisionService.GetDivisionsByCompanyID(uint(compID))
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

	compID, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	division, err := h.divisionService.GetDivisionByID(uint(divisionID), uint(compID))
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

	var req services.UpdateDivisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: Company ID not found")
		return
	}

	compID, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	updatedDivision, err := h.divisionService.UpdateDivision(uint(divisionID), uint(compID), req)
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

	compID, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	if err := h.divisionService.DeleteDivision(uint(divisionID), uint(compID)); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete division")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Division deleted successfully", nil)
}