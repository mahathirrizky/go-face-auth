package main

import (
	"go-face-auth/config"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/routes"
	"go-face-auth/services"
	"go-face-auth/websocket"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)


func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, assuming environment variables are set.")
	}

	config.LoadMidtransConfig() // Load Midtrans configuration
	config.LoadEmailConfig()    // Load Email configuration

	// Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// --- Python Virtual Environment Setup ---
	venvPath := filepath.Join(".", ".venv")
	pythonExecutable := filepath.Join(venvPath, "bin", "python3") // For Linux/macOS
	// For Windows, it might be: pythonExecutable := filepath.Join(venvPath, "Scripts", "python.exe")

	// Check if venv exists
	if _, err := os.Stat(pythonExecutable); os.IsNotExist(err) {
		log.Println("Python virtual environment not found. Creating and installing dependencies...")

		// Create venv
		cmd := exec.Command("python3", "-m", "venv", venvPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println("Running:", cmd.Args)
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to create virtual environment: %v", err)
		}

		// Install dependencies
		cmd = exec.Command(pythonExecutable, "-m", "pip", "install", "-r", "requirements.txt")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println("Running:", cmd.Args)
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to install Python dependencies: %v", err)
		}
		log.Println("Python virtual environment setup complete.")
	} else {
		log.Println("Python virtual environment found. Skipping setup.")
	}

	// Dynamic Python Server Port
	pythonPort := os.Getenv("PYTHON_SERVER_PORT")
	if pythonPort == "" {
		pythonPort = "5000"
	}

	// --- Start Python face recognition server in a goroutine ---
	var pythonCmd *exec.Cmd
	go func() {
		log.Println("Starting Python face recognition server...")
		pythonCmd = exec.Command(pythonExecutable, "face_recognition_server.py")
		pythonCmd.Dir = "." // Run from current directory

		// Inherit environment variables and ensure specific ones are set
		pythonCmd.Env = os.Environ()
		// Note: os.Environ() already includes PYTHON_SERVER_PORT if it was set in the shell or .env loaded above

		pythonCmd.Stdout = os.Stdout // Redirect Python stdout to Go stdout
		pythonCmd.Stderr = os.Stderr // Redirect Python stderr to Go stderr

		err := pythonCmd.Start()
		if err != nil {
			log.Fatalf("Failed to start Python server: %v", err)
		}
		// Wait for the Python server to exit, if it does
		err = pythonCmd.Wait()
		if err != nil {
			log.Printf("Python face recognition server exited with error: %v", err)
		} else {
			log.Println("Python face recognition server stopped gracefully.")
		}
	}()

	// Ensure Python server is killed on Go app exit
	defer func() {
		if pythonCmd != nil && pythonCmd.Process != nil {
			log.Println("Attempting to kill Python server process...")
			// Try SIGTERM first
			if err := pythonCmd.Process.Signal(syscall.SIGTERM); err != nil {
				log.Printf("Failed to send SIGTERM to Python server: %v", err)
				// Fallback to Kill
				pythonCmd.Process.Kill()
			} else {
				log.Println("Sent SIGTERM to Python server.")
			}
		}
	}()

	// Wait for Python server to be ready using health-check polling
	log.Println("Waiting for Python server to become ready...")
	pythonAddr := "127.0.0.1:" + pythonPort
	pythonReady := false
	for i := 0; i < 30; i++ {
		conn, err := net.DialTimeout("tcp", pythonAddr, 1*time.Second)
		if err == nil {
			conn.Close()
			log.Println("Python server is ready and accepting connections.")
			pythonReady = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !pythonReady {
		log.Println("WARNING: Python server did not become ready within the expected time. Face recognition features may fail.")
	}

	// Initialize database connection
	database.InitDB()
	defer database.CloseDB()

	// Seed initial data (e.g., superadmin and subscription packages)
	database.SeedSuperAdmin()
	database.SeedSubscriptionPackages()

	// Initialize cron scheduler
	c := cron.New(cron.WithLocation(time.UTC)) // Use UTC for cron schedule

	// Initialize all repositories and services needed for the cron job
	companyRepo := repository.NewCompanyRepository(database.DB)
	employeeRepo := repository.NewEmployeeRepository(database.DB)
	attendanceRepo := repository.NewAttendanceRepository(database.DB)
	leaveRequestRepo := repository.NewLeaveRequestRepository(database.DB)
	shiftRepo := repository.NewShiftRepository(database.DB)
	faceImageRepo := repository.NewFaceImageRepository(database.DB)
	attendanceLocationRepo := repository.NewAttendanceLocationRepository(database.DB)
	divisionRepo := repository.NewDivisionRepository(database.DB)
	pythonClient := services.NewPythonClient()

	// Create an instance of the attendance service for the cron job
	cronAttendanceService := services.NewAttendanceService(employeeRepo, companyRepo, attendanceRepo, faceImageRepo, attendanceLocationRepo, leaveRequestRepo, shiftRepo, divisionRepo, pythonClient)

	// Schedule the MarkDailyAbsentees function to run at 03:00, 09:00, 15:00, 21:00 UTC
	_, err := c.AddFunc("0 3,9,15,21 * * *", func() {
		log.Println("Running scheduled MarkDailyAbsentees...")
		if err := cronAttendanceService.MarkDailyAbsentees(); err != nil {
			log.Printf("Error running MarkDailyAbsentees: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule MarkDailyAbsentees: %v", err)
	}

	// Start the cron scheduler in a goroutine
	c.Start()
	log.Println("Cron scheduler started.")

	r := gin.Default()

	routes.SetupRoutes(r, hub) // Pass the hub to SetupRoutes

	// Use PORT environment variable if available, fallback to :8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
