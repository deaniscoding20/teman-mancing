package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"teman-mancing/backend/internal/config"
	"teman-mancing/backend/internal/delivery/http/routes"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/email"
	"teman-mancing/backend/internal/infrastructure/logger"
	"teman-mancing/backend/internal/infrastructure/security"
	"teman-mancing/backend/internal/infrastructure/storage"
)

const AppVersion = "1.0.0-phase1-foundation"

func main() {
	// 1. Load Configurations
	cfg := config.Load()

	// 2. Initialize Structured Logger
	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Starting Teman Mancing API Server",
		zap.String("version", AppVersion),
		zap.String("env", cfg.AppEnv),
		zap.Int("port", cfg.AppPort),
	)

	// 3. Initialize Database Connection (optional fallback in dev mode)
	var dbConn *database.PostgresDB
	dbConn, err = database.NewPostgresDB(cfg, log)
	if err != nil {
		if cfg.IsProduction() {
			log.Fatal("Failed to connect to PostgreSQL database", zap.Error(err))
		} else {
			log.Warn("PostgreSQL not reachable on startup, continuing in degraded dev mode", zap.Error(err))
		}
	} else {
		defer func() {
			if err := dbConn.Close(); err != nil {
				log.Error("Failed to close DB connection", zap.Error(err))
			}
		}()
	}

	// 4. Initialize MinIO Object Storage
	var minioClient *storage.MinIOClient
	minioClient, err = storage.NewMinIOClient(cfg, log)
	if err != nil {
		if cfg.IsProduction() {
			log.Fatal("Failed to initialize MinIO client", zap.Error(err))
		} else {
			log.Warn("MinIO client initialization failed, continuing in dev mode", zap.Error(err))
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		if err := minioClient.EnsureBuckets(ctx); err != nil {
			log.Warn("Could not ensure MinIO buckets on startup", zap.Error(err))
		}
		cancel()
	}

	// 5. Setup Token Service, Mailer & Repositories
	tokenService := security.NewTokenService(cfg.JWTSecret, cfg.JWTExpiration)
	mailer := email.NewMailer(email.Config{
		Host:        cfg.SMTPHost,
		Port:        cfg.SMTPPort,
		User:        cfg.SMTPUser,
		Password:    cfg.SMTPPassword,
		UseSSL:      cfg.SMTPUseSSL,
		SenderEmail: cfg.SMTPSender,
		SenderName:  cfg.SMTPSenderName,
		BaseURL:     cfg.AppBaseURL,
	}, log)

	// 6. Setup Gin Router & Middlewares
	router := routes.SetupRouter(routes.RouterDeps{
		Config:       cfg,
		Logger:       log,
		DB:           dbConn,
		MinIO:        minioClient,
		Mailer:       mailer,
		TokenService: tokenService,
		Version:      AppVersion,
	})

	// 6. Setup HTTP Server
	serverAddr := fmt.Sprintf(":%d", cfg.AppPort)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 7. Run Server in Goroutine
	go func() {
		log.Info("HTTP server listening", zap.String("addr", serverAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server failed to listen and serve", zap.Error(err))
		}
	}()

	// 8. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down HTTP server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Teman Mancing API Server exited cleanly")
}
