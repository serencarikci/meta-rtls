package services

import (
	"net/http"
	"strings"

	"github.com/denizyetis/meta-rtls/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Services
}

func NewHandler(svc *Services) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoot(router *gin.Engine) {
	router.Use(h.funcQueryMiddleware())
}

func (h *Handler) funcQueryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		funcName := strings.ToLower(strings.TrimSpace(c.Query("func")))
		if funcName == "" {
			c.Next()
			return
		}

		switch funcName {
		case "getversion":
			response.OK(c, h.svc.GetVersion())
		case "getconfig":
			response.OK(c, h.svc.GetConfig())
		default:
			response.Fail(c, http.StatusBadRequest, "unknown func")
		}
		c.Abort()
	}
}
