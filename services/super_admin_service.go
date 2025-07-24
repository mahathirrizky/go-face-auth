package services

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/models"
	"log"
	"time"

	"go-face-auth/database/repository"
)

func GetSuperAdminDashboardSummary() (map[string]interface{}, error) {
	var totalCompanies int64
	if err := database.DB.Model(&models.CompaniesTable{}).Count(&totalCompanies).Error; err != nil {
		return nil, fmt.Errorf("error counting total companies: %w", err)
	}

	var activeSubscriptions int64
	if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ?", "active").Count(&activeSubscriptions).Error; err != nil {
		return nil, fmt.Errorf("error counting active subscriptions: %w", err)
	}

	var expiredSubscriptions int64
	if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ? OR subscription_status = ?", "expired", "expired_trial").Count(&expiredSubscriptions).Error; err != nil {
		return nil, fmt.Errorf("error counting expired subscriptions: %w", err)
	}

	var trialSubscriptions int64
	if err := database.DB.Model(&models.CompaniesTable{}).Where("subscription_status = ?", "trial").Count(&trialSubscriptions).Error; err != nil {
		return nil, fmt.Errorf("error counting trial subscriptions: %w", err)
	}

	var recentCompanies []models.CompaniesTable
	if err := database.DB.Order("created_at DESC").Limit(5).Find(&recentCompanies).Error; err != nil {
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

func GetCompanies() ([]models.CompaniesTable, error) {
	var companies []models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").Find(&companies).Error; err != nil {
		return nil, fmt.Errorf("error fetching companies: %w", err)
	}
	return companies, nil
}

func GetSubscriptions() ([]models.CompaniesTable, error) {
	var companies []models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").Find(&companies).Error; err != nil {
		return nil, fmt.Errorf("error fetching subscriptions: %w", err)
	}
	return companies, nil
}

type MonthlyRevenue struct {
	Month        string  `json:"month"`
	Year         string  `json:"year"`
	TotalRevenue float64 `json:"total_revenue"`
}

func GetRevenueSummary(startDateStr, endDateStr string) ([]MonthlyRevenue, error) {
	var monthlyRevenue []MonthlyRevenue

	query := database.DB.Model(&models.InvoiceTable{}).Where("status = ?", "paid")

	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format. Use YYYY-MM-DD")
		}
		query = query.Where("created_at >= ?", startDate)
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format. Use YYYY-MM-DD")
		}
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("created_at <= ?", endDate)
	}

	if err := query.Select(
		"DATE_FORMAT(created_at, '%Y-%m') AS month, DATE_FORMAT(created_at, '%Y') AS year, SUM(amount) AS total_revenue").Group("month, year").Order("year DESC, month DESC").Scan(&monthlyRevenue).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve revenue summary: %w", err)
	}

	return monthlyRevenue, nil
}

func GetCustomPackageRequests(page, pageSize int, search string) ([]models.CustomPackageRequest, int64, error) {
	return repository.GetCustomPackageRequestsPaginated(page, pageSize, search)
}

func UpdateCustomPackageRequestStatus(requestID uint, newStatus string) error {
	request, err := repository.GetCustomPackageRequestByID(requestID)
	if err != nil || request == nil {
		return fmt.Errorf("custom package request not found")
	}

	request.Status = newStatus

	if err := repository.UpdateCustomPackageRequest(request); err != nil {
		return fmt.Errorf("failed to update request status: %w", err)
	}

	return nil
}
