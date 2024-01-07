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
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/database/mongodb"
	"shotwot_backend/pkg/firebase"
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

	mongoClient, err := mongodb.NewClient("mongodb+srv://goshotwot:Tokyo.3619@shotwottest.rpgjdmf.mongodb.net/?retryWrites=true&w=majority")
	if err != nil {
		logger.Error(err)
		return
	}

	db := mongoClient.Database("shotwot")

	repos := repository.NewRepositories(db)
	tokenManager, err := jwtauth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}
	adminTokenManager, err := jwtauth.NewAdminManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}
	authClient, err := firebase.NewAuthClient()
	if err != nil {
		logger.Error(err)
		return
	}
	services := service.NewServices(
		service.Deps{
			Repos:             repos,
			TokenManager:      tokenManager,
			AdminTokenManager: adminTokenManager,
			AccessTokenTTL:    cfg.Auth.JWT.AccessTokenTTL,
			RefreshTokenTTL:   cfg.Auth.JWT.RefreshTokenTTL,
			AuthClient:        authClient,
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
