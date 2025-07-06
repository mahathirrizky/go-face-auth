package routes

import (
	"go-face-auth/handlers"
	"go-face-auth/middleware"
	"go-face-auth/websocket"
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NoCache is a middleware function that sets headers to prevent caching.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Next()
}

func SetupRoutes(r *gin.Engine, hub *websocket.Hub) {
	// Apply the NoCache middleware to all routes
	r.Use(NoCache)

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://admin.localhost:5173", "http://localhost:5173"}, // Allow specific origins for development
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

	// Create a rate limiter middleware
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 10,
	})
	limiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests. Please try again later.",
			})
		},
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	})

	// Public API routes (no authentication required)
	apiPublic := r.Group("/api")
	apiPublic.Use(limiter) // Apply rate limiting to all public routes
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
		apiPublic.POST("/forgot-password-employee", handlers.ForgotPasswordEmployee)
		apiPublic.POST("/reset-password", handlers.ResetPassword)
	}

	// Authenticated API routes
	apiAuthenticated := r.Group("/api")
	apiAuthenticated.Use(middleware.AuthMiddleware()) // Apply JWT authentication middleware
	{
		apiAuthenticated.GET("/company-details", handlers.GetCompanyDetails)

		apiAuthenticated.GET("/dashboard-summary", func(c *gin.Context) {
			handlers.GetDashboardSummary(hub, c)
		})

		// Attendance routes
		apiAuthenticated.POST("/attendance", func(c *gin.Context) {
			handlers.HandleAttendance(hub, c)
		})
		apiAuthenticated.GET("/attendances", handlers.GetAttendances)

		// Company routes
		apiAuthenticated.POST("/companies", handlers.CreateCompany)
		apiAuthenticated.GET("/company/:id", handlers.GetCompanyByID)

		// Employee routes
		apiAuthenticated.POST("/employees", handlers.CreateEmployee)
		apiAuthenticated.GET("/employee/:id", handlers.GetEmployeeByID)
		apiAuthenticated.GET("/companies/:company_id/employees", handlers.GetEmployeesByCompanyID)
		apiAuthenticated.GET("/companies/:company_id/employees/search", handlers.SearchEmployees)

		// Face Image routes
		apiAuthenticated.POST("/face-images", handlers.UploadFaceImage) // For multipart form data
		apiAuthenticated.GET("/employees/:employee_id/face-images", handlers.GetFaceImagesByEmployeeID)

		// Shift routes
		apiAuthenticated.POST("/shifts", handlers.CreateShift)
		apiAuthenticated.GET("/shifts", handlers.GetShiftsByCompany)
		apiAuthenticated.PUT("/shifts/:id", handlers.UpdateShift)
		apiAuthenticated.DELETE("/shifts/:id", handlers.DeleteShift)

		// Leave Request routes (Employee)
		apiAuthenticated.POST("/leave-requests", handlers.ApplyLeave)
		apiAuthenticated.GET("/my-leave-requests", handlers.GetMyLeaveRequests)

		// Leave Request routes (Admin)
		apiAuthenticated.GET("/company-leave-requests", handlers.GetAllCompanyLeaveRequests)
		apiAuthenticated.PUT("/leave-requests/:id/review", handlers.ReviewLeaveRequest)
	}

	// WebSocket Face Recognition route
	r.GET("/ws/face-recognition", handlers.FaceRecognitionWebSocketHandler)

	// WebSocket Dashboard Update route
	r.GET("/ws/dashboard", func(c *gin.Context) {
		handlers.ServeWs(hub, c)
	})

	// Catch-all route for SPA (Vue.js routing)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}