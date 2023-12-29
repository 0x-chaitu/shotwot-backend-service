package firebase

import (
	"context"
	"shotwot_backend/pkg/logger"

	firebase "firebase.google.com/go"
	fireauth "firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type AuthClient struct {
	*fireauth.Client
}

func NewAuthClient() (*AuthClient, error) {
	opt := option.WithCredentialsFile("/home/chaitu/Projects/showot/shotwot_backend/shotwot_test_firebase.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		logger.Errorf("Failed to create Firebase auth client: %v", err)
	}

	return &AuthClient{authClient}, nil
}
