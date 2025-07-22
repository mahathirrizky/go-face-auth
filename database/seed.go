package database

import (
	"go-face-auth/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedSuperAdmin creates a default superadmin if one does not already exist.
func SeedSuperAdmin() {
	// Check if a superadmin already exists
	var superAdmin models.SuperAdminTable
	result := DB.Where("email = ?", "superadmin@example.com").First(&superAdmin)

	if result.Error == nil {
		log.Println("Superadmin 'superadmin@example.com' already exists. Skipping seeding.")
		return
	}

	// If not found, create a new superadmin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash superadmin password: %v", err)
	}

	newSuperAdmin := models.SuperAdminTable{
		Email:    "superadmin@example.com",
		Password: string(hashedPassword),
		Role:     "super_admin",
	}

	if err := DB.Create(&newSuperAdmin).Error; err != nil {
		log.Fatalf("Failed to create superadminn: %v", err)
	}

	log.Println("Superadmin superadminin@example.com' created successfully.")
}

// SeedSubscriptionPackages creates default subscription packages if they do not already exist.
func SeedSubscriptionPackages() {
	packages := []models.SubscriptionPackageTable{
		{
			PackageName:  "Standard",
			PriceMonthly: 100000.00,
			PriceYearly:  1000000.00, // Example yearly price for Standard
			MaxEmployees: 50,
			MaxLocations: 2, // Updated MaxLocations for Standard
			Features:     "All Basic features, advanced reporting, 3 admin users",
			IsActive:     true,
		},
		{
			PackageName:  "Premium",
			PriceMonthly: 250000.00,
			PriceYearly:  2500000.00, // Example yearly price for Premium
			MaxEmployees: 200,
			MaxLocations: 3, // Updated MaxLocations for Premium
			Features:     "All Standard features, unlimited admin users, API access, priority support",
			IsActive:     true,
		},
	}

	for _, pkg := range packages {
		var existingPackage models.SubscriptionPackageTable
		result := DB.Where("package_name = ?", pkg.PackageName).First(&existingPackage)

		if result.Error == gorm.ErrRecordNotFound {
			// Package does not exist, create it
			if err := DB.Create(&pkg).Error; err != nil {
				log.Printf("Failed to create subscription package %s: %v", pkg.PackageName, err)
			} else {
				log.Printf("Subscription package %s created successfully.", pkg.PackageName)
			}
		} else if result.Error != nil {
			log.Printf("Error checking for subscription package %s: %v", pkg.PackageName, result.Error)
		} else {
			// Package exists, update its fields
			existingPackage.PriceMonthly = pkg.PriceMonthly
			existingPackage.PriceYearly = pkg.PriceYearly
			existingPackage.MaxEmployees = pkg.MaxEmployees
			existingPackage.MaxLocations = pkg.MaxLocations // Updated MaxLocations
			existingPackage.Features = pkg.Features
			existingPackage.IsActive = pkg.IsActive
			if err := DB.Save(&existingPackage).Error; err != nil {
				log.Printf("Failed to update subscription package %s: %v", pkg.PackageName, err)
			} else {
				log.Printf("Subscription package %s updated successfully.", pkg.PackageName)
			}
		}
	}
}
