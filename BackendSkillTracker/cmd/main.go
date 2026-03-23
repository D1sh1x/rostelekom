package main

import (
	"context"
	"io"
	"os"
	"skilltracker/internal/config"
	"skilltracker/internal/handler"
	"skilltracker/internal/service"
	"skilltracker/internal/storage/postgres"
	"skilltracker/internal/transport"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title SkillTracker API
// @version 1.0
// @description API Server for SkillTracker Application

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	store, err := postgres.New(cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init storage")
	}

	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("failed to create logs directory")
	}
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open log file")
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	logger := log.Logger.With().Str("app", "skilltracker").Logger()

	srv := service.New(store, logger, []byte(cfg.Auth.JWTSecret))

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}
	if err := srv.SeedAdmin(context.Background(), adminPassword); err != nil {
		logger.Error().Err(err).Msg("failed to seed admin account")
	}

	h := handler.NewHandler(srv)
	httpSrv := transport.NewServer([]byte(cfg.Auth.JWTSecret), h, cfg)
	logger.Info().Msg("Server Running")
	if err := transport.Run(httpSrv); err != nil {
		logger.Error().Err(err).Msg("server shutdown error")
	}
	logger.Info().Msg("Server Stopped")
}
