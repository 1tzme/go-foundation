package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hot-coffee/pkg/logger"
)

func main() {
	// Initialize logger
	loggerConfig := logger.Config{
		Level:         getLogLevel(),
		Format:        getEnv("LOG_FORMAT", "json"),
		Output:        getEnv("LOG_OUTPUT", "stdout"),
		EnableCaller:  getEnv("LOG_ENABLE_CALLER", "true") == "true",
		EnableColors:  getEnv("LOG_ENABLE_COLORS", "false") == "true",
		Environment:   getEnv("ENVIRONMENT", "development"),
		EnableMetrics: true,
		SensitiveKeys: []string{"password", "token", "secret", "key", "authorization"},
	}

	appLogger := logger.New(loggerConfig)
	appLogger.Info("Starting Hot Coffee application",
		"environment", loggerConfig.Environment,
		"log_level", loggerConfig.Level)

	// TODO: Initialize repositories with logger
	// orderRepo := dal.NewOrderRepository(appLogger)
	// menuRepo := dal.NewMenuRepository(appLogger)
	// inventoryRepo := dal.NewInventoryRepository(appLogger)

	// TODO: Initialize services with logger
	// orderService := service.NewOrderService(orderRepo, inventoryRepo, appLogger)
	// menuService := service.NewMenuService(menuRepo, appLogger)
	// inventoryService := service.NewInventoryService(inventoryRepo, appLogger)

	// TODO: Initialize handlers with logger
	// orderHandler := handler.NewOrderHandler(orderService, appLogger)
	// menuHandler := handler.NewMenuHandler(menuService, appLogger)
	// inventoryHandler := handler.NewInventoryHandler(inventoryService, appLogger)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Simple JSON response
		w.Write([]byte(`{"status":"healthy","service":"hot-coffee"}`))
	})

	// TODO: Add API routes
	// api := "/api/v1"
	// mux.HandleFunc(api+"/orders", orderHandler.CreateOrder)
	// mux.HandleFunc(api+"/menu", menuHandler.GetMenu)
	// mux.HandleFunc(api+"/inventory", inventoryHandler.GetInventory)

	// Setup server with logging middleware
	handler := appLogger.HTTPMiddleware(mux)

	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "localhost")

	server := &http.Server{
		Addr:         host + ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("Starting HTTP server",
			"host", host,
			"port", port,
			"address", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Setup graceful shutdown
	setupGracefulShutdown(server, appLogger)
}

// setupGracefulShutdown handles graceful server shutdown
func setupGracefulShutdown(server *http.Server, logger *logger.Logger) {
	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)

	// Register the channel to receive specific signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-quit
	logger.Info("Received shutdown signal", "signal", sig.String())

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	logger.Info("Shutting down server gracefully...")

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		return
	}

	logger.Info("Server shutdown completed successfully")
}

// getLogLevel returns the log level from environment variable
func getLogLevel() logger.LogLevel {
	level := getEnv("LOG_LEVEL", "info")
	switch level {
	case "debug":
		return logger.LevelDebug
	case "info":
		return logger.LevelInfo
	case "warn":
		return logger.LevelWarn
	case "error":
		return logger.LevelError
	default:
		return logger.LevelInfo
	}
}

// getEnv returns environment variable value or default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
