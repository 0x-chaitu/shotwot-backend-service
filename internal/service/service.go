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
	IdToken  string `json:"idToken"`
}

type Tokens struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type Users interface {
	SignUp(ctx context.Context, input string) (*Tokens, error)
	SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	Update(ctx context.Context, input *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	// RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	// Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Admins interface {
	// SignUp(ctx context.Context, input string) (*Tokens, error)
	SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	// Update(ctx context.Context, input *domain.User) (*domain.User, error)
	// GetUser(ctx context.Context, id string) (*domain.User, error)
	// Delete(ctx context.Context, id string) error
	// RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	// Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Auth interface {
	UserIdentity(token string) (*jwtauth.CustomClaims, error)
}

type AdminAuth interface {
	AdminIdentity(token string) (*jwtauth.CustomAdminClaims, error)
}

type Services struct {
	Users  Users
	Admins Admins

	Auth      Auth
	AdminAuth AdminAuth
}

type Deps struct {
	Repos             *repository.Repositories
	TokenManager      jwtauth.TokenManager
	AdminTokenManager jwtauth.AdminTokenManager

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	AuthClient *firebase.AuthClient
}

func NewServices(deps Deps) *Services {
	userService := NewUsersService(deps.Repos.Users, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.AuthClient)
	authService := NewAuthService(deps.TokenManager)
	adminAuthService := NewAdminAuthService(deps.AdminTokenManager)
	return &Services{
		Users:     userService,
		Auth:      authService,
		AdminAuth: adminAuthService,
	}
}
