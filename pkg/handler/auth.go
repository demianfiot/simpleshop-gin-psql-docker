package handler

import (
	"prac/todo"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, "invalid input body") // err.Error()
		return
	}
	ctx := c.Request.Context()
	id, err := h.services.Authorization.CreateUser(ctx, input)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"id": id,
	})

}

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var input SignInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	ctx := c.Request.Context()
	token, err := h.services.Authorization.GenerateToken(ctx, input.Email, input.Password)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"token": token,
	})

}
