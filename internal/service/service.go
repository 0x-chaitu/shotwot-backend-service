package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/firebase"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/s3"
	"time"
)

type AccountAuthInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	IdToken  string `json:"idToken"`
	Role     int    `json:"role"`
}

type AuthResponse struct {
	*Tokens
	User *domain.User `json:"user"`
}

type Tokens struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type Users interface {
	SignUp(ctx context.Context, input AccountAuthInput) (*AuthResponse, error)
	SignIn(ctx context.Context, input AccountAuthInput) (*AuthResponse, error)
	Update(ctx context.Context, input *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	GetOrCreateByPhone(ctx context.Context, user *domain.User) (*AuthResponse, error)

	GetUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)

	SearchUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)

	TotalUsers(ctx context.Context) (int64, error)

	Download(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)

	// RefreshAuthResponse(ctx context.Context, refreshToken string) (AuthResponse, error)
	// Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Admins interface {
	CreateAdmin(ctx context.Context, input AccountAuthInput) error
	SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error)
	Update(ctx context.Context, input *domain.Admin) (*domain.Admin, error)
	GetAdmin(ctx context.Context, id string) (*domain.Admin, error)
	Delete(ctx context.Context, id string) error

	GetAllAdmins(ctx context.Context) ([]*domain.Admin, error)

	// RefreshAuthResponse(ctx context.Context, refreshToken string) (AuthResponse, error)
	// Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Briefs interface {
	Create(ctx context.Context, input *domain.BriefInput) (*domain.BriefRes, error)

	Update(ctx context.Context, input *domain.BriefInput) (*domain.BriefRes, error)

	GetBriefs(ctx context.Context, predicate *helper.BriefPredicate) ([]*domain.Brief, error)

	GetBrief(ctx context.Context, id string) (*domain.Brief, error)

	DeleteBrief(ctx context.Context, id string) error

	CreateDraft(ctx context.Context, brief *domain.BriefInput) (*domain.BriefRes, error)
}

type BriefApplications interface {
	Create(ctx context.Context, input domain.BriefApplication) (*domain.BriefApplication, error)

	UpdateBriefApplication(ctx context.Context, input *domain.BriefApplication) (*domain.BriefApplication, error)

	GetBriefApplications(ctx context.Context, id string) ([]*domain.BriefApplication, error)

	GetBriefApplication(ctx context.Context, id string) (*domain.UserBriefAppliedDetails, error)

	GetUserBriefApplications(ctx context.Context, id string) ([]*domain.UserBriefAppliedDetails, error)
}

type Assets interface {
	Create(ctx context.Context, asset *domain.AssetInput) (*domain.AssetRes, error)
}

type SavedBriefs interface {
	CreateOrUpdate(ctx context.Context, input *domain.SavedBriefInput) (*domain.SavedBriefRes, error)

	GetSavedBriefs(ctx context.Context, userId string) ([]*domain.SavedBrief, error)
}

type Auth interface {
	UserIdentity(token string) (*jwtauth.CustomClaims, error)
}

type AdminAuth interface {
	AdminIdentity(token string) (*jwtauth.CustomAdminClaims, error)
}

type MasterClass interface {
	CreatePlaylist(ctx context.Context, input *domain.PlaylistInput) (*domain.PlaylistResp, error)
	GetPlaylists(ctx context.Context, predicate *helper.PlaylistPredicate) ([]*domain.Playlist, error)
}

type Services struct {
	Users  Users
	Admins Admins

	Briefs            Briefs
	BriefApplications BriefApplications
	SavedBriefs       SavedBriefs
	MasterClass       MasterClass
	Assets            Assets

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

	WasabiS3Client *s3.S3Client
}

func NewServices(deps Deps) *Services {
	userService := NewUsersService(deps.Repos.Users, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.AuthClient)
	adminService := NewAdminsService(deps.Repos.Admins, deps.AdminTokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.AuthClient)

	briefService := NewBriefsService(deps.Repos.Briefs, deps.WasabiS3Client)
	briefApplications := NewBriefApplicationsService(deps.Repos.BriefApplications, deps.WasabiS3Client)
	savedBriefs := NewSavedBriefsService(deps.Repos.SavedBriefs)
	masterClass := NewMasterClassService(deps.Repos.MasterClass, deps.WasabiS3Client)
	assetServie := NewAssetsService(deps.Repos.Assets, deps.WasabiS3Client)

	authService := NewAuthService(deps.TokenManager)
	adminAuthService := NewAdminAuthService(deps.AdminTokenManager)
	return &Services{
		Users:             userService,
		Admins:            adminService,
		Briefs:            briefService,
		BriefApplications: briefApplications,
		SavedBriefs:       savedBriefs,
		MasterClass:       masterClass,
		Assets:            assetServie,

		Auth:      authService,
		AdminAuth: adminAuthService,
	}
}
