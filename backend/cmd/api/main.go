package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rafif/healy-backend/internal/delivery/websocket"
	deliveryHttp "github.com/rafif/healy-backend/internal/delivery/http"
	"github.com/rafif/healy-backend/internal/repository/postgres"
	"github.com/rafif/healy-backend/internal/usecase"
	"github.com/rafif/healy-backend/pkg/config"
	"github.com/rafif/healy-backend/pkg/database"
	"github.com/rafif/healy-backend/pkg/jwt"
)

func main() {
	// 1. Initialize context
	ctx := context.Background()

	// 2. Load Config
	cfg := config.LoadConfig()

	// 3. Initialize Database Pool
	if cfg.DatabaseURL == "" {
		log.Println("DATABASE_URL is empty, will probably fail if no default is provided")
	}
	dbPool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()

	// 4. Initialize Repositories
	// Inject the concrete PostgreSQL implementations here, 
	// satisfying the interfaces expected by the usecase layer.
	telemetryRepo := postgres.NewTelemetryRepository(dbPool)
	userRepo := postgres.NewUserRepository(dbPool)

	// 5. Initialize Services/Utils
	tokenGenerator := jwt.NewJWTGenerator(cfg)

	// 6. Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// 7. Initialize Usecases
	// Inject the hub's broadcast channel into the telemetry usecase
	telemetryUsecase := usecase.NewTelemetryUsecase(telemetryRepo, hub.Broadcast)
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenGenerator)

	// 8. Initialize Delivery (HTTP Router)
	router := deliveryHttp.SetupRouter(cfg, hub, telemetryUsecase, authUsecase)

	// 9. Setup Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: router,
	}

	// 10. Graceful Shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
