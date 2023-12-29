package service

import (
	"context"
	"errors"
	"fmt"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/firebase"
	"shotwot_backend/pkg/logger"
	"time"

	fireauth "firebase.google.com/go/auth"
)

type AccountsService struct {
	repo            repository.Accounts
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	authClient      *firebase.AuthClient
}

func NewAccountsService(repo repository.Accounts, tokenManager auth.TokenManager, accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration, authClient *firebase.AuthClient) *AccountsService {
	return &AccountsService{
		repo:            repo,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		authClient:      authClient,
	}
}

func (s *AccountsService) SignUp(ctx context.Context, input AccountSignUpInput) (*Tokens, error) {
	params := (&fireauth.UserToCreate{}).
		Email(input.Email).
		EmailVerified(false).
		Password(input.Password)
	usr, err := s.authClient.CreateUser(ctx, params)
	if err != nil {
		logger.Info(err)
	}
	account := &domain.Account{
		Email: input.Email,
		Id:    usr.UID,
	}
	_, err = s.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}
	if err != nil {
		if errors.Is(err, domain.ErrAccountAlreadyExists) {
			return nil, err
		}
		return nil, err
	}
	return s.createSession(ctx, usr.UID)

}

func (s *AccountsService) SignIn(ctx context.Context, input AccountSignInInput) (*Tokens, error) {

	return s.createSession(ctx, "")
}

func (s *AccountsService) createSession(ctx context.Context, userId string) (*Tokens, error) {
	token, err := s.tokenManager.NewJWT(fmt.Sprint(userId), s.accessTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshtoken, err := s.tokenManager.NewJWT(fmt.Sprint(userId), s.refreshTokenTTL)
	if err != nil {
		return nil, err
	}
	return &Tokens{
		AccessToken:  token,
		RefreshToken: refreshtoken,
	}, nil
}
