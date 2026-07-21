package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/denizyetis/meta-rtls/internal/config"
	"github.com/denizyetis/meta-rtls/internal/modules/analysis"
	"github.com/denizyetis/meta-rtls/internal/modules/identity"
	"github.com/denizyetis/meta-rtls/internal/modules/location"
	"github.com/denizyetis/meta-rtls/internal/modules/metadata"
	"github.com/denizyetis/meta-rtls/internal/modules/rtlsconfig"
	"github.com/denizyetis/meta-rtls/internal/modules/tenant"
	"github.com/denizyetis/meta-rtls/internal/platform/auth"
	"github.com/denizyetis/meta-rtls/internal/platform/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
	mqtt   *location.MQTTWorker
	sim    *location.Simulator
}

func New(cfg *config.Config, db *sql.DB, logger *slog.Logger) (*App, error) {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.JWTTTLMinutes)

	idRepo := identity.NewRepository(db)
	idSvc := identity.NewService(idRepo, tokens)
	idHandler := identity.NewHandler(idSvc)

	tenantHandler := tenant.NewHandler(tenant.NewService(tenant.NewRepository(db)))
	rtlsHandler := rtlsconfig.NewHandler(rtlsconfig.NewService(rtlsconfig.NewRepository(db)))
	metaSvc := metadata.NewService(metadata.NewRepository(db))
	metaHandler := metadata.NewHandler(metaSvc)

	hub := location.NewHub()
	locRepo := location.NewRepository(db)
	locSvc := location.NewService(locRepo, hub, logger)
	mqttWorker := location.NewMQTTWorker(cfg.MQTTBroker, cfg.MQTTClientID, cfg.MQTTTopic, locSvc, logger)
	sim := location.NewSimulator(locRepo, locSvc, mqttWorker, logger)
	locHandler := location.NewHandler(locSvc, sim, tokens)
	analysisSvc := analysis.NewService(analysis.NewRepository(db))
	analysisHandler := analysis.NewHandler(analysisSvc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(securityHeaders())
	r.Use(requestLogger(logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{
			"service": "metartls",
			"status":  "ok",
			"env":     cfg.AppEnv,
		})
	})

	r.GET("/ready", func(c *gin.Context) {
		pingCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(pingCtx); err != nil {
			response.Fail(c, http.StatusServiceUnavailable, "oracle not ready")
			return
		}
		response.OK(c, gin.H{
			"service": "metartls",
			"status":  "ready",
			"oracle":  "up",
		})
	})

	api := r.Group("/api/v1")
	idHandler.Register(api, api.Group("", auth.Middleware(tokens)))

	protected := api.Group("", auth.Middleware(tokens))
	tenantHandler.Register(protected)
	rtlsHandler.Register(protected)
	metaHandler.Register(protected)
	locHandler.Register(api, protected)
	analysisHandler.Register(protected)

	a := &App{router: r, mqtt: mqttWorker, sim: sim}

	bootCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := idSvc.BootstrapDemoUsers(bootCtx); err != nil {
		logger.Warn("demo user bootstrap skipped or failed", "err", err)
	}
	if err := metaSvc.BootstrapDemoMetadata(bootCtx); err != nil {
		logger.Warn("demo metadata bootstrap skipped or failed", "err", err)
	}
	if err := analysisSvc.Bootstrap(bootCtx); err != nil {
		logger.Warn("analysis bootstrap skipped or failed", "err", err)
	}

	mqttWorker.Start()
	if err := sim.Start(bootCtx); err != nil {
		logger.Warn("simulator start skipped or failed", "err", err)
	}

	return a, nil
}

func (a *App) Router() http.Handler { return a.router }

func (a *App) Close() {
	if a.sim != nil {
		a.sim.Stop()
	}
	if a.mqtt != nil {
		a.mqtt.Stop()
	}
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("X-XSS-Protection", "0")
		c.Next()
	}
}

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("http",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}
