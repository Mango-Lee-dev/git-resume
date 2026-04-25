package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wootaiklee/git-resume/internal/api"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	Long: `Start a REST API server for the Git Resume Analyzer.

The server provides endpoints for async analysis, querying results,
exporting data, and viewing statistics.

Example:
  git-resume serve
  git-resume serve --port=8080
  git-resume serve --host=0.0.0.0 --port=3000`,
	RunE: runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("host", "localhost", "server host")
	serveCmd.Flags().Int("port", 8080, "server port")
	serveCmd.Flags().StringSlice("cors-origins", []string{"*"}, "allowed CORS origins")
	serveCmd.Flags().Int("workers", 2, "number of background workers for analysis jobs")

	viper.BindPFlag("server.host", serveCmd.Flags().Lookup("host"))
	viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.cors_origins", serveCmd.Flags().Lookup("cors-origins"))
	viper.BindPFlag("server.workers", serveCmd.Flags().Lookup("workers"))
}

func runServe(cmd *cobra.Command, args []string) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := api.DefaultServerConfig()

	// Override with flags/config
	if host := viper.GetString("server.host"); host != "" {
		cfg.Host = host
	}
	if port := viper.GetInt("server.port"); port > 0 {
		cfg.Port = port
	}
	if origins := viper.GetStringSlice("server.cors_origins"); len(origins) > 0 {
		cfg.AllowedOrigins = origins
	}
	if workers := viper.GetInt("server.workers"); workers > 0 {
		cfg.WorkerCount = workers
	}

	// Database path
	cfg.DBPath = getServerDBPath()

	// API key
	cfg.ClaudeAPIKey = viper.GetString("CLAUDE_API_KEY")

	server, err := api.NewServer(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	fmt.Printf("Server started at http://%s\n", server.Addr())
	fmt.Println("Press Ctrl+C to stop")

	server.WaitForShutdown()
	return nil
}

func getServerDBPath() string {
	dbPath := viper.GetString("db")
	if dbPath == "" {
		dbPath = viper.GetString("DB_PATH")
	}
	if dbPath == "" {
		dbPath = "./data/cache.db"
	}
	return dbPath
}
