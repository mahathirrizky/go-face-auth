package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

type CustomPackageRequestService interface {
	CreateCustomPackageRequest(companyID, adminID uint, phone, message string) (*models.CustomPackageRequest, error)
}

type customPackageRequestService struct {
	companyRepo              repository.CompanyRepository
	adminCompanyRepo         repository.AdminCompanyRepository
	customPackageRequestRepo repository.CustomPackageRequestRepository
}

func NewCustomPackageRequestService(companyRepo repository.CompanyRepository, adminCompanyRepo repository.AdminCompanyRepository, customPackageRequestRepo repository.CustomPackageRequestRepository) CustomPackageRequestService {
	return &customPackageRequestService{
		companyRepo:              companyRepo,
		adminCompanyRepo:         adminCompanyRepo,
		customPackageRequestRepo: customPackageRequestRepo,
	}
}

func (s *customPackageRequestService) CreateCustomPackageRequest(companyID, adminID uint, phone, message string) (*models.CustomPackageRequest, error) {
	company, err := s.companyRepo.GetCompanyByID(int(companyID))
	if err != nil || company == nil {
		return nil, fmt.Errorf("failed to retrieve company information")
	}

	admin, err := s.adminCompanyRepo.GetAdminCompanyByID(int(adminID))
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

	if err := s.customPackageRequestRepo.CreateCustomPackageRequest(customRequest); err != nil {
		return nil, fmt.Errorf("failed to submit custom package request: %w", err)
	}

	return customRequest, nil
}