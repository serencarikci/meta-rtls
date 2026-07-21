package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/denizyetis/meta-rtls/internal/config"
	"github.com/denizyetis/meta-rtls/internal/modules/identity"
	"github.com/denizyetis/meta-rtls/internal/modules/metadata"
	"github.com/denizyetis/meta-rtls/internal/modules/rtlsconfig"
	"github.com/denizyetis/meta-rtls/internal/modules/tenant"
	"github.com/denizyetis/meta-rtls/internal/platform/auth"
	"github.com/denizyetis/meta-rtls/internal/platform/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	cfg    *config.Config
	db     *sql.DB
	logger *slog.Logger
	router *gin.Engine
	idSvc  *identity.Service
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

	r := gin.New()
	r.Use(gin.Recovery())
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
		response.OK(c, gin.H{"service": "metartls", "status": "ok"})
	})

	api := r.Group("/api/v1")
	idHandler.Register(api, api.Group("", auth.Middleware(tokens)))

	protected := api.Group("", auth.Middleware(tokens))
	tenantHandler.Register(protected)
	rtlsHandler.Register(protected)
	metaHandler.Register(protected)

	a := &App{cfg: cfg, db: db, logger: logger, router: r, idSvc: idSvc}

	bootCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := idSvc.BootstrapDemoUsers(bootCtx); err != nil {
		logger.Warn("demo user bootstrap skipped or failed", "err", err)
	}
	if err := metaSvc.BootstrapDemoMetadata(bootCtx); err != nil {
		logger.Warn("demo metadata bootstrap skipped or failed", "err", err)
	}

	return a, nil
}

func (a *App) Router() http.Handler { return a.router }

func (a *App) Close() {
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
