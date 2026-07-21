package rtlsconfig

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
	protected.GET("/sites", h.listSites)
	protected.POST("/sites", h.createSite)
	protected.GET("/buildings", h.listBuildings)
	protected.GET("/floors", h.listFloors)
	protected.GET("/floors/:floorId/zones", h.listZones)
}

func (h *Handler) listSites(c *gin.Context) {
	items, err := h.svc.ListSites(c.Request.Context(), auth.TenantID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list sites")
		return
	}
	if items == nil {
		items = []Site{}
	}
	response.OK(c, items)
}

func (h *Handler) createSite(c *gin.Context) {
	var req CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.CreateSite(c.Request.Context(), auth.TenantID(c), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to create site")
		return
	}
	response.Created(c, item)
}

func (h *Handler) listBuildings(c *gin.Context) {
	items, err := h.svc.ListBuildings(c.Request.Context(), auth.TenantID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list buildings")
		return
	}
	if items == nil {
		items = []Building{}
	}
	response.OK(c, items)
}

func (h *Handler) listFloors(c *gin.Context) {
	items, err := h.svc.ListFloors(c.Request.Context(), auth.TenantID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list floors")
		return
	}
	if items == nil {
		items = []Floor{}
	}
	response.OK(c, items)
}

func (h *Handler) listZones(c *gin.Context) {
	items, err := h.svc.ListZones(c.Request.Context(), auth.TenantID(c), c.Param("floorId"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list zones")
		return
	}
	if items == nil {
		items = []Zone{}
	}
	response.OK(c, items)
}
