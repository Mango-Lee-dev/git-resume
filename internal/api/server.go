package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wootaiklee/git-resume/internal/api/jobs"
	"github.com/wootaiklee/git-resume/internal/api/session"
	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/internal/service"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	AllowedOrigins []string
	DBPath         string
	ClaudeAPIKey   string
	WorkerCount    int
}

// DefaultServerConfig returns default configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:           "localhost",
		Port:           8080,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    120 * time.Second,
		AllowedOrigins: []string{"*"},
		DBPath:         "./data/cache.db",
		WorkerCount:    2,
	}
}

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	deps       *Dependencies
	logger     *slog.Logger
}

// NewServer creates a new server instance
func NewServer(cfg *ServerConfig, logger *slog.Logger) (*Server, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	// Initialize database
	database, err := db.New(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	cache := db.NewCache(database)

	// Initialize services
	analyzer := service.NewAnalyzer(cache, cfg.ClaudeAPIKey)
	resultsService := service.NewResultsService(cache)

	// Initialize job manager
	workerCount := cfg.WorkerCount
	if workerCount <= 0 {
		workerCount = 2
	}
	jobManager := jobs.NewManager(analyzer, workerCount)

	// Initialize session manager (24 hour TTL)
	sessionManager := session.NewManager(24 * time.Hour)

	deps := &Dependencies{
		DB:             database,
		Cache:          cache,
		JobManager:     jobManager,
		ResultsService: resultsService,
		SessionManager: sessionManager,
		Logger:         logger,
	}

	router := NewRouter(deps, cfg.AllowedOrigins)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		deps:       deps,
		logger:     logger,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("starting server", "addr", s.httpServer.Addr)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Error("server error", "error", err)
		}
	}()

	return nil
}

// Addr returns the server address
func (s *Server) Addr() string {
	return s.httpServer.Addr
}

// WaitForShutdown waits for shutdown signal and gracefully stops the server
func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.Shutdown()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	s.logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("server shutdown error", "error", err)
	}

	// Shutdown job manager
	s.deps.JobManager.Shutdown()

	// Close database
	s.deps.DB.Close()

	s.logger.Info("server stopped")
}
