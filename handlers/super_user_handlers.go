package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
)

// SuperUserHandler handles HTTP requests related to super users.
type SuperUserHandler struct {
}

// NewSuperUserHandler creates a new SuperUserHandler.
func NewSuperUserHandler() *SuperUserHandler {
	return &SuperUserHandler{}
}

// CreateSuperUser handles the creation of a new super user.
func (h *SuperUserHandler) CreateSuperUser(c *gin.Context) {
	// Check if a super user already exists
	existingSuperUsers, err := repository.GetAllSuperUsers()
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(existingSuperUsers) > 0 {
		helper.SendError(c, http.StatusConflict, "Only one super user is allowed.")
		return
	}

	var superUser models.SuperUserTable
	if err := c.BindJSON(&superUser); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := repository.CreateSuperUser(&superUser); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Super user created successfully", superUser)
}

// GetSuperUserByID handles fetching a super user by its ID.
func (h *SuperUserHandler) GetSuperUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid super user ID")
		return
	}

	superUser, err := repository.GetSuperUserByID(id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if superUser == nil {
		helper.SendError(c, http.StatusNotFound, "Super user not found")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Super user fetched successfully", superUser)
}

// GetSuperUserByEmail handles fetching a super user by its email.
func (h *SuperUserHandler) GetSuperUserByEmail(c *gin.Context) {
	email := c.Param("email")

	superUser, err := repository.GetSuperUserByEmail(email)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if superUser == nil {
		helper.SendError(c, http.StatusNotFound, "Super user not found for this email")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Super user fetched successfully", superUser)
}

// UpdateSuperUser handles updating an existing super user.
func (h *SuperUserHandler) UpdateSuperUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid super user ID")
		return
	}

	var superUser models.SuperUserTable
	if err := c.BindJSON(&superUser); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedSuperUser, err := repository.UpdateSuperUser(id, &superUser)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return	
	}

	helper.SendSuccess(c, http.StatusOK, "Super user updated successfully", updatedSuperUser)
}

// DeleteSuperUser handles deleting a super user by its ID.
func (h *SuperUserHandler) DeleteSuperUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid super user ID")
		return
	}

	if err := repository.DeleteSuperUser(id); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Super user deleted successfully", nil)
}

// GetAllSuperUsers handles fetching all super users.
func (h *SuperUserHandler) GetAllSuperUsers(c *gin.Context) {
	superUsers, err := repository.GetAllSuperUsers()
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Super users fetched successfully", superUsers)
}
