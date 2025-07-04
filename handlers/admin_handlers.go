package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// --- Company Handlers ---

type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func CreateCompany(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	company := &models.CompaniesTable{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := repository.CreateCompany(company); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create company.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Company created successfully.", company)
}

func GetCompanyByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID.")
		return
	}

	company, err := repository.GetCompanyByID(id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company.")
		return
	}

	if company == nil {
		helper.SendError(c, http.StatusNotFound, "Company not found.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Company retrieved successfully.", company)
}

// --- Employee Handlers ---

type CreateEmployeeRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	
}

func CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get company ID from JWT claims
	companyID, exists := c.Get("company_id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token")
		return
	}
	compID := companyID.(int)

	// Retrieve company and its subscription package
	var company models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").First(&company, compID).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company information")
		return
	}

	// Check current employee count
	var employeeCount int64
	if err := database.DB.Model(&models.EmployeesTable{}).Where("company_id = ?", compID).Count(&employeeCount).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to count existing employees")
		return
	}

	// Check if adding a new employee would exceed the package limit
	if employeeCount >= int64(company.SubscriptionPackage.MaxEmployees) {
		helper.SendError(c, http.StatusForbidden, "Employee limit reached for your subscription package")
		return
	}

	// Hash the employee password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	employee := &models.EmployeesTable{
		CompanyID: compID,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		
	}

	if err := repository.CreateEmployee(employee); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create employee")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Employee created successfully", employee)
}

func GetEmployeeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	employee, err := repository.GetEmployeeByID(id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employee.")
		return
	}

	if employee == nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee retrieved successfully.", employee)
}

func GetEmployeesByCompanyID(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID.")
		return
	}

	employees, err := repository.GetEmployeesByCompanyID(companyID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employees.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employees retrieved successfully.", employees)
}

// --- Face Image Handlers ---

func UploadFaceImage(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.PostForm("employee_id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	// Check if employee exists
	_, err = repository.GetEmployeeByID(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Image file is required.")
		return
	}

	// Create a unique filename
	filename := "employee_" + strconv.Itoa(employeeID) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(file.Filename)
	
	// Define the path to save the image
	savePath := filepath.Join("images", "employee_faces", filename)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create image directory.")
		return
	}

	// Save the file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to save image file.")
		return
	}

	faceImage := &models.FaceImagesTable{
		EmployeeID: employeeID,
		ImagePath:  savePath,
	}

	if err := repository.CreateFaceImage(faceImage); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to record face image in database.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Face image uploaded and recorded successfully.", gin.H{
		"employee_id": employeeID,
		"image_path":  savePath,
	})
}

func GetFaceImagesByEmployeeID(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("employee_id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	faceImages, err := repository.GetFaceImagesByEmployeeID(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve face images.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Face images retrieved successfully.", faceImages)
}
