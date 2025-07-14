package database

import (
	"fmt"
	"log"
	"os"

	"go-face-auth/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	gorm "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database with GORM: %v", err)
	}

	DB = db

	log.Println("Successfully connected to MySQL database with GORM!")

	// AutoMigrate will create/update tables based on your models
	log.Println("Running GORM AutoMigrate...")
	err = DB.AutoMigrate(
		&models.CompaniesTable{},
		&models.EmployeesTable{},
		&models.FaceImagesTable{},
		&models.AttendancesTable{},
		&models.AdminCompaniesTable{},
		&models.SuperAdminTable{},
		&models.SubscriptionPackageTable{},
		&models.InvoiceTable{},
		&models.PasswordResetTokenTable{},
		&models.ShiftsTable{},
		&models.LeaveRequest{},
		&models.AttendanceLocation{},
		&models.BroadcastMessage{},
		&models.EmployeeBroadcastRead{},
	)
	if err != nil {
		log.Fatalf("Error running GORM AutoMigrate: %v", err)
	}
	log.Println("GORM AutoMigrate completed.")
}

func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error getting underlying sql.DB from GORM: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
		log.Println("Database connection closed.")
	}
}
