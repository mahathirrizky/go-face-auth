package routes

import (
	"go-face-auth/handlers"
	"go-face-auth/middleware" // Import the middleware package
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors" // Import the cors middleware
)

// NoCache is a middleware function that sets headers to prevent caching.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Next()
}

func SetupRoutes(r *gin.Engine) {
	// Apply the NoCache middleware to all routes
	r.Use(NoCache)

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins for development
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve static files (like index.html, CSS, JS)
	// This will serve index.html when accessing the root URL (e.g., http://localhost:8080/)
	

	r.Static("/assets", "./frontend/dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// Serve employee face images statically
	r.Static("/images/employee_faces", "./images/employee_faces")

	// Public API routes (no authentication required)
	apiPublic := r.Group("/api")
	{
		apiPublic.POST("/login/superuser", handlers.LoginSuperUser)
		apiPublic.POST("/login/admin-company", handlers.LoginAdminCompany)
		apiPublic.POST("/login/employee", handlers.LoginEmployee)
		apiPublic.GET("/subscription-packages", handlers.GetSubscriptionPackages)
		apiPublic.POST("/register-company", handlers.RegisterCompany)
		apiPublic.POST("/payment-confirmation", handlers.HandlePaymentConfirmation)
		apiPublic.POST("/midtrans/create-transaction", handlers.CreateMidtransTransaction)
		apiPublic.GET("/invoices/:order_id", handlers.GetInvoiceByOrderID)
		apiPublic.POST("/forgot-password", handlers.ForgotPassword)
		apiPublic.POST("/reset-password", handlers.ResetPassword)
	}

	// Authenticated API routes
	apiAuthenticated := r.Group("/api")
	apiAuthenticated.Use(middleware.AuthMiddleware()) // Apply JWT authentication middleware
	{
		// Attendance routes
		apiAuthenticated.POST("/attendance", handlers.HandleAttendance)

		// Company routes
		apiAuthenticated.POST("/companies", handlers.CreateCompany)
		apiAuthenticated.GET("/company/:id", handlers.GetCompanyByID)

		// Employee routes
		apiAuthenticated.POST("/employees", handlers.CreateEmployee)
		apiAuthenticated.GET("/employee/:id", handlers.GetEmployeeByID)
		apiAuthenticated.GET("/companies/:company_id/employees", handlers.GetEmployeesByCompanyID)

		// Face Image routes
		apiAuthenticated.POST("/face-images", handlers.UploadFaceImage) // For multipart form data
		apiAuthenticated.GET("/employees/:employee_id/face-images", handlers.GetFaceImagesByEmployeeID)
		
	}

	// WebSocket Face Recognition route
	r.GET("/ws/face-recognition", handlers.FaceRecognitionWebSocketHandler)

	// Catch-all route for SPA (Vue.js routing)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}