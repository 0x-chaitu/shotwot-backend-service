package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"shotwot_backend/internal/config"
	delivery "shotwot_backend/internal/delivery/http"
	"shotwot_backend/internal/repository"
	"shotwot_backend/internal/server"
	"shotwot_backend/internal/service"
	"shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/database/mongodb"
	"shotwot_backend/pkg/logger"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)

		return
	}

	mongoClient, err := mongodb.NewClient("mongodb+srv://chaitu:chaitu@cluster0.4cmhaeu.mongodb.net/?retryWrites=true&w=majority")
	if err != nil {
		logger.Error(err)

		return
	}

	db := mongoClient.Database("shotwot")

	repos := repository.NewRepositories(db)
	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	services := service.NewServices(
		service.Deps{
			Repos:           repos,
			TokenManager:    tokenManager,
			AccessTokenTTL:  cfg.Auth.JWT.AccessTokenTTL,
			RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL,
		})
	handlers := delivery.NewHandler(services)

	srv := server.NewServer(cfg, handlers.Init(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
