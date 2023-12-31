package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	jwtauth "shotwot_backend/pkg/auth"
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

type Users interface {
	SignUp(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	Update(ctx context.Context, input *domain.User) (*domain.User, error)
	// RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	// Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Auth interface {
	UserIdentity(token string) (*jwtauth.CustomClaims, error)
}

type Services struct {
	Users Users
	Auth  Auth
}

type Deps struct {
	Repos        *repository.Repositories
	TokenManager jwtauth.TokenManager

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	AuthClient *firebase.AuthClient
}

func NewServices(deps Deps) *Services {
	userService := NewUsersService(deps.Repos.Users, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.AuthClient)
	authService := NewAuthService(deps.TokenManager)
	return &Services{
		Users: userService,
		Auth:  authService,
	}
}
