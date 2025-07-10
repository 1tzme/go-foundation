package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"hot-coffee/internal/handler"
	"hot-coffee/pkg/logger"
)

// loadEnvFile loads environment variables from .env file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

func main() {
	envErr := loadEnvFile(".env")

	loggerConfig := logger.Config{
		Level:         getLogLevel(),
		Format:        getEnv("LOG_FORMAT", "json"),
		Output:        getEnv("LOG_OUTPUT", "stdout"),
		EnableCaller:  getEnv("LOG_ENABLE_CALLER", "true") == "true",
		EnableColors:  getEnv("LOG_ENABLE_COLORS", "false") == "true",
		Environment:   getEnv("ENVIRONMENT", "development"),
		EnableMetrics: true,
		SensitiveKeys: []string{"password", "token", "secret", "key", "authorization"}, //abstract for now
	}

	appLogger := logger.New(loggerConfig)

	if envErr != nil {
		appLogger.Warn("Failed to load .env file", "error", envErr)
	} else {
		appLogger.Debug(".env file loaded successfully")
	}

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

	// Temporary placeholder for inventoryService until real implementation is available
	type inventoryServicePlaceholder struct{}
	var inventoryService interface{} = &inventoryServicePlaceholder{}

	// TODO: Initialize handlers with logger
	// orderHandler := handler.NewOrderHandler(orderService, appLogger)
	// menuHandler := handler.NewMenuHandler(menuService, appLogger)

	inventoryHandler := handler.NewInventoryHandler(inventoryService, appLogger)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"hot-coffee"}`))
	})

	// TODO: Add API routes
	api := "/api/v1"
	// mux.HandleFunc(api+"/orders", orderHandler.CreateOrder)
	// mux.HandleFunc(api+"/menu", menuHandler.GetMenu)
	mux.HandleFunc(api+"/inventory", func(w http.ResponseWriter, r *http.Request) {
		inventoryHandler.GetInventory(w, r)
	})

	handler := appLogger.HTTPMiddleware(mux)

	initialPort := getEnv("PORT", "8080")
	host := getEnv("HOST", "localhost")

	port := initialPort

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Create error channel to handle server errors
	serverErrors := make(chan error, 1)

	// Try to start server, with fallback ports if needed
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		// Update server address with current port
		server.Addr = host + ":" + port

		// Start the server in a goroutine
		go func() {
			appLogger.Info("Starting HTTP server",
				"host", host,
				"port", port,
				"address", server.Addr)

			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				appLogger.Error("Server error", "error", err)
				serverErrors <- err
			}
		}()

		// Wait a moment to see if the server starts successfully
		select {
		case err := <-serverErrors:
			// If port is in use, try the next port
			if strings.Contains(err.Error(), "address already in use") && i < maxRetries-1 {
				portNum := 8080 + i + 1
				port = fmt.Sprintf("%d", portNum)
				appLogger.Warn("Port already in use, trying alternative port",
					"current_port", server.Addr,
					"next_port", port)
				continue
			} else {
				// Other error or out of retries
				appLogger.Error("Failed to start server after retries", "error", err)
				return
			}
		case <-time.After(200 * time.Millisecond):
			// Server started successfully
			appLogger.Info("Server started successfully", "port", port)
			break
		}

		// If we get here, the server started successfully
		break
	}

	// Wait for shutdown or server error
	select {
	case err := <-serverErrors:
		appLogger.Error("Could not start server", "error", err)
		return
	default:
		setupGracefulShutdown(server, appLogger)
	}
}

// setupGracefulShutdown handles graceful server shutdown
func setupGracefulShutdown(server *http.Server, logger *logger.Logger) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info("Received shutdown signal", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
