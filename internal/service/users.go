package service

import (
	"context"
	"fmt"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/firebase"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"
	"time"

	fireauth "firebase.google.com/go/v4/auth"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type UsersService struct {
	repo               repository.Users
	tokenManager       jwtauth.TokenManager
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
	firebaseAuthClient *firebase.AuthClient
}

func (a AccountAuthInput) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required, validation.Length(6, 64)),
	)
}

func NewUsersService(repo repository.Users, tokenManager jwtauth.TokenManager, accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration, firebaseAuthClient *firebase.AuthClient) *UsersService {
	return &UsersService{
		repo:               repo,
		tokenManager:       tokenManager,
		accessTokenTTL:     accessTokenTTL,
		refreshTokenTTL:    refreshTokenTTL,
		firebaseAuthClient: firebaseAuthClient,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input AccountAuthInput) (*AuthResponse, error) {
	token, err := s.firebaseAuthClient.VerifyIDToken(ctx, input.IdToken)
	if err != nil {
		if fireauth.IsEmailAlreadyExists(err) {
			return nil, domain.ErrAccountAlreadyExists
		}
		logger.Errorf("firebase user creation error %v", err)
		return nil, err
	}
	account := &domain.User{
		Email:    token.Claims["email"].(string),
		UserId:   token.UID,
		Created:  time.Now(),
		UserName: token.UID,
	}

	user, err := s.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}
	tokens, err := s.createSession(ctx, user.Id.Hex())
	return &AuthResponse{
		User:   account,
		Tokens: tokens,
	}, err

}

func (s *UsersService) SignIn(ctx context.Context, input AccountAuthInput) (*AuthResponse, error) {
	token, err := s.firebaseAuthClient.VerifyIDToken(ctx, input.IdToken)
	if err != nil {
		logger.Errorf("firebase user signin error %v", err)
		return nil, err
	}
	account := &domain.User{
		Email:    token.Claims["email"].(string),
		UserId:   token.UID,
		Created:  time.Now(),
		UserName: token.UID,
	}
	user, err := s.repo.GetOrCreate(ctx, account)
	if err != nil {
		return nil, err
	}
	tokens, err := s.createSession(ctx, user.Id.Hex())
	return &AuthResponse{
		User:   user,
		Tokens: tokens,
	}, err
}

func (s *UsersService) GetOrCreateByPhone(ctx context.Context, user *domain.User) (*AuthResponse, error) {
	user, err := s.repo.GetOrCreateByPhone(ctx, user)
	if err != nil {
		return nil, err
	}
	tokens, err := s.createSession(ctx, user.Id.Hex())
	return &AuthResponse{
		User:   user,
		Tokens: tokens,
	}, err
}

func (s *UsersService) Download(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error) {
	return s.repo.Download(ctx, predicate)
}
func (s *UsersService) Update(ctx context.Context, input *domain.User) (*domain.User, error) {
	user, err := s.repo.Update(ctx, input)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return user, nil
}

func (s *UsersService) Delete(ctx context.Context, id string) error {
	return s.firebaseAuthClient.DeleteUser(ctx, id)
}

func (s *UsersService) GetUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error) {
	return s.repo.GetUsers(ctx, predicate)
}

func (s *UsersService) SearchUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error) {
	return s.repo.SearchUsers(ctx, predicate)
}

func (s *UsersService) TotalUsers(ctx context.Context) (int64, error) {
	return s.repo.TotalUsers(ctx)
}

func (s *UsersService) createSession(_ context.Context, userId string) (*Tokens, error) {
	token, err := s.tokenManager.NewJWT(fmt.Sprint(userId), s.accessTokenTTL)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	refreshtoken, err := s.tokenManager.NewJWT(fmt.Sprint(userId), s.refreshTokenTTL)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &Tokens{
		AccessToken:  token,
		RefreshToken: refreshtoken,
	}, nil
}

func (s *UsersService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.Get(ctx, id)
}
