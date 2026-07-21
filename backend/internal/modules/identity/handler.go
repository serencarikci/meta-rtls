package identity

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

func (h *Handler) Register(public, protected *gin.RouterGroup) {
	public.POST("/auth/login", h.login)
	protected.GET("/auth/me", h.me)
}

func (h *Handler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.svc.Login(c.Request.Context(), req)
	if errors.Is(err, ERR_INVALID_CREDENTIALS) {
		response.Fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "login failed")
		return
	}
	response.OK(c, out)
}

func (h *Handler) me(c *gin.Context) {
	out, err := h.svc.Me(c.Request.Context(), auth.UserID(c))
	if errors.Is(err, ERR_NOT_FOUND) {
		response.Fail(c, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to load profile")
		return
	}
	response.OK(c, out)
}
