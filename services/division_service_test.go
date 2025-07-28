package services

import (
	"errors"

	"go-face-auth/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)



// TestCreateDivision tests the business logic for creating a division.
func TestCreateDivision(t *testing.T) {
	// Test case 1: Successful creation
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		req := CreateDivisionRequest{
			Name:      "Engineering",
			CompanyID: 1,
		}
		expectedDivision := &models.DivisionTable{
			Name:      "Engineering",
			CompanyID: 1,
		}

		// We tell the mock what to expect and what to return
		mocks.DivisionRepo.On("IsDivisionNameTaken", "Engineering", uint(1), uint(0)).Return(false, nil).Once()
		mocks.DivisionRepo.On("CreateDivision", mock.AnythingOfType("*models.DivisionTable")).Return(expectedDivision, nil).Once()

		// Act
		createdDivision, err := service.CreateDivision(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, createdDivision)
		assert.Equal(t, "Engineering", createdDivision.Name)
		mocks.DivisionRepo.AssertExpectations(t) // Verify that the expected methods were called
	})

	// Test case 2: Division name is already taken
	t.Run("Name Taken", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		req := CreateDivisionRequest{
			Name:      "Marketing",
			CompanyID: 1,
		}

		// We tell the mock to return 'true' for the name check
		mocks.DivisionRepo.On("IsDivisionNameTaken", "Marketing", uint(1), uint(0)).Return(true, nil).Once()

		// Act
		createdDivision, err := service.CreateDivision(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, createdDivision)
		assert.Equal(t, "division name is already taken", err.Error())
		mocks.DivisionRepo.AssertExpectations(t)
		// We don't expect CreateDivision to be called, and AssertExpectations will fail if it is.
	})

	// Test case 3: Repository returns an error during name check
	t.Run("Repository Error on Name Check", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		req := CreateDivisionRequest{
			Name:      "HR",
			CompanyID: 1,
		}
		repoError := errors.New("database connection lost")

		// We tell the mock to return an error for the name check
		mocks.DivisionRepo.On("IsDivisionNameTaken", "HR", uint(1), uint(0)).Return(false, repoError).Once()

		// Act
		createdDivision, err := service.CreateDivision(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, createdDivision)
		assert.Contains(t, err.Error(), "failed to check division name")
		mocks.DivisionRepo.AssertExpectations(t)
	})
}

func TestGetDivisionsByCompanyID(t *testing.T) {
	// Test case 1: Successful retrieval
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		companyID := uint(1)
		expectedDivisions := []models.DivisionTable{
			{ID: 1, Name: "Engineering", CompanyID: companyID},
			{ID: 2, Name: "Marketing", CompanyID: companyID},
		}

		mocks.DivisionRepo.On("GetDivisionsByCompanyID", companyID).Return(expectedDivisions, nil).Once()

		// Act
		divisions, err := service.GetDivisionsByCompanyID(companyID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, divisions)
		assert.Equal(t, expectedDivisions, divisions)
		mocks.DivisionRepo.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository Error", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		companyID := uint(1)
		repoError := errors.New("database error")

		mocks.DivisionRepo.On("GetDivisionsByCompanyID", companyID).Return(nil, repoError).Once()

		// Act
		divisions, err := service.GetDivisionsByCompanyID(companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, divisions)
		assert.Equal(t, repoError, err)
		mocks.DivisionRepo.AssertExpectations(t)
	})
}

func TestGetDivisionByID(t *testing.T) {
	// Test case 1: Successful retrieval
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		expectedDivision := &models.DivisionTable{ID: divisionID, Name: "HR", CompanyID: companyID}

		mocks.DivisionRepo.On("GetDivisionByID", divisionID).Return(expectedDivision, nil).Once()

		// Act
		division, err := service.GetDivisionByID(divisionID, companyID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, division)
		assert.Equal(t, expectedDivision, division)
		mocks.DivisionRepo.AssertExpectations(t)
	})

	// Test case 2: Division not found in repository
	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		divisionID := uint(999)
		companyID := uint(100)

		mocks.DivisionRepo.On("GetDivisionByID", divisionID).Return(nil, errors.New("record not found")).Once()

		// Act
		division, err := service.GetDivisionByID(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, division)
		assert.Contains(t, err.Error(), "record not found")
		mocks.DivisionRepo.AssertExpectations(t)
	})

	// Test case 3: Division found but belongs to a different company
	t.Run("Wrong Company", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		// Division belongs to company 200, but we are querying for company 100
		foundDivision := &models.DivisionTable{ID: divisionID, Name: "HR", CompanyID: 200}

		mocks.DivisionRepo.On("GetDivisionByID", divisionID).Return(foundDivision, nil).Once()

		// Act
		division, err := service.GetDivisionByID(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, division)
		assert.Equal(t, "division not found in this company", err.Error())
		mocks.DivisionRepo.AssertExpectations(t)
	})

	// Test case 4: Repository error
	t.Run("Repository Error", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		repoError := errors.New("db connection failed")

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(nil, repoError).Once()

		// Act
		division, err := service.GetDivisionByID(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, division)
		assert.Equal(t, repoError, err)
		mockRepo.DivisionRepo.AssertExpectations(t)
	})
}

func TestUpdateDivision(t *testing.T) {
	// Test case 1: Successful update of name and description
	t.Run("Success Update Name and Description", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Old Name", CompanyID: companyID, Description: "Old Desc"}
		req := UpdateDivisionRequest{Name: "New Name", Description: "New Desc"}
		updatedDivision := &models.DivisionTable{ID: divisionID, Name: "New Name", CompanyID: companyID, Description: "New Desc"}

		mocks.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mocks.DivisionRepo.On("IsDivisionNameTaken", "New Name", companyID, divisionID).Return(false, nil).Once()
		mocks.DivisionRepo.On("UpdateDivision", mock.AnythingOfType("*models.DivisionTable")).Return(updatedDivision, nil).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "New Desc", result.Description)
		mocks.DivisionRepo.AssertExpectations(t)
	})

	// Test case 2: Successful update of only description (name not provided)
	t.Run("Success Update Only Description", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Original Name", CompanyID: companyID, Description: "Old Desc"}
		req := UpdateDivisionRequest{Description: "Updated Desc"}
		updatedDivision := &models.DivisionTable{ID: divisionID, Name: "Original Name", CompanyID: companyID, Description: "Updated Desc"}

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.DivisionRepo.On("UpdateDivision", mock.AnythingOfType("*models.DivisionTable")).Return(updatedDivision, nil).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Original Name", result.Name)
		assert.Equal(t, "Updated Desc", result.Description)
		mockRepo.DivisionRepo.AssertExpectations(t)
	})

	// Test case 3: Successful update of shifts
	t.Run("Success Update Shifts", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Test Div", CompanyID: companyID}
		shiftIDs := []uint{1, 2}
		shifts := []models.ShiftsTable{{ID: 1}, {ID: 2}}
		req := UpdateDivisionRequest{ShiftIDs: shiftIDs}
		updatedDivision := &models.DivisionTable{ID: divisionID, Name: "Test Div", CompanyID: companyID, Shifts: shifts}

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.ShiftRepo.On("GetShiftsByIDs", shiftIDs).Return(shifts, nil).Once()
		mockRepo.DivisionRepo.On("UpdateDivision", mock.AnythingOfType("*models.DivisionTable")).Return(updatedDivision, nil).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, shifts, result.Shifts)
		mockRepo.DivisionRepo.AssertExpectations(t)
		mockRepo.ShiftRepo.AssertExpectations(t)
	})

	// Test case 4: Successful update of locations
	t.Run("Success Update Locations", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Test Div", CompanyID: companyID}
		locationIDs := []uint{1, 2}
				locations := []models.AttendanceLocation{{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}}}
		req := UpdateDivisionRequest{LocationIDs: locationIDs}
		updatedDivision := &models.DivisionTable{ID: divisionID, Name: "Test Div", CompanyID: companyID, Locations: locations}

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.AttendanceLocationRepo.On("GetLocationsByIDs", locationIDs).Return(locations, nil).Once()
		mockRepo.DivisionRepo.On("UpdateDivision", mock.AnythingOfType("*models.DivisionTable")).Return(updatedDivision, nil).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, locations, result.Locations)
		mockRepo.DivisionRepo.AssertExpectations(t)
		mockRepo.AttendanceLocationRepo.AssertExpectations(t)
	})

	// Test case 5: Division not found
	t.Run("Division Not Found", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		req := UpdateDivisionRequest{Name: "New Name"}

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(nil, errors.New("division not found")).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "division not found")
		mockRepo.DivisionRepo.AssertExpectations(t)
	})

	// Test case 6: Division name already taken
	t.Run("Name Already Taken", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Old Name", CompanyID: companyID}
		req := UpdateDivisionRequest{Name: "Taken Name"}

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.DivisionRepo.On("IsDivisionNameTaken", "Taken Name", companyID, divisionID).Return(true, nil).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "division name is already taken", err.Error())
		mockRepo.DivisionRepo.AssertExpectations(t)
	})

	// Test case 7: Repository error during name check
	t.Run("Repo Error Name Check", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Old Name", CompanyID: companyID}
		req := UpdateDivisionRequest{Name: "New Name"}
		repoError := errors.New("db error during name check")

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.DivisionRepo.On("IsDivisionNameTaken", "New Name", companyID, divisionID).Return(false, repoError).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to check division name")
		mockRepo.DivisionRepo.AssertExpectations(t)
	})

	// Test case 8: Repository error during update
	t.Run("Repo Error Update", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Old Name", CompanyID: companyID}
		req := UpdateDivisionRequest{Name: "New Name"}
		repoError := errors.New("db error during update")

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.DivisionRepo.On("IsDivisionNameTaken", "New Name", companyID, divisionID).Return(false, nil).Once()
		mockRepo.DivisionRepo.On("UpdateDivision", mock.AnythingOfType("*models.DivisionTable")).Return(nil, repoError).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoError, err)
		mockRepo.DivisionRepo.AssertExpectations(t)
	})

	// Test case 9: Error fetching shifts
	t.Run("Error Fetching Shifts", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Test Div", CompanyID: companyID}
		shiftIDs := []uint{1, 2}
		req := UpdateDivisionRequest{ShiftIDs: shiftIDs}
		repoError := errors.New("shift db error")

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.ShiftRepo.On("GetShiftsByIDs", shiftIDs).Return(nil, repoError).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to retrieve shifts")
		mockRepo.DivisionRepo.AssertExpectations(t)
		mockRepo.ShiftRepo.AssertExpectations(t)
	})

	// Test case 10: Not all shifts found
	t.Run("Not All Shifts Found", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Test Div", CompanyID: companyID}
		shiftIDs := []uint{1, 2}
		shifts := []models.ShiftsTable{{ID: 1}}
		req := UpdateDivisionRequest{ShiftIDs: shiftIDs}

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.ShiftRepo.On("GetShiftsByIDs", shiftIDs).Return(shifts, nil).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "one or more selected shifts not found")
		mockRepo.DivisionRepo.AssertExpectations(t)
		mockRepo.ShiftRepo.AssertExpectations(t)
	})

	// Test case 11: Error fetching locations
	t.Run("Error Fetching Locations", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Test Div", CompanyID: companyID}
		locationIDs := []uint{1, 2}
		req := UpdateDivisionRequest{LocationIDs: locationIDs}
		repoError := errors.New("location db error")

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.AttendanceLocationRepo.On("GetLocationsByIDs", locationIDs).Return(nil, repoError).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to retrieve locations")
		mockRepo.DivisionRepo.AssertExpectations(t)
		mockRepo.AttendanceLocationRepo.AssertExpectations(t)
	})

	// Test case 12: Not all locations found
	t.Run("Not All Locations Found", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Test Div", CompanyID: companyID}
		locationIDs := []uint{1, 2}
		locations := []models.AttendanceLocation{{Model: gorm.Model{ID: 1}}}
		req := UpdateDivisionRequest{LocationIDs: locationIDs}

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.AttendanceLocationRepo.On("GetLocationsByIDs", locationIDs).Return(locations, nil).Once()

		// Act
		result, err := service.UpdateDivision(divisionID, companyID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "one or more selected locations not found")
		mockRepo.DivisionRepo.AssertExpectations(t)
		mockRepo.AttendanceLocationRepo.AssertExpectations(t)
	})
}

func TestDeleteDivision(t *testing.T) {
	// Test case 1: Successful deletion
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mocks := NewMockRepositories()
		service := NewDivisionService(mocks.DivisionRepo, mocks.ShiftRepo, mocks.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "To Delete", CompanyID: companyID}

		mocks.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mocks.DivisionRepo.On("DeleteDivision", divisionID).Return(nil).Once()

		// Act
		err := service.DeleteDivision(divisionID, companyID)

		// Assert
		assert.NoError(t, err)
		mocks.DivisionRepo.AssertExpectations(t)
	})

	// Test case 2: Division not found
	t.Run("Division Not Found", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(999)
		companyID := uint(100)

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(nil, errors.New("division not found")).Once()

		// Act
		err := service.DeleteDivision(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "division not found")
		mockRepo.DivisionRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error during delete
	t.Run("Repo Error Delete", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockRepositories()
		service := NewDivisionService(mockRepo.DivisionRepo, mockRepo.ShiftRepo, mockRepo.AttendanceLocationRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "To Delete", CompanyID: companyID}
		repoError := errors.New("db error during delete")

		mockRepo.DivisionRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.DivisionRepo.On("DeleteDivision", divisionID).Return(repoError).Once()

		// Act
		err := service.DeleteDivision(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.DivisionRepo.AssertExpectations(t)
	})
}