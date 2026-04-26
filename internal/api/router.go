package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/wootaiklee/git-resume/internal/api/handlers"
	"github.com/wootaiklee/git-resume/internal/api/jobs"
	"github.com/wootaiklee/git-resume/internal/api/middleware"
	"github.com/wootaiklee/git-resume/internal/api/session"
	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/internal/service"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	DB             *db.DB
	Cache          *db.Cache
	JobManager     *jobs.Manager
	ResultsService *service.ResultsService
	SessionManager *session.Manager
	Logger         *slog.Logger
}

// NewRouter creates and configures the HTTP router
func NewRouter(deps *Dependencies, corsOrigins []string) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(deps.Logger))
	r.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Request-ID", "X-Session-ID", "X-API-Key"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(deps.DB)
	templatesHandler := handlers.NewTemplatesHandler()
	statsHandler := handlers.NewStatsHandler(deps.ResultsService)
	resultsHandler := handlers.NewResultsHandler(deps.ResultsService)
	exportHandler := handlers.NewExportHandler(deps.ResultsService)
	analyzeHandler := handlers.NewAnalyzeHandler(deps.JobManager)
	jobsHandler := handlers.NewJobsHandler(deps.JobManager)
	sessionsHandler := handlers.NewSessionsHandler(deps.SessionManager)

	// Health endpoints (no auth)
	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Ready)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Session management (no auth required)
		r.Post("/sessions", sessionsHandler.Create)
		r.Get("/sessions/{id}", sessionsHandler.Get)
		r.Delete("/sessions/{id}", sessionsHandler.Delete)
		r.Post("/sessions/{id}/extend", sessionsHandler.Extend)

		// Templates (no auth - public info)
		r.Get("/templates", templatesHandler.List)

		// Protected routes (require session)
		r.Group(func(r chi.Router) {
			r.Use(middleware.SessionAuth(deps.SessionManager))

			// Analysis
			r.Post("/analyze", analyzeHandler.Start)

			// Jobs
			r.Get("/jobs", jobsHandler.List)
			r.Get("/jobs/{id}", jobsHandler.Get)
			r.Delete("/jobs/{id}", jobsHandler.Cancel)

			// Results
			r.Get("/results", resultsHandler.List)
			r.Get("/results/{id}", resultsHandler.Get)

			// Export
			r.Get("/export", exportHandler.Export)

			// Stats
			r.Get("/stats", statsHandler.GetStats)
		})
	})

	return r
}
