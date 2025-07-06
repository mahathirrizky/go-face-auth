package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-face-auth/database"
	"go-face-auth/helper"
	"go-face-auth/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken parses and validates a JWT token string.
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		// Fallback for development if env var is not set
		jwtSecret = []byte("supersecretjwtkeythatshouldbechangedinproduction")
		log.Println("WARNING: JWT_SECRET environment variable not set for token validation. Using default secret.")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect: HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// AuthMiddleware validates JWT tokens from the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			helper.SendError(c, http.StatusUnauthorized, "Authorization header required.")
			c.Abort()
			return
		}

		// Check if the token string starts with "Bearer "
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			helper.SendError(c, http.StatusUnauthorized, "Invalid token format. Must be 'Bearer <token>'.")
			c.Abort()
			return
		}

		tokenString = tokenString[7:] // Remove "Bearer " prefix
		log.Printf("AuthMiddleware: Received token (first 10 chars): %s", tokenString[:10])

		claims, err := ValidateToken(tokenString)
		if err != nil {
			log.Printf("AuthMiddleware: Token validation error: %v", err)
			helper.SendError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		log.Printf("AuthMiddleware: Token valid. Claims: ID=%v, Role=%v, CompanyID=%v", claims["id"], claims["role"], claims["companyID"])

		// Subscription status check
		companyID, ok := claims["companyID"].(float64) // JWT numbers are float64
		if !ok || companyID == 0 {
			// If companyID is not present (e.g., for superuser), skip subscription check
			c.Set("id", claims["id"])
			c.Set("role", claims["role"])
			c.Set("companyID", claims["companyID"])
			c.Next()
			return
		}

		var company models.CompaniesTable
		if err := database.DB.First(&company, int(companyID)).Error; err != nil {
			helper.SendError(c, http.StatusForbidden, "Could not retrieve company information.")
			c.Abort()
			return
		}

		switch company.SubscriptionStatus {
		case "trial":
			if company.TrialEndDate != nil && time.Now().After(*company.TrialEndDate) {
				// Trial has expired, update status and block access
				company.SubscriptionStatus = "expired_trial"
				database.DB.Save(&company)
				helper.SendError(c, http.StatusForbidden, "Your free trial has expired. Please subscribe to continue.")
				c.Abort()
				return
			}
		case "active":
			// Active subscription, allow access
			break
		default:
			// Any other status (expired_trial, inactive, pending) is blocked
			helper.SendError(c, http.StatusForbidden, "Access denied. Please check your subscription status.")
			c.Abort()
			return
		}

		// Set claims in gin context
		c.Set("id", claims["id"])
		c.Set("role", claims["role"])
		c.Set("companyID", claims["companyID"]) // Set companyID in context
		c.Next() // Proceed to the next handler
	}
}

// RoleAuthMiddleware checks if the user has one of the allowed roles.
func RoleAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			helper.SendError(c, http.StatusForbidden, "Role information not found.")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			helper.SendError(c, http.StatusForbidden, "Invalid role format.")
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		helper.SendError(c, http.StatusForbidden, "Insufficient permissions.")
		c.Abort()
	}
}