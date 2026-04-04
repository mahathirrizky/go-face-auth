package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"go-face-auth/models"

	"time"

	"gorm.io/driver/mysql"
	gorm "gorm.io/gorm"
)

var DB *gorm.DB

// ErrRecordNotFound is a custom error for when a record is not found.
var ErrRecordNotFound = errors.New("record not found")

func InitDB() {


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

	// Memory & Connection Pool Optimization
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)                  // Keep up to 10 connections idle
		sqlDB.SetMaxOpenConns(100)                 // No more than 100 concurrent DB connections
		sqlDB.SetConnMaxLifetime(time.Minute * 30) // Recycle connections after 30 mins
	}

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
		&models.CustomOffer{},
		&models.CustomPackageRequest{},
		&models.DivisionTable{},
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
