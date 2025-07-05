package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-face-auth/helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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

		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		if len(jwtSecret) == 0 {
			// Fallback for development if env var is not set
			jwtSecret = []byte("supersecretjwtkeythatshouldbechangedinproduction")
			log.Println("WARNING: JWT_SECRET environment variable not set for AuthMiddleware. Using default secret.")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what we expect: HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			log.Printf("AuthMiddleware: Token parsing error: %v", err)
			helper.SendError(c, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Printf("AuthMiddleware: Token valid. Claims: ID=%v, Role=%v, CompanyID=%v", claims["id"], claims["role"], claims["companyID"])
			// Set claims in gin context
			c.Set("id", claims["id"])
			c.Set("role", claims["role"])
			c.Set("companyID", claims["companyID"]) // Set companyID in context
			c.Next() // Proceed to the next handler
		} else {
			log.Printf("AuthMiddleware: Invalid token claims or not valid. Token valid: %t, Claims ok: %t", token.Valid, ok)
			helper.SendError(c, http.StatusUnauthorized, "Invalid token claims.")
			c.Abort()
			return
		}
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