package routes

import (
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/handlers"
	"go-face-auth/helper"
	"go-face-auth/middleware"
	"go-face-auth/services"
	"go-face-auth/websocket"
	"log"
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
		AllowOrigins:     []string{"https://4commander.my.id", "https://admin.4commander.my.id", "https://superadmin.4commander.my.id", "https://api.4commander.my.id", "http://admin.localhost:5173", "http://localhost:5173", "http://superadmin.localhost:5173", "http://admin.localhost:8080", "http://localhost:8080", "http://superadmin.localhost:8080", "http://admin.localhost", "http://localhost", "https://d3eeb2f5afa3.ngrok-free.app", "https://admin.d3eeb2f5afa3.ngrok-free.app", "https://superadmin.d3eeb2f5afa3.ngrok-free.app"},
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

	// Background worker for Superadmin Dashboard Updates
	go func() {
		for range hub.TriggerSuperAdminUpdate {
			summary, err := superAdminService.GetSuperAdminDashboardSummary()
			if err != nil {
				log.Printf("Failed to generate superadmin dashboard summary: %v", err)
				continue
			}
			hub.SendSuperAdminDashboardUpdate(summary)
		}
	}()

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
		apiPublic.GET("/offer/:token", customOfferHandler.HandleGetCustomOfferByToken)

		apiPublic.GET("/check-subscriptions", adminCompanyHandler.CheckAndNotifySubscriptions)
	}

	// Authenticated API routes
	apiAuthenticated := r.Group("/api")
	apiAuthenticated.Use(middleware.AuthMiddleware())

	// Admin (including superadmins) routes
	adminRoutes := apiAuthenticated.Group("/")
	adminRoutes.Use(middleware.RoleAuthMiddleware("admin", "superadmin"))
	{
		adminRoutes.GET("/company-details", companyHandler.GetCompanyDetails)
		adminRoutes.PUT("/company-details", companyHandler.UpdateCompanyDetails)
		adminRoutes.GET("/company/:companyId/subscription-status", companyHandler.GetCompanySubscriptionStatus)
		adminRoutes.PUT("/admin/change-password", adminCompanyHandler.ChangeAdminPassword)

		// Dashboard
		adminRoutes.GET("/dashboard-summary", func(c *gin.Context) {
			adminCompanyHandler.GetDashboardSummary(hub, c)
		})

		// Invoice History routes (Admin)
		adminRoutes.GET("/invoices", paymentHandler.GetCompanyInvoices)
		adminRoutes.GET("/invoices/:order_id/download", paymentHandler.DownloadInvoicePDF)

		// Company routes
		adminRoutes.POST("/companies", companyHandler.CreateCompany)
		adminRoutes.GET("/company/:companyId", companyHandler.GetCompanyByID)

		// Attendance Location routes (Admin)
		adminRoutes.GET("/company/locations", locationHandler.GetAttendanceLocations)
		adminRoutes.POST("/company/locations", locationHandler.CreateAttendanceLocation)
		adminRoutes.PUT("/company/locations/:location_id", locationHandler.UpdateAttendanceLocation)
		adminRoutes.DELETE("/company/locations/:location_id", locationHandler.DeleteAttendanceLocation)

		// Employee management routes
		adminRoutes.POST("/employees", employeeHandler.CreateEmployee)
		adminRoutes.GET("/employees/:employeeID", employeeHandler.GetEmployeeByID)
		adminRoutes.PUT("/employees/:employeeID", employeeHandler.UpdateEmployee)
		adminRoutes.DELETE("/employees/:employeeID", employeeHandler.DeleteEmployee)
		adminRoutes.GET("/companies/:company_id/employees", employeeHandler.GetEmployeesByCompanyID)
		adminRoutes.GET("/companies/:company_id/employees/search", employeeHandler.SearchEmployees)
		adminRoutes.GET("/companies/:company_id/employees/pending", employeeHandler.GetPendingEmployees)
		adminRoutes.POST("/employees/:employee_id/resend-password-email", employeeHandler.ResendPasswordEmail)

		// Bulk Employee Import
		adminRoutes.GET("/employees/template", employeeHandler.GenerateEmployeeTemplate)
		adminRoutes.POST("/employees/bulk", employeeHandler.BulkCreateEmployees)

		// Face Image routes (admin trigger)
		adminRoutes.POST("/employee/register-face", employeeHandler.UploadFaceImage)
		adminRoutes.GET("/employees/:employeeID/face-images", employeeHandler.GetFaceImagesByEmployeeID)

		// Shift routes
		adminRoutes.POST("/shifts", shiftHandler.CreateShift)
		adminRoutes.GET("/shifts", shiftHandler.GetShiftsByCompany)
		adminRoutes.PUT("/shifts/:id", shiftHandler.UpdateShift)
		adminRoutes.DELETE("/shifts/:id", shiftHandler.DeleteShift)
		adminRoutes.POST("/shifts/set-default", shiftHandler.SetDefaultShift)

		// Division routes
		adminRoutes.POST("/admin/divisions", divisionHandler.CreateDivision)
		adminRoutes.GET("/admin/divisions", divisionHandler.GetDivisions)
		adminRoutes.GET("/admin/divisions/:id", divisionHandler.GetDivisionByID)
		adminRoutes.PUT("/admin/divisions/:id", divisionHandler.UpdateDivision)
		adminRoutes.DELETE("/admin/divisions/:id", divisionHandler.DeleteDivision)

		// Attendance routes
		adminRoutes.POST("/attendance", func(c *gin.Context) {
			attendanceHandler.HandleAttendance(hub, c)
		})
		adminRoutes.GET("/attendances", attendanceHandler.GetAttendances)
		adminRoutes.GET("/employees/:employeeID/attendances", attendanceHandler.GetEmployeeAttendanceHistory)
		adminRoutes.GET("/employees/:employeeID/attendances/export", attendanceHandler.ExportEmployeeAttendanceToExcel)
		adminRoutes.GET("/attendances/export", attendanceHandler.ExportAllAttendancesToExcel)
		adminRoutes.GET("/attendances/unaccounted", attendanceHandler.GetUnaccountedEmployees)
		adminRoutes.GET("/attendances/unaccounted/export", attendanceHandler.ExportUnaccountedToExcel)
		adminRoutes.GET("/attendances/overtime", attendanceHandler.GetOvertimeAttendances)
		adminRoutes.GET("/attendances/overtime/export", attendanceHandler.ExportOvertimeToExcel)
		adminRoutes.POST("/attendances/correction", attendanceHandler.CorrectAttendance)

		// Leave Request routes (Admin)
		adminRoutes.GET("/company-leave-requests", leaveRequestHandler.GetAllCompanyLeaveRequests)
		adminRoutes.GET("/company-leave-requests/export", leaveRequestHandler.ExportCompanyLeaveRequestsToExcel)
		adminRoutes.PUT("/leave-requests/:id/review", leaveRequestHandler.ReviewLeaveRequest(hub))
		adminRoutes.PUT("/leave-requests/:id/admin-cancel", leaveRequestHandler.AdminCancelApprovedLeaveHandler)

		// Overtime Attendance routes
		adminRoutes.POST("/overtime/check-in", func(c *gin.Context) {
			attendanceHandler.HandleOvertimeCheckIn(hub, c)
		})
		adminRoutes.POST("/overtime/check-out", func(c *gin.Context) {
			attendanceHandler.HandleOvertimeCheckOut(hub, c)
		})

		// Broadcast routes
		adminRoutes.POST("/broadcasts", func(c *gin.Context) {
			broadcastHandler.BroadcastMessage(hub, c)
		})
		adminRoutes.GET("/broadcasts", broadcastHandler.GetBroadcasts)
		adminRoutes.POST("/broadcasts/:id/read", broadcastHandler.MarkBroadcastAsRead)
		adminRoutes.POST("/custom-package-requests", customPackageRequestHandler.HandleCustomPackageRequest(hub))
	}

	// Superadmin-specific routes
	superAdminRoutes := apiAuthenticated.Group("/superadmin")
	superAdminRoutes.Use(middleware.RoleAuthMiddleware("superadmin"))
	{
		superAdminRoutes.GET("/dashboard-summary", superAdminHandler.GetSuperAdminDashboardSummary)
		superAdminRoutes.GET("/companies", superAdminHandler.GetCompanies)
		superAdminRoutes.GET("/subscriptions", superAdminHandler.GetSubscriptions)
		superAdminRoutes.GET("/revenue-summary", superAdminHandler.GetRevenueSummary)
		superAdminRoutes.POST("/subscription-packages", subscriptionPackageHandler.CreateSubscriptionPackage)
		superAdminRoutes.PUT("/subscription-packages/:id", subscriptionPackageHandler.UpdateSubscriptionPackage)
		superAdminRoutes.DELETE("/subscription-packages/:id", subscriptionPackageHandler.DeleteSubscriptionPackage)
		superAdminRoutes.GET("/subscription-packages", subscriptionPackageHandler.GetSubscriptionPackages)
		superAdminRoutes.POST("/custom-offers", customOfferHandler.HandleCreateCustomOffer)
		superAdminRoutes.GET("/custom-package-requests", superAdminHandler.GetCustomPackageRequests)
		superAdminRoutes.PUT("/custom-package-requests/:id/:status", superAdminHandler.UpdateCustomPackageRequestStatus)
	}

	// Employee-specific routes (also accessible by superadmin/admin if desired via role middleware)
	employeeRoutes := apiAuthenticated.Group("/employee")
	employeeRoutes.Use(middleware.RoleAuthMiddleware("employee", "admin", "superadmin"))
	{
		employeeRoutes.GET("/profile", employeeHandler.GetEmployeeProfile)
		employeeRoutes.PUT("/profile", employeeHandler.UpdateEmployeeProfile)
		employeeRoutes.PUT("/change-password", employeeHandler.ChangeEmployeePassword)
		employeeRoutes.GET("/dashboard-summary", employeeHandler.GetEmployeeDashboardSummary)
		// Allow employees to register their own face image
		employeeRoutes.POST("/register-face", employeeHandler.UploadFaceImage)
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
