package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"go-face-auth/database/repository"
	"go-face-auth/helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// SuperAdminLoginRequest represents the request body for super admin login.
type SuperAdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AdminCompanyLoginRequest represents the request body for admin company login.
type AdminCompanyLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// EmployeeLoginRequest represents the request body for employee login.
type EmployeeLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// generateToken generates a JWT token with given claims.
func generateToken(id int, role string, companyID int) (string, error) {
	claims := jwt.MapClaims{
		"id":        id,
		"role":      role,
		"companyID": companyID, // Add companyID to claims
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("supersecretjwtkeythatshouldbechangedinproduction")
		log.Println("WARNING: JWT_SECRET environment variable not set for token generation. Using default secret.")
	}

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}
	log.Printf("Token generated successfully (first 10 chars): %s", tokenString[:10])
	return tokenString, nil
}

// LoginSuperAdmin handles super admin authentication and JWT token generation.
func LoginSuperAdmin(c *gin.Context) {
	var req SuperAdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	superAdmin, err := repository.GetSuperAdminByEmail(req.Email)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve super user.")
		return
	}
	if superAdmin == nil {
		helper.SendError(c, http.StatusUnauthorized, "Invalid credentials.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(superAdmin.Password), []byte(req.Password)); err != nil {
		helper.SendError(c, http.StatusUnauthorized, "Invalid credentials.")
		return
	}

	tokenString, err := generateToken(superAdmin.ID, superAdmin.Role, 0)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate token.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Super admin login successful.", gin.H{
		"token": tokenString,
		"user":  superAdmin,
	})
}

// LoginAdminCompany handles admin company authentication and JWT token generation.
func LoginAdminCompany(c *gin.Context) {
	var req AdminCompanyLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	adminCompany, err := repository.GetAdminCompanyByEmail(req.Email)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve admin company.")
		return
	}
	if adminCompany == nil {
		helper.SendError(c, http.StatusUnauthorized, "Invalid credentials.")
		return
	}

	// Check if email is confirmed
	if !adminCompany.IsConfirmed {
		helper.SendError(c, http.StatusUnauthorized, "Email not confirmed. Please check your inbox for a confirmation link.")
		return
	}

	// Check company subscription status
	// if adminCompany.Company.SubscriptionStatus != "active" {
	// 	helper.SendError(c, http.StatusForbidden, "Company subscription is not active. Please complete payment.")
	// 	return
	// }

	if err := bcrypt.CompareHashAndPassword([]byte(adminCompany.Password), []byte(req.Password)); err != nil {
		helper.SendError(c, http.StatusUnauthorized, "Kata sandi salah.")
		return
	}

	log.Printf("Attempting to generate token for AdminCompany ID: %d, Role: %s, CompanyID: %d", adminCompany.ID, adminCompany.Role, adminCompany.CompanyID)
	tokenString, err := generateToken(adminCompany.ID, adminCompany.Role, adminCompany.CompanyID)
	if err != nil {
		log.Printf("Error generating token for AdminCompany ID %d: %v", adminCompany.ID, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate token.")
		return
	}
	log.Printf("Generated token (first 10 chars): %s", tokenString[:10])

		helper.SendSuccess(c, http.StatusOK, "Admin company login successful.", gin.H{
			"token": tokenString,
			"user":  adminCompany,
		})
}

// LoginEmployee handles employee authentication and JWT token generation.
func LoginEmployee(c *gin.Context) {
	var req EmployeeLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	employee, err := repository.GetEmployeeByEmail(req.Email)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employee.")
		return
	}
	if employee == nil {
		helper.SendError(c, http.StatusUnauthorized, "Invalid credentials.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(req.Password)); err != nil {
		helper.SendError(c, http.StatusUnauthorized, "Invalid credentials.")
		return
	}

	tokenString, err := generateToken(employee.ID, "employee", employee.CompanyID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate token.")
		return
	}

	// Get company details to include in the response
	locations, err := repository.GetAttendanceLocationsByCompanyID(employee.CompanyID)
	if err != nil {
		log.Printf("Error retrieving attendance locations for company %d: %v", employee.CompanyID, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company location information.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee login successful.", gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":                   employee.ID,
			"email":                employee.Email,
			"name":                 employee.Name,
			"position":             employee.Position,
			"role":                 "employee",
			"companyID":            employee.CompanyID,
			"attendance_locations": locations, // Return all valid locations
		},
	})
}
