package database

import (
	"go-face-auth/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedSuperUser creates a default superuser if one does not already exist.
func SeedSuperUser() {
	// Check if a superuser already exists
	var superUser models.SuperUserTable
	result := DB.Where("email = ?", "superuser@example.com").First(&superUser)

	if result.Error == nil {
		log.Println("Superuser 'superuser@example.com' already exists. Skipping seeding.")
		return
	}

	// If not found, create a new superuser
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash superuser password: %v", err)
	}

	newSuperUser := models.SuperUserTable{
		Email:    "superuser@example.com",
		Password: string(hashedPassword),
		Role:     "super_admin",
	}

	if err := DB.Create(&newSuperUser).Error; err != nil {
		log.Fatalf("Failed to create superuser: %v", err)
	}

	log.Println("Superuser 'superuser@example.com' created successfully.")
}

// SeedSubscriptionPackages creates default subscription packages if they do not already exist.
func SeedSubscriptionPackages() {
	packages := []models.SubscriptionPackageTable{
		{
			Name:         "Basic",
			Price:        50000.00,
			MaxEmployees: 10,
			Features:     "Basic attendance tracking, 1 admin user",
		},
		{
			Name:         "Standard",
			Price:        100000.00,
			MaxEmployees: 50,
			Features:     "All Basic features, advanced reporting, 3 admin users",
		},
		{
			Name:         "Premium",
			Price:        250000.00,
			MaxEmployees: 200,
			Features:     "All Standard features, unlimited admin users, API access, priority support",
		},
	}

	for _, pkg := range packages {
		var existingPackage models.SubscriptionPackageTable
		result := DB.Where("name = ?", pkg.Name).First(&existingPackage)

		if result.Error == gorm.ErrRecordNotFound {
			// Package does not exist, create it
			if err := DB.Create(&pkg).Error; err != nil {
				log.Printf("Failed to create subscription package %s: %v", pkg.Name, err)
			} else {
				log.Printf("Subscription package %s created successfully.", pkg.Name)
			}
		} else if result.Error != nil {
			log.Printf("Error checking for subscription package %s: %v", pkg.Name, result.Error)
		} else {
			log.Printf("Subscription package %s already exists. Skipping seeding.", pkg.Name)
		}
	}
}
