package routes

import (
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/handlers"
	"go-face-auth/helper"
	"go-face-auth/middleware"
	"go-face-auth/services"
	"go-face-auth/websocket"
	"net/http"
	"os"
	"path/filepath"
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
		AllowOrigins:     []string{"https://4commander.my.id", "https://admin.4commander.my.id", "https://superadmin.4commander.my.id", "https://api.4commander.my.id", "http://admin.localhost:5173", "http://localhost:5173", "http://superadmin.localhost:5173", "http://admin.localhost:8080", "http://localhost:8080", "http://superadmin.localhost:8080", "http://admin.localhost", "http://localhost","https://d3eeb2f5afa3.ngrok-free.app","https://admin.d3eeb2f5afa3.ngrok-free.app","https://superadmin.d3eeb2f5afa3.ngrok-free.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// =================================================================
	// DEPENDENCY INJECTION SETUP
	// =================================================================
	db := database.DB // Your gorm.DB instance

	// Repositories
	adminCompanyRepo := repository.NewAdminCompanyRepository(db)
	attendanceLocationRepo := repository.NewAttendanceLocationRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	broadcastRepo := repository.NewBroadcastRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	customOfferRepo := repository.NewCustomOfferRepository(db)
	customPackageRequestRepo := repository.NewCustomPackageRequestRepository(db)
	divisionRepo := repository.NewDivisionRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	faceImageRepo := repository.NewFaceImageRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)
	leaveRequestRepo := repository.NewLeaveRequestRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	shiftRepo := repository.NewShiftRepository(db)
	subscriptionPackageRepo := repository.NewSubscriptionPackageRepository(db)
	superAdminRepo := repository.NewSuperAdminRepository(db)

		pythonClient := services.NewPythonClient()

	// Services
	authService := services.NewAuthService(superAdminRepo, adminCompanyRepo, employeeRepo, attendanceLocationRepo)
	adminCompanyService := services.NewAdminCompanyService(adminCompanyRepo, companyRepo, employeeRepo, attendanceRepo, leaveRequestRepo)
	attendanceService := services.NewAttendanceService(employeeRepo, companyRepo, attendanceRepo, faceImageRepo, attendanceLocationRepo, leaveRequestRepo, shiftRepo, divisionRepo, pythonClient)
	broadcastService := services.NewBroadcastService(broadcastRepo)
	companyService := services.NewCompanyService(companyRepo, adminCompanyRepo, subscriptionPackageRepo, shiftRepo)
	customOfferService := services.NewCustomOfferService(customOfferRepo)
	customPackageRequestService := services.NewCustomPackageRequestService(companyRepo, adminCompanyRepo, customPackageRequestRepo)
	divisionService := services.NewDivisionService(divisionRepo, shiftRepo, attendanceLocationRepo)
	employeeService := services.NewEmployeeService(employeeRepo, companyRepo, shiftRepo, passwordResetRepo, faceImageRepo, attendanceRepo, leaveRequestRepo, attendanceLocationRepo)
	initialPasswordSetupService := services.NewInitialPasswordSetupService(passwordResetRepo, employeeRepo)
	leaveRequestService := services.NewLeaveRequestService(employeeRepo, leaveRequestRepo, adminCompanyRepo)
	locationService := services.NewLocationService(companyRepo, attendanceLocationRepo)
	passwordResetService := services.NewPasswordResetService(adminCompanyRepo, employeeRepo, passwordResetRepo)
	paymentService := services.NewPaymentService(invoiceRepo, companyRepo, subscriptionPackageRepo, customOfferRepo, adminCompanyRepo, helper.NewPDFGenerator())
	shiftService := services.NewShiftService(shiftRepo, companyRepo)
	subscriptionPackageService := services.NewSubscriptionPackageService(subscriptionPackageRepo)
	superAdminService := services.NewSuperAdminService(companyRepo, invoiceRepo, customPackageRequestRepo, superAdminRepo)

	// Handlers
	adminCompanyHandler := handlers.NewAdminCompanyHandler(adminCompanyService)
	// adminHandler is removed
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService, adminCompanyService) // Use adminCompanyService for dashboard summary
	authHandler := handlers.NewAuthHandler(authService)
	broadcastHandler := handlers.NewBroadcastHandler(broadcastService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	customOfferHandler := handlers.NewCustomOfferHandler(customOfferService)
	customPackageRequestHandler := handlers.NewCustomPackageRequestHandler(customPackageRequestService)
	divisionHandler := handlers.NewDivisionHandler(divisionService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService, shiftService)
	initialPasswordSetupHandler := handlers.NewInitialPasswordSetupHandler(initialPasswordSetupService)
	leaveRequestHandler := handlers.NewLeaveRequestHandler(leaveRequestService, adminCompanyService) // Use adminCompanyService for dashboard summary
	locationHandler := handlers.NewLocationHandler(locationService)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	shiftHandler := handlers.NewShiftHandler(shiftService)
	subscriptionPackageHandler := handlers.NewSubscriptionPackageHandler(subscriptionPackageService)
	superAdminHandler := handlers.NewSuperAdminHandler(superAdminService)
	// Employee WebSocket Handler (no service dependencies for now)
	employeeWebSocketHandler := handlers.NewEmployeeWebSocketHandler()

	// =================================================================

	// Serve static files (like index.html, CSS, JS)
	// This will serve index.html when accessing the root URL (e.g., http://localhost:8080/)

	r.Static("/assets", "./frontend/dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// Serve employee face images statically
	storageBaseDir := os.Getenv("STORAGE_BASE_PATH")
	if storageBaseDir == "" {
		storageBaseDir = "/tmp/go_face_auth_data" // Fallback for development/testing
	}
	r.Static("/images/employee_faces", filepath.Join(storageBaseDir, "employee_faces"))

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
		apiPublic.POST("/login/superadmin", authHandler.LoginSuperAdmin)
		apiPublic.POST("/login/admin-company", authHandler.LoginAdminCompany)
		apiPublic.POST("/login/employee", authHandler.LoginEmployee)
		apiPublic.GET("/subscription-packages", subscriptionPackageHandler.GetSubscriptionPackages)
		apiPublic.POST("/register-company", companyHandler.RegisterCompany(hub))
		apiPublic.POST("/payment-confirmation", paymentHandler.HandlePaymentConfirmation(hub))
		apiPublic.POST("/midtrans/create-transaction", paymentHandler.CreateMidtransTransaction)
		apiPublic.GET("/invoices/:order_id", paymentHandler.GetInvoiceByOrderID)
		apiPublic.POST("/forgot-password", passwordResetHandler.ForgotPassword)
		apiPublic.POST("/forgot-password-employee", passwordResetHandler.ForgotPasswordEmployee)
		apiPublic.POST("/reset-password", passwordResetHandler.ResetPassword)
		apiPublic.POST("/initial-password-setup", initialPasswordSetupHandler.InitialPasswordSetup) // New route for initial password setup
		apiPublic.GET("/confirm-email", companyHandler.ConfirmEmail)

		apiPublic.GET("/check-subscriptions", adminCompanyHandler.CheckAndNotifySubscriptions)
	}

	// Authenticated API routes
	apiAuthenticated := r.Group("/api")
	apiAuthenticated.Use(middleware.AuthMiddleware()) // Apply JWT authentication middleware
	{
		apiAuthenticated.GET("/company-details", companyHandler.GetCompanyDetails)
		apiAuthenticated.PUT("/company-details", companyHandler.UpdateCompanyDetails)
		apiAuthenticated.GET("/company/:companyId/subscription-status", companyHandler.GetCompanySubscriptionStatus)
		apiAuthenticated.PUT("/admin/change-password", adminCompanyHandler.ChangeAdminPassword)

		apiAuthenticated.GET("/superadmin/dashboard-summary", superAdminHandler.GetSuperAdminDashboardSummary)
		apiAuthenticated.GET("/superadmin/companies", superAdminHandler.GetCompanies)
		apiAuthenticated.GET("/superadmin/subscriptions", superAdminHandler.GetSubscriptions)
		apiAuthenticated.GET("/superadmin/revenue-summary", superAdminHandler.GetRevenueSummary)

		// Invoice History routes (Admin)
		apiAuthenticated.GET("/invoices", paymentHandler.GetCompanyInvoices)
		apiAuthenticated.GET("/invoices/:order_id/download", paymentHandler.DownloadInvoicePDF)

		// Subscription Package routes (SuperAdmin)
		apiAuthenticated.POST("/superadmin/subscription-packages", subscriptionPackageHandler.CreateSubscriptionPackage)
		apiAuthenticated.PUT("/superadmin/subscription-packages/:id", subscriptionPackageHandler.UpdateSubscriptionPackage)
		apiAuthenticated.DELETE("/superadmin/subscription-packages/:id", subscriptionPackageHandler.DeleteSubscriptionPackage)
		apiAuthenticated.GET("/superadmin/subscription-packages", subscriptionPackageHandler.GetSubscriptionPackages)
		apiAuthenticated.POST("/superadmin/custom-offers", customOfferHandler.HandleCreateCustomOffer) // New route for superadmin to create custom offers

		apiAuthenticated.GET("/dashboard-summary", func(c *gin.Context) {
			adminCompanyHandler.GetDashboardSummary(hub, c)
		})

		// Attendance routes
		apiAuthenticated.POST("/attendance", func(c *gin.Context) {
			attendanceHandler.HandleAttendance(hub, c)
		})
		apiAuthenticated.GET("/attendances", attendanceHandler.GetAttendances)
		apiAuthenticated.GET("/employees/:employeeID/attendances", attendanceHandler.GetEmployeeAttendanceHistory)
		apiAuthenticated.GET("/employees/:employeeID/attendances/export", attendanceHandler.ExportEmployeeAttendanceToExcel)
		apiAuthenticated.GET("/attendances/export", attendanceHandler.ExportAllAttendancesToExcel)
		apiAuthenticated.GET("/attendances/unaccounted", attendanceHandler.GetUnaccountedEmployees)
		apiAuthenticated.GET("/attendances/unaccounted/export", attendanceHandler.ExportUnaccountedToExcel)
		apiAuthenticated.GET("/attendances/overtime", attendanceHandler.GetOvertimeAttendances)
		apiAuthenticated.GET("/attendances/overtime/export", attendanceHandler.ExportOvertimeToExcel)
		apiAuthenticated.POST("/attendances/correction", attendanceHandler.CorrectAttendance) // New route for manual correction by admin

		// Company routes
		apiAuthenticated.POST("/companies", companyHandler.CreateCompany)
		apiAuthenticated.GET("/company/:companyId", companyHandler.GetCompanyByID)

		// Attendance Location routes (Admin)
		apiAuthenticated.GET("/company/locations", locationHandler.GetAttendanceLocations)
		apiAuthenticated.POST("/company/locations", locationHandler.CreateAttendanceLocation)
		apiAuthenticated.PUT("/company/locations/:location_id", locationHandler.UpdateAttendanceLocation)
		apiAuthenticated.DELETE("/company/locations/:location_id", locationHandler.DeleteAttendanceLocation)

		// Employee routes
		apiAuthenticated.POST("/employees", employeeHandler.CreateEmployee)
		apiAuthenticated.GET("/employees/:employeeID", employeeHandler.GetEmployeeByID)
		apiAuthenticated.PUT("/employees/:employeeID", employeeHandler.UpdateEmployee)
		apiAuthenticated.DELETE("/employees/:employeeID", employeeHandler.DeleteEmployee)
		apiAuthenticated.GET("/companies/:company_id/employees", employeeHandler.GetEmployeesByCompanyID)
		apiAuthenticated.GET("/companies/:company_id/employees/search", employeeHandler.SearchEmployees)
		apiAuthenticated.GET("/companies/:company_id/employees/pending", employeeHandler.GetPendingEmployees)       // New route for pending employees
		apiAuthenticated.POST("/employees/:employee_id/resend-password-email", employeeHandler.ResendPasswordEmail) // New route to resend password email

		// Bulk Employee Import
		apiAuthenticated.GET("/employees/template", employeeHandler.GenerateEmployeeTemplate)
		apiAuthenticated.POST("/employees/bulk", employeeHandler.BulkCreateEmployees)

		// Employee Profile route
		apiAuthenticated.GET("/employee/profile", employeeHandler.GetEmployeeProfile)
		apiAuthenticated.PUT("/employee/profile", employeeHandler.UpdateEmployeeProfile)
		apiAuthenticated.PUT("/employee/change-password", employeeHandler.ChangeEmployeePassword)

		// Employee Dashboard Summary route
		apiAuthenticated.GET("/employee/dashboard-summary", employeeHandler.GetEmployeeDashboardSummary)

		// Face Image routes
		apiAuthenticated.POST("/employee/register-face", employeeHandler.UploadFaceImage) // For multipart form data
		apiAuthenticated.GET("/employees/:employeeID/face-images", employeeHandler.GetFaceImagesByEmployeeID)

		// Shift routes
		apiAuthenticated.POST("/shifts", shiftHandler.CreateShift)
		apiAuthenticated.GET("/shifts", shiftHandler.GetShiftsByCompany)
		apiAuthenticated.PUT("/shifts/:id", shiftHandler.UpdateShift)
		apiAuthenticated.DELETE("/shifts/:id", shiftHandler.DeleteShift)
		apiAuthenticated.POST("/shifts/set-default", shiftHandler.SetDefaultShift)

		// Division routes
		apiAuthenticated.POST("/admin/divisions", divisionHandler.CreateDivision)
		apiAuthenticated.GET("/admin/divisions", divisionHandler.GetDivisions)
		apiAuthenticated.GET("/admin/divisions/:id", divisionHandler.GetDivisionByID)
		apiAuthenticated.PUT("/admin/divisions/:id", divisionHandler.UpdateDivision)
		apiAuthenticated.DELETE("/admin/divisions/:id", divisionHandler.DeleteDivision)

		// Leave Request routes (Employee)
		apiAuthenticated.POST("/leave-requests", leaveRequestHandler.ApplyLeave)
		apiAuthenticated.GET("/my-leave-requests", leaveRequestHandler.GetMyLeaveRequests)
		apiAuthenticated.DELETE("/leave-requests/:id", leaveRequestHandler.CancelLeaveRequest)

		// Leave Request routes (Admin)
		apiAuthenticated.GET("/company-leave-requests", leaveRequestHandler.GetAllCompanyLeaveRequests)
		apiAuthenticated.GET("/company-leave-requests/export", leaveRequestHandler.ExportCompanyLeaveRequestsToExcel)
		apiAuthenticated.PUT("/leave-requests/:id/review", leaveRequestHandler.ReviewLeaveRequest(hub))
		apiAuthenticated.PUT("/leave-requests/:id/admin-cancel", leaveRequestHandler.AdminCancelApprovedLeaveHandler)

		// Overtime Attendance routes
		apiAuthenticated.POST("/overtime/check-in", func(c *gin.Context) {
			attendanceHandler.HandleOvertimeCheckIn(hub, c)
		})
		apiAuthenticated.POST("/overtime/check-out", func(c *gin.Context) {
			attendanceHandler.HandleOvertimeCheckOut(hub, c)
		})

		// Broadcast routes
		apiAuthenticated.POST("/broadcasts", func(c *gin.Context) {
			broadcastHandler.BroadcastMessage(hub, c)
		})
		apiAuthenticated.GET("/broadcasts", broadcastHandler.GetBroadcasts)
		apiAuthenticated.POST("/broadcasts/:id/read", broadcastHandler.MarkBroadcastAsRead)
		apiAuthenticated.POST("/custom-package-requests", customPackageRequestHandler.HandleCustomPackageRequest(hub)) // New route for custom package requests

		// SuperAdmin Custom Package Request routes
		apiAuthenticated.GET("/superadmin/custom-package-requests", superAdminHandler.GetCustomPackageRequests)
		apiAuthenticated.PUT("/superadmin/custom-package-requests/:id/:status", superAdminHandler.UpdateCustomPackageRequestStatus)
		apiAuthenticated.GET("/offer/:token", customOfferHandler.HandleGetCustomOfferByToken) // Moved to authenticated route
	}

	// WebSocket Dashboard Update route
	r.GET("/ws/dashboard", func(c *gin.Context) {
		handlers.ServeWs(hub, c)
	})

	// WebSocket SuperAdmin Dashboard Update route
	r.GET("/ws/superadmin-dashboard", func(c *gin.Context) {
		superAdminHandler.SuperAdminDashboardWebSocketHandler(hub, c)
	})

	// WebSocket Employee Notifications route
	r.GET("/ws/employee-notifications", func(c *gin.Context) {
		employeeWebSocketHandler.HandleEmployeeWebSocket(hub, c)
	})

	// Catch-all route for SPA (Vue.js routing)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}