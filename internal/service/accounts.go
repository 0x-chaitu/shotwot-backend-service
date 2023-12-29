package service

import (
	"context"
	"errors"
	"fmt"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/hash"
	"time"
)

type AccountsService struct {
	repo            repository.Accounts
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAccountsService(repo repository.Accounts, tokenManager auth.TokenManager, accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration) *AccountsService {
	return &AccountsService{
		repo:            repo,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AccountsService) SignUp(ctx context.Context, input AccountSignUpInput) (*Tokens, error) {

	passwordHash, err := hash.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	account := &domain.Account{
		Email:    input.Email,
		Name:     input.Name,
		Password: passwordHash,
	}

	user, err := s.repo.Create(ctx, account)
	if err != nil {
		if errors.Is(err, domain.ErrAccountAlreadyExists) {
			return nil, err
		}
		return nil, err
	}
	return s.createSession(ctx, user.Id)

}

func (s *AccountsService) SignIn(ctx context.Context, input AccountSignInInput) (*Tokens, error) {
	user, err := s.repo.GetByCredentials(ctx, input.Name, input.Password)
	if err != nil {
		return nil, err
	}
	err = hash.Match(input.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("details not valid")
	}
	return s.createSession(ctx, user.Id)
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
