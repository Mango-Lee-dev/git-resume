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
	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/internal/service"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	DB             *db.DB
	Cache          *db.Cache
	JobManager     *jobs.Manager
	ResultsService *service.ResultsService
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
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Request-ID"},
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

	// Health endpoints (no auth)
	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Ready)

	// API routes
	r.Route("/api", func(r chi.Router) {
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

		// Templates
		r.Get("/templates", templatesHandler.List)

		// Stats
		r.Get("/stats", statsHandler.GetStats)
	})

	return r
}
