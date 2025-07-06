package main

import (
	"go-face-auth/config"
	"go-face-auth/database"
	"go-face-auth/routes"
	"go-face-auth/websocket"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func killProcessOnPort(port string) {
	cmd := exec.Command("lsof", "-t", "-i", ":"+port)
	output, err := cmd.Output()
	if err != nil {
		// No process found or lsof error
		return
	}

	pids := strings.TrimSpace(string(output))
	if pids == "" {
		return
	}

	log.Printf("Found process(es) on port %s: %s", port, pids)
	// Kill the process(es)
	killCmd := exec.Command("kill", strings.Fields(pids)...)
	killCmd.Stdout = os.Stdout
	killCmd.Stderr = os.Stderr
	if err := killCmd.Run(); err != nil {
		log.Printf("Failed to kill process(es) on port %s: %v", port, err)
	} else {
		log.Printf("Successfully killed process(es) on port %s.", port)
	}
}

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

	// Ensure port 5000 is free before starting Python server
	killProcessOnPort("5000")

	// --- Start Python face recognition server in a goroutine ---
	var pythonCmd *exec.Cmd
	go func() {
		log.Println("Starting Python face recognition server...")
		pythonCmd = exec.Command(pythonExecutable, "face_recognition_server.py")
		pythonCmd.Dir = "." // Run from current directory
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
			if err := pythonCmd.Process.Kill(); err != nil {
				log.Printf("Failed to kill Python server process: %v", err)
			} else {
				log.Println("Python server process killed.")
			}
		}
	}()

	// Give Python server a moment to start up
	time.Sleep(10 * time.Second)

	// Initialize database connection
	database.InitDB()
	defer database.CloseDB()

	// Seed initial data (e.g., superuser and subscription packages)
	database.SeedSuperUser()
	database.SeedSubscriptionPackages()

	r := gin.Default()

	routes.SetupRoutes(r, hub) // Pass the hub to SetupRoutes

	r.Run("localhost:8080") // listen and serve on localhost:8080
}
