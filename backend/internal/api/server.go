package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/dojo-harvester/backend/internal/config"
	"github.com/dojo-harvester/backend/internal/executor"
	"github.com/dojo-harvester/backend/internal/logger"
	"github.com/dojo-harvester/backend/internal/storage"
	"github.com/dojo-harvester/backend/internal/uploader"
	"github.com/dojo-harvester/backend/internal/websocket"
)

type Server struct {
	router   *gin.Engine
	srv      *http.Server
	store    *storage.Store
	cfg      *config.Config
	ws       *websocket.Manager
	exec     *executor.Executor
	uploader *uploader.Manager
}

func NewServer(cfg *config.Config, store *storage.Store) *Server {
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	wsManager := websocket.NewManager(cfg)
	execService := executor.NewExecutor(cfg)
	uploadManager := uploader.NewManager(store, cfg)

	// Middleware
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware())
	router.Use(corsMiddleware(cfg.Server.AllowOrigins))

	s := &Server{
		router:   router,
		store:    store,
		cfg:      cfg,
		ws:       wsManager,
		exec:     execService,
		uploader: uploadManager,
	}

	s.setupRoutes()

	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	return s
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	{
		api.GET("/status", s.handleStatus)
	}
	s.registerSessionRoutes(api)
	s.registerRecordingRoutes(api)
	s.registerTelemetryRoutes(api)
	s.registerDatasetRoutes(api)
	s.registerUploadRoutes(api)
	s.registerTaskRoutes(api)

	// WebSocket endpoint
	s.router.GET("/ws", s.ws.HandleWebSocket)
}

func (s *Server) handleStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) Start() error {
	logger.Get().Info("Starting HTTP server", 
		zap.String("host", s.cfg.Server.Host), 
		zap.Int("port", s.cfg.Server.Port),
	)
	
	ctx := context.Background()
	go s.ws.Run(ctx)
	
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Get().Info("Stopping HTTP server")
	return s.srv.Shutdown(ctx)
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		logger.Get().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

func corsMiddleware(allowOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
