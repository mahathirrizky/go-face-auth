package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

func CreateCustomPackageRequest(companyID, adminID uint, phone, message string) (*models.CustomPackageRequest, error) {
	company, err := repository.GetCompanyByID(int(companyID))
	if err != nil || company == nil {
		return nil, fmt.Errorf("failed to retrieve company information")
	}

	admin, err := repository.GetAdminCompanyByID(int(adminID))
	if err != nil || admin == nil {
		return nil, fmt.Errorf("failed to retrieve admin information")
	}

	customRequest := &models.CustomPackageRequest{
		CompanyID:   companyID,
		Email:       admin.Email,
		Phone:       phone,
		CompanyName: company.Name,
		Message:     message,
		Status:      "pending",
	}

	if err := repository.CreateCustomPackageRequest(customRequest); err != nil {
		return nil, fmt.Errorf("failed to submit custom package request: %w", err)
	}

	return customRequest, nil
}
