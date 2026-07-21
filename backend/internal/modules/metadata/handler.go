package metadata

import (
	"github.com/denizyetis/meta-rtls/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Register(protected *gin.RouterGroup) {
	protected.GET("/metadata/definitions", h.listDefinitions)
}

func (h *Handler) listDefinitions(c *gin.Context) {
	response.OK(c, gin.H{
		"items":   []any{},
		"message": "Metadata engine comes in Phase 2. Schema tables are already in Oracle.",
	})
}

func HealthPayload() gin.H {
	return gin.H{"service": "metartls", "status": "ok"}
}
