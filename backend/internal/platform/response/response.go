package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, envelope{Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, envelope{Data: data})
}

func Fail(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, envelope{Error: msg})
}
