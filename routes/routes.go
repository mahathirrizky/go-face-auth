package routes

import (
	"go-face-auth/handlers"

	"github.com/gin-gonic/gin"
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

	// Serve static files (like index.html, CSS, JS)
	// This will serve index.html when accessing the root URL (e.g., http://localhost:8080/)
	

	r.Static("/assets", "./frontend/dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// Serve employee face images statically
	r.Static("/images/employee_faces", "./images/employee_faces")

	// API routes
	api := r.Group("/api")
	{
		// Attendance routes
		api.POST("/attendance", handlers.HandleAttendance)

		// Company routes
		api.POST("/companies", handlers.CreateCompany)
		api.GET("/company/:id", handlers.GetCompanyByID)

		// Employee routes
		api.POST("/employees", handlers.CreateEmployee)
		api.GET("/employee/:id", handlers.GetEmployeeByID)
		api.GET("/companies/:company_id/employees", handlers.GetEmployeesByCompanyID)

		// Face Image routes
		api.POST("/face-images", handlers.UploadFaceImage) // For multipart form data
		api.GET("/employees/:employee_id/face-images", handlers.GetFaceImagesByEmployeeID)
	}

	

	// WebSocket Face Recognition route
	r.GET("/ws/face-recognition", handlers.FaceRecognitionWebSocketHandler)

	// Catch-all route for SPA (Vue.js routing)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}