package tenant

import (
	"errors"
	"net/http"

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
	protected.GET("/tenants", h.list)
	protected.GET("/tenants/:id", h.get)
	protected.POST("/tenants", h.create)
}

func (h *Handler) list(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list tenants")
		return
	}
	if items == nil {
		items = []Tenant{}
	}
	response.OK(c, items)
}

func (h *Handler) get(c *gin.Context) {
	item, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if errors.Is(err, ERR_NOT_FOUND) {
		response.Fail(c, http.StatusNotFound, "tenant not found")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to get tenant")
		return
	}
	response.OK(c, item)
}

func (h *Handler) create(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to create tenant")
		return
	}
	response.Created(c, item)
}
