package analysis

import (
	"net/http"

	"github.com/denizyetis/meta-rtls/internal/platform/auth"
	"github.com/denizyetis/meta-rtls/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(protected *gin.RouterGroup) {
	protected.GET("/analysis/requirements", h.listRequirements)
	protected.GET("/analysis/compare", h.compare)
	protected.POST("/analysis/impact", h.impact)
	protected.GET("/analysis/change-requests", h.listChangeRequests)
}

func (h *Handler) listRequirements(c *gin.Context) {
	items, err := h.svc.ListRequirements(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list requirements")
		return
	}
	response.OK(c, items)
}

func (h *Handler) compare(c *gin.Context) {
	item, err := h.svc.CompareProfiles(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to compare profiles")
		return
	}
	response.OK(c, item)
}

func (h *Handler) impact(c *gin.Context) {
	var req ImpactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.AnalyzeImpact(c.Request.Context(), auth.TenantID(c), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "impact analysis failed")
		return
	}
	response.OK(c, item)
}

func (h *Handler) listChangeRequests(c *gin.Context) {
	items, err := h.svc.ListChangeRequests(c.Request.Context(), auth.TenantID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list change requests")
		return
	}
	response.OK(c, items)
}
