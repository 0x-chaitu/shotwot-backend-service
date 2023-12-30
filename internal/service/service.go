package service

import (
	"context"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/firebase"
	"time"
)

type AccountAuthInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Tokens struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type Accounts interface {
	SignUp(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	// RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	// Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Services struct {
	Accounts Accounts
}

type Deps struct {
	Repos        *repository.Repositories
	TokenManager auth.TokenManager

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	AuthClient *firebase.AuthClient
}

func NewServices(deps Deps) *Services {
	accountService := NewAccountsService(deps.Repos.Accounts, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.AuthClient)
	return &Services{
		Accounts: accountService,
	}
}
