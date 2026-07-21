package location

import (
	"net/http"

	"github.com/denizyetis/meta-rtls/internal/platform/auth"
	"github.com/denizyetis/meta-rtls/internal/platform/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	svc      *Service
	sim      *Simulator
	tokens   *auth.TokenService
	upgrader websocket.Upgrader
}

func NewHandler(svc *Service, sim *Simulator, tokens *auth.TokenService) *Handler {
	return &Handler{
		svc:    svc,
		sim:    sim,
		tokens: tokens,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *Handler) Register(api *gin.RouterGroup, protected *gin.RouterGroup) {
	protected.GET("/locations/latest", h.latest)
	protected.GET("/simulator/status", h.simStatus)
	protected.POST("/simulator/start", h.simStart)
	protected.POST("/simulator/stop", h.simStop)

	api.GET("/ws/locations", h.wsLocations)
}

func (h *Handler) latest(c *gin.Context) {
	items := h.svc.LatestForTenant(auth.TenantID(c))
	if items == nil {
		items = []LivePosition{}
	}
	response.OK(c, items)
}

func (h *Handler) simStatus(c *gin.Context) {
	response.OK(c, gin.H{"running": h.sim.Running()})
}

func (h *Handler) simStart(c *gin.Context) {
	if err := h.sim.Start(c.Request.Context()); err != nil {
		response.Fail(c, http.StatusInternalServerError, "could not start simulator")
		return
	}
	response.OK(c, gin.H{"running": true})
}

func (h *Handler) simStop(c *gin.Context) {
	h.sim.Stop()
	response.OK(c, gin.H{"running": false})
}

func (h *Handler) wsLocations(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		response.Fail(c, http.StatusUnauthorized, "token query required")
		return
	}
	claims, err := h.tokens.Parse(tokenStr)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "invalid token")
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	h.svc.hub.Add(conn, claims.TenantID)

	for _, pos := range h.svc.LatestForTenant(claims.TenantID) {
		_ = conn.WriteJSON(pos)
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.svc.hub.Remove(conn)
			return
		}
	}
}
