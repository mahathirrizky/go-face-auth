package services

import (
	"fmt"
	"go-face-auth/models"
	"log"
	"time"

	"go-face-auth/database/repository"
)

// SuperAdminService defines the interface for super admin related business logic.
type SuperAdminService interface {
	GetSuperAdminDashboardSummary() (map[string]interface{}, error)
	GetCompanies() ([]models.CompaniesTable, error)
	GetSubscriptions() ([]models.CompaniesTable, error)
	GetRevenueSummary(startDateStr, endDateStr string) ([]MonthlyRevenue, error)
	GetCustomPackageRequests(page, pageSize int, search string) ([]models.CustomPackageRequest, int64, error)
	UpdateCustomPackageRequestStatus(requestID uint, newStatus string) error
}

// superAdminService is the concrete implementation of SuperAdminService.
type superAdminService struct {
	companyRepo              repository.CompanyRepository
	invoiceRepo              repository.InvoiceRepository
	customPackageRequestRepo repository.CustomPackageRequestRepository
	superAdminRepo           repository.SuperAdminRepository
}

// NewSuperAdminService creates a new instance of SuperAdminService.
func NewSuperAdminService(companyRepo repository.CompanyRepository, invoiceRepo repository.InvoiceRepository, customPackageRequestRepo repository.CustomPackageRequestRepository, superAdminRepo repository.SuperAdminRepository) SuperAdminService {
	return &superAdminService{
		companyRepo:              companyRepo,
		invoiceRepo:              invoiceRepo,
		customPackageRequestRepo: customPackageRequestRepo,
		superAdminRepo:           superAdminRepo,
	}
}

func (s *superAdminService) GetSuperAdminDashboardSummary() (map[string]interface{}, error) {
	totalCompanies, err := s.superAdminRepo.GetTotalCompaniesCount()
	if err != nil {
		return nil, fmt.Errorf("error counting total companies: %w", err)
	}

	activeSubscriptions, err := s.superAdminRepo.GetCompaniesCountBySubscriptionStatus("active")
	if err != nil {
		return nil, fmt.Errorf("error counting active subscriptions: %w", err)
	}

	expiredSubscriptions, err := s.superAdminRepo.GetExpiredAndTrialExpiredCompaniesCount()
	if err != nil {
		return nil, fmt.Errorf("error counting expired subscriptions: %w", err)
	}

	trialSubscriptions, err := s.superAdminRepo.GetCompaniesCountBySubscriptionStatus("trial")
	if err != nil {
		return nil, fmt.Errorf("error counting trial subscriptions: %w", err)
	}

	recentCompanies, err := s.superAdminRepo.GetRecentCompanies(5)
	if err != nil {
		log.Printf("Error fetching recent companies: %v", err)
	}

	recentActivities := make([]map[string]interface{}, len(recentCompanies))
	for i, company := range recentCompanies {
		recentActivities[i] = map[string]interface{}{
			"id":          company.ID,
			"description": fmt.Sprintf("Company %s registered", company.Name),
			"timestamp":   company.CreatedAt.UnixMilli(),
		}
	}

	return map[string]interface{}{
		"total_companies":       totalCompanies,
		"active_subscriptions":  activeSubscriptions,
		"expired_subscriptions": expiredSubscriptions,
		"trial_subscriptions":   trialSubscriptions,
		"recent_activities":     recentActivities,
	}, nil
}

func (s *superAdminService) GetCompanies() ([]models.CompaniesTable, error) {
	companies, err := s.superAdminRepo.GetAllCompaniesWithPreload()
	if err != nil {
		return nil, fmt.Errorf("error fetching companies: %w", err)
	}
	return companies, nil
}

func (s *superAdminService) GetSubscriptions() ([]models.CompaniesTable, error) {
	companies, err := s.superAdminRepo.GetAllCompaniesWithPreload()
	if err != nil {
		return nil, fmt.Errorf("error fetching subscriptions: %w", err)
	}
	return companies, nil
}

type MonthlyRevenue struct {
	Month        string  `json:"month"`
	Year         string  `json:"year"`
	TotalRevenue float64 `json:"total_revenue"`
}

func (s *superAdminService) GetRevenueSummary(startDateStr, endDateStr string) ([]MonthlyRevenue, error) {
	var startDate *time.Time
	if startDateStr != "" {
		parsedTime, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format. Use YYYY-MM-DD")
		}
		startDate = &parsedTime
	}

	var endDate *time.Time
	if endDateStr != "" {
		parsedTime, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format. Use YYYY-MM-DD")
		}
		// Add 23h 59m 59s to include the whole end day
		endOfDay := parsedTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endDate = &endOfDay
	}

	monthlyRevenueData, err := s.superAdminRepo.GetPaidInvoicesMonthlyRevenue(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve revenue summary: %w", err)
	}

	var monthlyRevenue []MonthlyRevenue
	for _, data := range monthlyRevenueData {
		monthlyRevenue = append(monthlyRevenue, MonthlyRevenue{
			Month:        data.Month,
			Year:         data.Year,
			TotalRevenue: data.TotalRevenue,
		})
	}

	return monthlyRevenue, nil
}

func (s *superAdminService) GetCustomPackageRequests(page, pageSize int, search string) ([]models.CustomPackageRequest, int64, error) {
	return s.customPackageRequestRepo.GetCustomPackageRequestsPaginated(page, pageSize, search)
}

func (s *superAdminService) UpdateCustomPackageRequestStatus(requestID uint, newStatus string) error {
	request, err := s.customPackageRequestRepo.GetCustomPackageRequestByID(requestID)
	if err != nil || request == nil {
		return fmt.Errorf("custom package request not found")
	}

	request.Status = newStatus

	if err := s.customPackageRequestRepo.UpdateCustomPackageRequest(request); err != nil {
		return fmt.Errorf("failed to update request status: %w", err)
	}

	return nil
}