package handler

import (
	"errors"
	"net/http"
	"prac/todo"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create a new user
// @Security ApiKeyAuth
// @Description Create a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param input body todo.User true "User data"
// @Success 200 {integer} integer 1
// @Failure 400 {object} myerror
// @Failure 500 {object} myerror
// @Router /api/users  [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var input todo.User
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	id, err := h.services.User.CreateUser(ctx, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// GetAllUsers godoc
// @Summary Get users list
// @Security ApiKeyAuth
// @Description Get all existing users
// @Tags users
// @Produce json
// @Success 200 {object} UsersResponse
// @Failure 500 {object} myerror
// @Router /api/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.services.User.GetAllUsers(ctx)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, UsersResponse{
		Users: users,
	})
}

// або без тексту і без дто вийде {array} []User
type UsersResponse struct {
	Users []todo.User `json:"users"`
}

// GetUserByID godoc
// @Summary Get user by ID
// @Security ApiKeyAuth
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} todo.User
// @Failure 400 {object} myerror
// @Failure 500 {object} myerror
// @Router /api/users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}
	ctx := c.Request.Context()
	user, err := h.services.User.GetUserByID(ctx, uint(id))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetUserByEmail godoc
// @Summary Get user by email
// @Security ApiKeyAuth
// @Description Get user by email
// @Tags users
// @Accept json
// @Produce json
// @Param emial path string true "User Emial"
// @Success 200 {object} todo.User
// @Failure 400 {object} myerror
// @Failure 500 {object} myerror
// @Router /api/users/{email} [get]
func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	ctx := c.Request.Context()
	user, err := h.services.User.GetUserByEmail(ctx, email)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUser godoc
// @Summary Update user
// @Security ApiKeyAuth
// @Description Update user by id
// @Tags users
// @Accept json
// @Produce json
// @Param input body todo.UpdateUserInput true "Updated data"
// @Success 200 {object} todo.User
// @Failure 400 {object} myerror
// @Failure 500 {object} myerror
// @Router /api/users/{id} [patch]
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}

	var input todo.UpdateUserInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	updatedUser, err := h.services.User.UpdateUser(ctx, uint(id), input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": updatedUser,
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Security ApiKeyAuth
// @Description Delete user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "message: user deleted successfully"
// @Failure 400 {object} myerror
// @Failure 500 {object} myerror
// @Router /api/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}
	ctx := c.Request.Context()
	err = h.services.User.DeleteUser(ctx, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

//

// GetProfile godoc
// @Summary Get profile
// @Security ApiKeyAuth
// @Description Get user profile by id form ctx
// @Tags profile
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} todo.User
// @Failure 400 {object} myerror
// @Failure 500 {object} myerror
// @Router /api/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	ctx := c.Request.Context()
	user, err := h.services.User.GetUserByID(ctx, userID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	})
}

// UpdateProfile godoc
// @Summary Update Profile
// @Security ApiKeyAuth
// @Description Update Profile
// @Tags profile
// @Accept json
// @Param input body todo.UpdateUserInput true "Updated data"
// @Success 200 {object} todo.User
// @Failure 400 {object} myerror
// @Failure 500 {object} myerror
// @Router /api/profile [patch]
func (h *Handler) UpdateProfile(c *gin.Context) { // з бівера зробити
	userID, err := h.getUserID(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var input todo.UpdateUserInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	updatedUser, err := h.services.User.UpdateUser(ctx, userID, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    updatedUser,
		"message": "Profile updated successfully",
	})
}

func (h *Handler) getUserID(c *gin.Context) (uint, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}
	// sminyty type danih
	idUint, ok := id.(uint)
	if ok {
		return idUint, nil
	}

	idInt, ok := id.(int)
	if ok {
		return uint(idInt), nil
	}

	return 0, errors.New("user id is of invalid type")
}
