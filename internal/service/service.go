package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/firebase"
	"shotwot_backend/pkg/helper"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AccountAuthInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	IdToken  string `json:"idToken"`
	Role     int    `json:"role"`
}

type Tokens struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type Users interface {
	SignUp(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	Update(ctx context.Context, input *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error

	GetUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)

	SearchUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)

	TotalUsers(ctx context.Context) (int64, error)

	Download(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)

	// RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	// Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Admins interface {
	CreateAdmin(ctx context.Context, input AccountAuthInput) error
	SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	Update(ctx context.Context, input *domain.Admin) (*domain.Admin, error)
	GetAdmin(ctx context.Context, id string) (*domain.Admin, error)
	Delete(ctx context.Context, id string) error

	GetAllAdmins(ctx context.Context) ([]*domain.Admin, error)

	// RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	// Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Briefs interface {
	Create(ctx context.Context, input *domain.Brief) (*domain.Brief, error)

	Update(ctx context.Context, input *domain.Brief) (*domain.Brief, error)

	GetBriefs(ctx context.Context, predicate *helper.BriefPredicate) ([]*domain.Brief, error)

	DeleteBrief(ctx context.Context, id string) error
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

	Briefs Briefs

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

	WasabiS3Client *s3.Client
}

func NewServices(deps Deps) *Services {
	userService := NewUsersService(deps.Repos.Users, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.AuthClient)
	adminService := NewAdminsService(deps.Repos.Admins, deps.AdminTokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.AuthClient)

	briefService := NewBriefsService(deps.Repos.Briefs)

	authService := NewAuthService(deps.TokenManager)
	adminAuthService := NewAdminAuthService(deps.AdminTokenManager)
	return &Services{
		Users:  userService,
		Admins: adminService,
		Briefs: briefService,

		Auth:      authService,
		AdminAuth: adminAuthService,
	}
}
