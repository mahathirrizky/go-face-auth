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
		AllowOrigins:     []string{"https://4commander.my.id", "https://admin.4commander.my.id", "https://superadmin.4commander.my.id", "https://api.4commander.my.id", "http://admin.localhost:5173", "http://localhost:5173", "http://superadmin.localhost:5173", "http://admin.localhost:8080", "http://localhost:8080", "http://superadmin.localhost:8080", "http://admin.localhost", "http://localhost"},
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
		apiPublic.POST("/login/superadmin", handlers.LoginSuperAdmin)
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
		apiPublic.GET("/confirm-email", handlers.ConfirmEmail)

		// Admin/Internal route for checking subscriptions (should be protected in production)
		adminCompanyHandler := handlers.NewAdminCompanyHandler()
		apiPublic.GET("/check-subscriptions", adminCompanyHandler.CheckAndNotifySubscriptions)
	}

	// Authenticated API routes
	apiAuthenticated := r.Group("/api")
	apiAuthenticated.Use(middleware.AuthMiddleware()) // Apply JWT authentication middleware
	{
		apiAuthenticated.GET("/company-details", handlers.GetCompanyDetails)
		apiAuthenticated.PUT("/company-details", handlers.UpdateCompanyDetails)

		apiAuthenticated.GET("/superadmin/dashboard-summary", handlers.GetSuperAdminDashboardSummary)
		apiAuthenticated.GET("/superadmin/companies", handlers.GetCompanies)
		apiAuthenticated.GET("/superadmin/subscriptions", handlers.GetSubscriptions)
		apiAuthenticated.GET("/superadmin/revenue-summary", handlers.GetRevenueSummary)

		// Subscription Package routes (SuperAdmin)
		apiAuthenticated.POST("/superadmin/subscription-packages", handlers.CreateSubscriptionPackage)
		apiAuthenticated.PUT("/superadmin/subscription-packages/:id", handlers.UpdateSubscriptionPackage)
		apiAuthenticated.DELETE("/superadmin/subscription-packages/:id", handlers.DeleteSubscriptionPackage)
		apiAuthenticated.GET("/superadmin/subscription-packages", handlers.GetSubscriptionPackages)

		apiAuthenticated.GET("/dashboard-summary", func(c *gin.Context) {
			handlers.GetDashboardSummary(hub, c)
		})

		// Attendance routes
		apiAuthenticated.POST("/attendance", func(c *gin.Context) {
			handlers.HandleAttendance(hub, c)
		})
		apiAuthenticated.GET("/attendances", handlers.GetAttendances)
		apiAuthenticated.GET("/employees/:employeeID/attendances", handlers.GetEmployeeAttendanceHistory)
		apiAuthenticated.GET("/employees/:employeeID/attendances/export", handlers.ExportEmployeeAttendanceToExcel)
		apiAuthenticated.GET("/attendances/export", handlers.ExportAllAttendancesToExcel)
		apiAuthenticated.GET("/unaccounted-employees", handlers.GetUnaccountedEmployees)
		apiAuthenticated.GET("/overtime-attendances", handlers.GetOvertimeAttendances)

		// Company routes
		apiAuthenticated.POST("/companies", handlers.CreateCompany)
		apiAuthenticated.GET("/company/:id", handlers.GetCompanyByID)

		// Attendance Location routes (Admin)
		apiAuthenticated.GET("/company/locations", handlers.GetAttendanceLocations)
		apiAuthenticated.POST("/company/locations", handlers.CreateAttendanceLocation)
		apiAuthenticated.PUT("/company/locations/:location_id", handlers.UpdateAttendanceLocation)
		apiAuthenticated.DELETE("/company/locations/:location_id", handlers.DeleteAttendanceLocation)

		// Employee routes
		apiAuthenticated.POST("/employees", handlers.CreateEmployee)
		apiAuthenticated.GET("/employees/:employeeID", handlers.GetEmployeeByID)
		apiAuthenticated.PUT("/employees/:employeeID", handlers.UpdateEmployee)
		apiAuthenticated.DELETE("/employees/:employeeID", handlers.DeleteEmployee)
		apiAuthenticated.GET("/companies/:company_id/employees", handlers.GetEmployeesByCompanyID)
		apiAuthenticated.GET("/companies/:company_id/employees/search", handlers.SearchEmployees)
		apiAuthenticated.GET("/companies/:company_id/employees/pending", handlers.GetPendingEmployees) // New route for pending employees
		apiAuthenticated.POST("/employees/:employee_id/resend-password-email", handlers.ResendPasswordEmail) // New route to resend password email

		// Face Image routes
		apiAuthenticated.POST("/face-images", handlers.UploadFaceImage) // For multipart form data
		apiAuthenticated.GET("/employees/:employeeID/face-images", handlers.GetFaceImagesByEmployeeID)

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

		// Overtime Attendance routes
		apiAuthenticated.POST("/overtime/check-in", func(c *gin.Context) {
			handlers.HandleOvertimeCheckIn(hub, c)
		})
		apiAuthenticated.POST("/overtime/check-out", func(c *gin.Context) {
			handlers.HandleOvertimeCheckOut(hub, c)
		})

		// Broadcast routes
		apiAuthenticated.POST("/broadcast", func(c *gin.Context) {
			handlers.BroadcastMessage(hub, c)
		})
	}

	// WebSocket Face Recognition route
	r.GET("/ws/face-recognition", handlers.FaceRecognitionWebSocketHandler)

	// WebSocket Dashboard Update route
	r.GET("/ws/dashboard", func(c *gin.Context) {
		handlers.ServeWs(hub, c)
	})

	// WebSocket SuperAdmin Dashboard Update route
	r.GET("/ws/superadmin-dashboard", func(c *gin.Context) {
		handlers.SuperAdminDashboardWebSocketHandler(hub, c)
	})

	// WebSocket Employee Notifications route
	r.GET("/ws/employee-notifications", func(c *gin.Context) {
		handlers.EmployeeWebSocketHandler(hub, c)
	})

	// Catch-all route for SPA (Vue.js routing)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}
