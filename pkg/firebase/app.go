package firebase

import (
	"context"
	"shotwot_backend/pkg/logger"

	firebase "firebase.google.com/go/v4"
	fireauth "firebase.google.com/go/v4/auth"
)

type AuthClient struct {
	*fireauth.Client
}

func NewAuthClient() (*AuthClient, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		logger.Errorf("Failed to create Firebase auth client: %v", err)
	}

	return &AuthClient{authClient}, nil
}
