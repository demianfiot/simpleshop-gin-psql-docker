package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type myerror struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, massage string) {
	logrus.Errorf(massage)
	c.AbortWithStatusJSON(statusCode, myerror{Message: massage})
}
