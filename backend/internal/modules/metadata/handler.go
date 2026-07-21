package metadata

import (
	"errors"
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
	protected.GET("/metadata/definitions", h.listDefinitions)
	protected.POST("/metadata/definitions", h.createDefinition)
	protected.GET("/metadata/definitions/:id", h.getDefinition)
	protected.GET("/metadata/definitions/:id/versions", h.listVersions)
	protected.POST("/metadata/definitions/:id/versions", h.createVersion)
	protected.GET("/metadata/versions/:versionId/fields", h.listFields)
	protected.POST("/metadata/versions/:versionId/fields", h.createField)
	protected.POST("/metadata/validate", h.validate)
	protected.GET("/metadata/features", h.listFeatures)
}

func (h *Handler) listDefinitions(c *gin.Context) {
	items, err := h.svc.ListDefinitions(c.Request.Context(), auth.TenantID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list definitions")
		return
	}
	response.OK(c, items)
}

func (h *Handler) createDefinition(c *gin.Context) {
	var req CreateDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.CreateDefinition(c.Request.Context(), auth.TenantID(c), req)
	if errors.Is(err, ERR_INVALID_ENTITY) {
		response.Fail(c, http.StatusBadRequest, "entityType must be ASSET, PERSON, TAG, ZONE or DEVICE")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to create definition")
		return
	}
	response.Created(c, item)
}

func (h *Handler) getDefinition(c *gin.Context) {
	item, err := h.svc.GetDefinition(c.Request.Context(), auth.TenantID(c), c.Param("id"))
	if errors.Is(err, ERR_NOT_FOUND) {
		response.Fail(c, http.StatusNotFound, "definition not found")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to get definition")
		return
	}
	response.OK(c, item)
}

func (h *Handler) listVersions(c *gin.Context) {
	items, err := h.svc.ListVersions(c.Request.Context(), auth.TenantID(c), c.Param("id"))
	if errors.Is(err, ERR_NOT_FOUND) {
		response.Fail(c, http.StatusNotFound, "definition not found")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list versions")
		return
	}
	response.OK(c, items)
}

func (h *Handler) createVersion(c *gin.Context) {
	var req CreateVersionRequest
	_ = c.ShouldBindJSON(&req)
	item, err := h.svc.CreateVersion(c.Request.Context(), auth.TenantID(c), c.Param("id"), req)
	if errors.Is(err, ERR_NOT_FOUND) {
		response.Fail(c, http.StatusNotFound, "definition not found")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to create version")
		return
	}
	response.Created(c, item)
}

func (h *Handler) listFields(c *gin.Context) {
	items, err := h.svc.ListFields(c.Request.Context(), auth.TenantID(c), c.Param("versionId"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list fields")
		return
	}
	response.OK(c, items)
}

func (h *Handler) createField(c *gin.Context) {
	var req CreateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.CreateField(c.Request.Context(), auth.TenantID(c), c.Param("versionId"), req)
	if errors.Is(err, ERR_INVALID_TYPE) {
		response.Fail(c, http.StatusBadRequest, "dataType is not allowed")
		return
	}
	if errors.Is(err, ERR_NOT_FOUND) {
		response.Fail(c, http.StatusNotFound, "version not found")
		return
	}
	if errors.Is(err, ERR_BAD_REQUEST) {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to create field")
		return
	}
	response.Created(c, item)
}

func (h *Handler) validate(c *gin.Context) {
	var req ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.Validate(c.Request.Context(), auth.TenantID(c), req)
	if errors.Is(err, ERR_NOT_FOUND) {
		response.Fail(c, http.StatusNotFound, "definition not found")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "validation failed")
		return
	}
	response.OK(c, item)
}

func (h *Handler) listFeatures(c *gin.Context) {
	items, err := h.svc.ListFeatures(c.Request.Context(), auth.TenantID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list features")
		return
	}
	response.OK(c, items)
}
