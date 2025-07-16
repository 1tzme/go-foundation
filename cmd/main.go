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
	"hot-coffee/internal/repositories"
	"hot-coffee/internal/service"
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

	orderRepo := repositories.NewOrderRepository(appLogger)
	menuRepo := repositories.NewMenuRepository(appLogger)
	inventoryRepo := repositories.NewInventoryRepository(appLogger)
	aggregationRepo := repositories.NewAggregationRepository(orderRepo, menuRepo, appLogger)

	orderService := service.NewOrderService(orderRepo, appLogger)
	menuService := service.NewMenuService(menuRepo, inventoryRepo, appLogger)
	inventoryService := service.NewInventoryService(inventoryRepo, appLogger)
	aggregationService := service.NewAggregationService(aggregationRepo, appLogger)

	orderHandler := handler.NewOrderHandler(orderService, appLogger)
	menuHandler := handler.NewMenuHandler(menuService, appLogger)
	inventoryHandler := handler.NewInventoryHandler(inventoryService, appLogger)
	aggregationHandler := handler.NewAggregationHandler(aggregationService, appLogger)

	mux := http.NewServeMux()

	api := "/api/v1"

	// Aggregation routes: GET total sales, GET popular items
	mux.HandleFunc(api+"/reports/total-sales", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			aggregationHandler.GetTotalSales(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(api+"/reports/popular-items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			aggregationHandler.GetPopularItems(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Order collection routes: POST (create), GET (all)
	mux.HandleFunc(api+"/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			orderHandler.CreateOrder(w, r)
			return
		}
		if r.Method == http.MethodGet {
			orderHandler.GetAllOrders(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Order item routes: GET (by id), PUT (update), DELETE (delete)
	mux.HandleFunc(api+"/orders/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a close order request: POST /api/v1/orders/{id}/close
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/close") {
			orderHandler.CloseOrder(w, r)
			return
		}

		// Regular order operations: GET, PUT, DELETE /api/v1/orders/{id}
		if r.Method == http.MethodGet {
			orderHandler.GetOrderByID(w, r)
			return
		}
		if r.Method == http.MethodPut {
			orderHandler.UpdateOrder(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			orderHandler.DeleteOrder(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Menu collection routes: POST (create), GET (all)
	mux.HandleFunc(api+"/menu", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			menuHandler.CreateMenuItem(w, r)
			return
		}
		if r.Method == http.MethodGet {
			menuHandler.GetAllMenuItems(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Menu item routes: GET (by id), PUT (update), DELETE (delete)
	mux.HandleFunc(api+"/menu/", func(w http.ResponseWriter, r *http.Request) {
		// /api/v1/menu/{id}
		if r.Method == http.MethodGet {
			menuHandler.GetMenuItem(w, r)
			return
		}
		if r.Method == http.MethodPut {
			menuHandler.UpdateMenuItem(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			menuHandler.DeleteMenuItem(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Inventory collection routes: POST (create), GET (all)
	mux.HandleFunc(api+"/inventory", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			inventoryHandler.CreateInventoryItem(w, r)
			return
		}
		if r.Method == http.MethodGet {
			inventoryHandler.GetAllInventoryItems(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Inventory item routes: GET (by id), PUT (update), DELETE (delete)
	mux.HandleFunc(api+"/inventory/", func(w http.ResponseWriter, r *http.Request) {
		// /api/v1/inventory/{id}
		if r.Method == http.MethodGet {
			inventoryHandler.GetInventoryItem(w, r)
			return
		}
		if r.Method == http.MethodPut {
			inventoryHandler.UpdateInventoryItem(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			inventoryHandler.DeleteInventoryItem(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
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

	serverErrors := make(chan error, 1)

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		server.Addr = host + ":" + port

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

		select {
		case err := <-serverErrors:
			if strings.Contains(err.Error(), "address already in use") && i < maxRetries-1 {
				portNum := 8080 + i + 1
				port = fmt.Sprintf("%d", portNum)
				appLogger.Warn("Port already in use, trying alternative port",
					"current_port", server.Addr,
					"next_port", port)
				continue
			} else {
				appLogger.Error("Failed to start server after retries", "error", err)
				return
			}
		case <-time.After(200 * time.Millisecond):
			appLogger.Info("Server started successfully", "port", port)
		}

		break
	}

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
