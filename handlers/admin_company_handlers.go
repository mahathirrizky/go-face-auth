package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
)

// AdminCompanyHandler handles HTTP requests related to admin companies.
type AdminCompanyHandler struct {
}

// NewAdminCompanyHandler creates a new AdminCompanyHandler.
func NewAdminCompanyHandler() *AdminCompanyHandler {
	return &AdminCompanyHandler{}
}

// CreateAdminCompany handles the creation of a new admin company.
func (h *AdminCompanyHandler) CreateAdminCompany(c *gin.Context) {
	var adminCompany models.AdminCompaniesTable
	if err := c.BindJSON(&adminCompany); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := repository.CreateAdminCompany(&adminCompany); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Admin company created successfully", adminCompany)
}

// GetAdminCompanyByCompanyID handles fetching an admin company by its CompanyID.
func (h *AdminCompanyHandler) GetAdminCompanyByCompanyID(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID")
		return
	}

	adminCompany, err := repository.GetAdminCompanyByCompanyID(companyID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if adminCompany == nil {
		helper.SendError(c, http.StatusNotFound, "Admin company not found for this company ID")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Admin company fetched successfully", adminCompany)
}

// GetAdminCompanyByEmployeeID handles fetching an admin company by its EmployeeID.
func (h *AdminCompanyHandler) GetAdminCompanyByEmployeeID(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	adminCompany, err := repository.GetAdminCompanyByEmployeeID(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if adminCompany == nil {
		helper.SendError(c, http.StatusNotFound, "Admin company not found for this employee ID")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Admin company fetched successfully", adminCompany)
}
