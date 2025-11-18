package handler

import (
	"errors"
	"net/http"
	"prac/pkg/handler/dto"
	"prac/todo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(c *gin.Context) {
	var input dto.CreateUserInput
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

func (h *Handler) GetAllUsers(c *gin.Context) {

	ctx := c.Request.Context()
	users, err := h.services.User.GetAllUsers(ctx)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

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
