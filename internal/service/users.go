package service

import (
	"context"
	"fmt"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/firebase"
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

func (s *UsersService) SignUp(ctx context.Context, input AccountAuthInput) (*Tokens, error) {
	token, err := s.firebaseAuthClient.VerifyIDToken(ctx, input.IdToken)
	if err != nil {
		if fireauth.IsEmailAlreadyExists(err) {
			return nil, domain.ErrAccountAlreadyExists
		}
		logger.Errorf("firebase user creation error %v", err)
		return nil, err
	}
	account := &domain.User{
		Email:   token.Claims["email"].(string),
		Id:      token.UID,
		Created: time.Now(),
	}

	err = s.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}
	return s.createSession(ctx, token.UID)

}

func (s *UsersService) SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error) {
	token, err := s.firebaseAuthClient.VerifyIDToken(ctx, input.IdToken)
	if err != nil {
		logger.Errorf("firebase user signin error %v", err)
		return nil, err
	}
	return s.createSession(ctx, token.UID)
}

func (s *UsersService) Update(ctx context.Context, input *domain.User) (*domain.User, error) {
	user, err := s.repo.Update(ctx, input)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info(user)
	return user, nil
}

func (s *UsersService) Delete(ctx context.Context, id string) error {
	return s.firebaseAuthClient.DeleteUser(ctx, id)
}

func (s *UsersService) createSession(ctx context.Context, userId string) (*Tokens, error) {
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
