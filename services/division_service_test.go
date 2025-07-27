package services

import (
	"errors"

	"go-face-auth/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDivisionRepository is a mock implementation of DivisionRepository for testing.
type MockDivisionRepository struct {
	mock.Mock
}

// We implement the interface methods on the mock struct.
func (m *MockDivisionRepository) CreateDivision(division *models.DivisionTable) (*models.DivisionTable, error) {
	args := m.Called(division)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DivisionTable), args.Error(1)
}

func (m *MockDivisionRepository) IsDivisionNameTaken(name string, companyID uint, currentDivisionID uint) (bool, error) {
	args := m.Called(name, companyID, currentDivisionID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDivisionRepository) GetDivisionsByCompanyID(companyID uint) ([]models.DivisionTable, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.DivisionTable), args.Error(1)
}

func (m *MockDivisionRepository) GetDivisionByID(divisionID uint) (*models.DivisionTable, error) {
	args := m.Called(divisionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DivisionTable), args.Error(1)
}

func (m *MockDivisionRepository) UpdateDivision(division *models.DivisionTable) (*models.DivisionTable, error) {
	args := m.Called(division)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DivisionTable), args.Error(1)
}

func (m *MockDivisionRepository) DeleteDivision(divisionID uint) error {
	args := m.Called(divisionID)
	return args.Error(0)
}

// TestCreateDivision tests the business logic for creating a division.
func TestCreateDivision(t *testing.T) {
	// Test case 1: Successful creation
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		input := &models.DivisionTable{
			Name:      "Engineering",
			CompanyID: 1,
		}

		// We tell the mock what to expect and what to return
		mockRepo.On("IsDivisionNameTaken", "Engineering", uint(1), uint(0)).Return(false, nil).Once()
		mockRepo.On("CreateDivision", input).Return(input, nil).Once()

		// Act
		createdDivision, err := service.CreateDivision(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, createdDivision)
		assert.Equal(t, "Engineering", createdDivision.Name)
		mockRepo.AssertExpectations(t) // Verify that the expected methods were called
	})

	// Test case 2: Division name is already taken
	t.Run("Name Taken", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		input := &models.DivisionTable{
			Name:      "Marketing",
			CompanyID: 1,
		}

		// We tell the mock to return 'true' for the name check
		mockRepo.On("IsDivisionNameTaken", "Marketing", uint(1), uint(0)).Return(true, nil).Once()

		// Act
		createdDivision, err := service.CreateDivision(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, createdDivision)
		assert.Equal(t, "division name is already taken", err.Error())
		mockRepo.AssertExpectations(t)
		// We don't expect CreateDivision to be called, and AssertExpectations will fail if it is.
	})

	// Test case 3: Repository returns an error during name check
	t.Run("Repository Error on Name Check", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		input := &models.DivisionTable{
			Name:      "HR",
			CompanyID: 1,
		}
		repoError := errors.New("database connection lost")

		// We tell the mock to return an error for the name check
		mockRepo.On("IsDivisionNameTaken", "HR", uint(1), uint(0)).Return(false, repoError).Once()

		// Act
		createdDivision, err := service.CreateDivision(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, createdDivision)
		assert.Contains(t, err.Error(), "failed to check division name")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetDivisionsByCompanyID(t *testing.T) {
	// Test case 1: Successful retrieval
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		companyID := uint(1)
		expectedDivisions := []models.DivisionTable{
			{ID: 1, Name: "Engineering", CompanyID: companyID},
			{ID: 2, Name: "Marketing", CompanyID: companyID},
		}

		mockRepo.On("GetDivisionsByCompanyID", companyID).Return(expectedDivisions, nil).Once()

		// Act
		divisions, err := service.GetDivisionsByCompanyID(companyID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, divisions)
		assert.Equal(t, expectedDivisions, divisions)
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		companyID := uint(1)
		repoError := errors.New("database error")

		mockRepo.On("GetDivisionsByCompanyID", companyID).Return(nil, repoError).Once()

		// Act
		divisions, err := service.GetDivisionsByCompanyID(companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, divisions)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetDivisionByID(t *testing.T) {
	// Test case 1: Successful retrieval
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		expectedDivision := &models.DivisionTable{ID: divisionID, Name: "HR", CompanyID: companyID}

		mockRepo.On("GetDivisionByID", divisionID).Return(expectedDivision, nil).Once()

		// Act
		division, err := service.GetDivisionByID(divisionID, companyID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, division)
		assert.Equal(t, expectedDivision, division)
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Division not found in repository
	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(999)
		companyID := uint(100)

		mockRepo.On("GetDivisionByID", divisionID).Return(nil, errors.New("record not found")).Once()

		// Act
		division, err := service.GetDivisionByID(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, division)
		assert.Contains(t, err.Error(), "record not found")
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Division found but belongs to a different company
	t.Run("Wrong Company", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		// Division belongs to company 200, but we are querying for company 100
		foundDivision := &models.DivisionTable{ID: divisionID, Name: "HR", CompanyID: 200}

		mockRepo.On("GetDivisionByID", divisionID).Return(foundDivision, nil).Once()

		// Act
		division, err := service.GetDivisionByID(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, division)
		assert.Equal(t, "division not found in this company", err.Error())
		mockRepo.AssertExpectations(t)
	})

	// Test case 4: Repository error
	t.Run("Repository Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		repoError := errors.New("db connection failed")

		mockRepo.On("GetDivisionByID", divisionID).Return(nil, repoError).Once()

		// Act
		division, err := service.GetDivisionByID(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, division)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateDivision(t *testing.T) {
	// Test case 1: Successful update
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Old Name", CompanyID: companyID, Description: "Old Desc"}
		updatedInput := &models.DivisionTable{ID: divisionID, Name: "New Name", CompanyID: companyID, Description: "New Desc"}

		mockRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.On("IsDivisionNameTaken", "New Name", companyID, divisionID).Return(false, nil).Once()
		mockRepo.On("UpdateDivision", mock.AnythingOfType("*models.DivisionTable")).Return(updatedInput, nil).Once()

		// Act
		result, err := service.UpdateDivision(updatedInput, companyID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "New Desc", result.Description)
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Division not found
	t.Run("Division Not Found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		updatedInput := &models.DivisionTable{ID: divisionID, Name: "New Name", CompanyID: companyID}

		mockRepo.On("GetDivisionByID", divisionID).Return(nil, errors.New("division not found")).Once()

		// Act
		result, err := service.UpdateDivision(updatedInput, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "division not found")
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Division name already taken
	t.Run("Name Already Taken", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Old Name", CompanyID: companyID}
		updatedInput := &models.DivisionTable{ID: divisionID, Name: "Taken Name", CompanyID: companyID}

		mockRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.On("IsDivisionNameTaken", "Taken Name", companyID, divisionID).Return(true, nil).Once()

		// Act
		result, err := service.UpdateDivision(updatedInput, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "division name is already taken", err.Error())
		mockRepo.AssertExpectations(t)
	})

	// Test case 4: Repository error during name check
	t.Run("Repo Error Name Check", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Old Name", CompanyID: companyID}
		updatedInput := &models.DivisionTable{ID: divisionID, Name: "New Name", CompanyID: companyID}
		repoError := errors.New("db error during name check")

		mockRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.On("IsDivisionNameTaken", "New Name", companyID, divisionID).Return(false, repoError).Once()

		// Act
		result, err := service.UpdateDivision(updatedInput, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to check division name")
		mockRepo.AssertExpectations(t)
	})

	// Test case 5: Repository error during update
	t.Run("Repo Error Update", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "Old Name", CompanyID: companyID}
		updatedInput := &models.DivisionTable{ID: divisionID, Name: "New Name", CompanyID: companyID}
		repoError := errors.New("db error during update")

		mockRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.On("IsDivisionNameTaken", "New Name", companyID, divisionID).Return(false, nil).Once()
		mockRepo.On("UpdateDivision", mock.AnythingOfType("*models.DivisionTable")).Return(nil, repoError).Once()

		// Act
		result, err := service.UpdateDivision(updatedInput, companyID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteDivision(t *testing.T) {
	// Test case 1: Successful deletion
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "To Delete", CompanyID: companyID}

		mockRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.On("DeleteDivision", divisionID).Return(nil).Once()

		// Act
		err := service.DeleteDivision(divisionID, companyID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Division not found
	t.Run("Division Not Found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(999)
		companyID := uint(100)

		mockRepo.On("GetDivisionByID", divisionID).Return(nil, errors.New("division not found")).Once()

		// Act
		err := service.DeleteDivision(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "division not found")
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error during delete
	t.Run("Repo Error Delete", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockDivisionRepository)
		service := NewDivisionService(mockRepo)

		divisionID := uint(1)
		companyID := uint(100)
		existingDivision := &models.DivisionTable{ID: divisionID, Name: "To Delete", CompanyID: companyID}
		repoError := errors.New("db error during delete")

		mockRepo.On("GetDivisionByID", divisionID).Return(existingDivision, nil).Once()
		mockRepo.On("DeleteDivision", divisionID).Return(repoError).Once()

		// Act
		err := service.DeleteDivision(divisionID, companyID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}