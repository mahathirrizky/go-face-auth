package repository

import (
	"go-face-auth/models"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type leaveRequestRepository struct {
	db *gorm.DB
}

func NewLeaveRequestRepository(db *gorm.DB) LeaveRequestRepository {
	return &leaveRequestRepository{db: db}
}

// CreateLeaveRequest inserts a new leave request into the database.
func (r *leaveRequestRepository) CreateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	result := r.db.Create(leaveRequest)
	if result.Error != nil {
		log.Printf("Error creating leave request: %v", result.Error)
		return result.Error
	}
	log.Printf("Leave request created with ID: %d", leaveRequest.ID)
	return nil
}

// GetLeaveRequestByID retrieves a leave request by its ID.
func (r *leaveRequestRepository) GetLeaveRequestByID(id uint) (*models.LeaveRequest, error) {
	var leaveRequest models.LeaveRequest
	result := r.db.Preload("Employee").First(&leaveRequest, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Leave request not found
		}
		log.Printf("Error getting leave request with ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &leaveRequest, nil
}

// GetAllLeaveRequests retrieves all leave requests.
func (r *leaveRequestRepository) GetAllLeaveRequests() ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	result := r.db.Preload("Employee").Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Error getting all leave requests: %v", result.Error)
		return nil, result.Error
	}
	return leaveRequests, nil
}

// GetLeaveRequestsByEmployeeID retrieves all leave requests for a given employee ID, optionally filtered by date range.
func (r *leaveRequestRepository) GetLeaveRequestsByEmployeeID(employeeID uint, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	query := r.db.Preload("Employee").Where("employee_id = ?", employeeID)

	if startDate != nil {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("end_date <= ?", endDate)
	}

	result := query.Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Error getting leave requests for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return leaveRequests, nil
}

// GetCompanyLeaveRequestsFiltered retrieves all leave requests for a company with filtering, without pagination.
func (r *leaveRequestRepository) GetCompanyLeaveRequestsFiltered(companyID int, status, search string, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	query := r.db.Preload("Employee").Where("employee_id IN (SELECT id FROM employees WHERE company_id = ?)", companyID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("LOWER(reason) LIKE ? OR LOWER(type) LIKE ? OR LOWER(Employee.name) LIKE ?", "%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%")
	}

	if startDate != nil {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != nil {
		// Add 23 hours, 59 minutes, 59 seconds to include the entire end day
		query = query.Where("end_date <= ?", endDate.Add(23*time.Hour+59*time.Minute+59*time.Second))
	}

	err := query.Order("created_at DESC").Find(&leaveRequests).Error
	return leaveRequests, err
}

// UpdateLeaveRequest updates an existing leave request record in the database.
func (r *leaveRequestRepository) UpdateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	result := r.db.Save(leaveRequest)
	if result.Error != nil {
		log.Printf("Error updating leave request: %v", result.Error)
		return result.Error
	}
	log.Printf("Leave request updated with ID: %d", leaveRequest.ID)
	return nil
}

// GetRecentLeaveRequestsByCompanyID retrieves recent leave requests for a given company ID.
func (r *leaveRequestRepository) GetRecentLeaveRequestsByCompanyID(companyID int, limit int) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	log.Printf("Repository: Fetching recent leave requests for company %d, limit %d", companyID, limit)
	result := r.db.Preload("Employee").Joins("join employees_tables on leave_requests.employee_id = employees_tables.id").Where("employees_tables.company_id = ?", companyID).Order("created_at DESC").Limit(limit).Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Repository: Error getting recent leave requests for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	log.Printf("Repository: Found %d recent leave requests for company %d.", len(leaveRequests), companyID)
	return leaveRequests, nil
}

// IsEmployeeOnApprovedLeave retrieves an approved leave request for a specific employee on a given date.
func (r *leaveRequestRepository) IsEmployeeOnApprovedLeave(employeeID int, date time.Time) (*models.LeaveRequest, error) {
	var leaveRequest models.LeaveRequest
	// Normalize date to start of day for comparison
	checkDate := date.Truncate(24 * time.Hour)

	result := r.db.Model(&models.LeaveRequest{}).Where("employee_id = ? AND status = ? AND start_date <= ? AND end_date >= ?", employeeID, "approved", checkDate, checkDate).First(&leaveRequest)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No approved leave record found
		}
		log.Printf("Error checking approved leave for employee %d on date %s: %v", employeeID, date.Format("2006-01-02"), result.Error)
		return nil, result.Error
	}
	return &leaveRequest, nil
}

// IsEmployeeOnApprovedLeaveDateRange checks if an employee has an approved leave or sick request within a specific date range.
func (r *leaveRequestRepository) IsEmployeeOnApprovedLeaveDateRange(employeeID int, startDate, endDate *time.Time) (bool, error) {
	var count int64
	query := r.db.Model(&models.LeaveRequest{}).Where("employee_id = ? AND status = ?", employeeID, "approved")

	// Check for overlap with the requested date range
	if startDate != nil {
		query = query.Where("end_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("start_date <= ?", *endDate)
	}

	result := query.Count(&count)
	if result.Error != nil {
		log.Printf("Error checking approved leave for employee %d in date range: %v", employeeID, result.Error)
		return false, result.Error
	}
	return count > 0, nil
}

// GetPendingLeaveRequestsByEmployeeID retrieves all pending leave requests for a given employee ID.
func (r *leaveRequestRepository) GetPendingLeaveRequestsByEmployeeID(employeeID int) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	result := r.db.Preload("Employee").Where("employee_id = ? AND status = ?", employeeID, "pending").Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Error getting pending leave requests for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return leaveRequests, nil
}

// GetCompanyLeaveRequestsPaginated retrieves paginated and filtered leave requests for a company.
func (r *leaveRequestRepository) GetCompanyLeaveRequestsPaginated(companyID int, status, search string, startDate, endDate *time.Time, page, pageSize int) ([]models.LeaveRequest, int64, error) {
	var leaveRequests []models.LeaveRequest
	var totalRecords int64

	// Base query
	query := r.db.Model(&models.LeaveRequest{}).
		Joins("JOIN employees_tables ON leave_requests.employee_id = employees_tables.id").
		Where("employees_tables.company_id = ?", companyID)

	// Apply filters
	if status != "" {
		query = query.Where("leave_requests.status = ?", status)
	}
	if search != "" {
		query = query.Where("employees_tables.name ILIKE ?", "%"+search+"%")
	}
	if startDate != nil {
		query = query.Where("leave_requests.start_date >= ?", startDate)
	}
	if endDate != nil {
		// Add 23 hours, 59 minutes, 59 seconds to include the entire end day
		query = query.Where("leave_requests.end_date <= ?", endDate.Add(23*time.Hour+59*time.Minute+59*time.Second))
	}

	// Get total records count
	if err := query.Count(&totalRecords).Error; err != nil {
		log.Printf("Error counting leave requests: %v", err)
		return nil, 0, err
	}

	// Apply pagination and order
	offset := (page - 1) * pageSize
	result := query.Preload("Employee").
		Order("leave_requests.created_at DESC").
		Offset(offset).
		Limit(pageSize).Find(&leaveRequests)

	if result.Error != nil {
		log.Printf("Error getting paginated leave requests: %v", result.Error)
		return nil, 0, result.Error
	}

	return leaveRequests, totalRecords, nil
}

func (r *leaveRequestRepository) GetOnLeaveEmployeesCountToday(companyID int) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02") // Format to YYYY-MM-DD for date comparison

	result := r.db.Model(&models.LeaveRequest{}).Where("employee_id IN (SELECT id FROM employees WHERE company_id = ?) AND status = ? AND ? BETWEEN start_date AND end_date", companyID, "approved", today).Count(&count)

	if result.Error != nil {
		log.Printf("Error counting on-leave employees today for company %d: %v", companyID, result.Error)
		return 0, result.Error
	}
	return count, nil
}