package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"

	"github.com/gin-gonic/gin"
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

	company := &models.Company{
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
	CompanyID      int    `json:"company_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	EmployeeIDNumber string `json:"employee_id_number" binding:"required"`
}

func CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	// Check if company exists
	_, err := repository.GetCompanyByID(req.CompanyID)
	if err != nil {
		helper.SendError(c, http.StatusNotFound, "Company not found.")
		return
	}

	employee := &models.Employee{
		CompanyID:      req.CompanyID,
		Name:           req.Name,
		Email:          req.Email,
		EmployeeIDNumber: req.EmployeeIDNumber,
	}

	if err := repository.CreateEmployee(employee); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create employee.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Employee created successfully.", employee)
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

	faceImage := &models.FaceImage{
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
