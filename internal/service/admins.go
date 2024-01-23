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

	"firebase.google.com/go/v4/auth"
	validation "github.com/go-ozzo/ozzo-validation"
)

type AdminsService struct {
	repo               repository.Admins
	tokenManager       jwtauth.AdminTokenManager
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
	firebaseAuthClient *firebase.AuthClient
}

func NewAdminsService(repo repository.Admins, tokenManager jwtauth.AdminTokenManager, accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration, firebaseAuthClient *firebase.AuthClient) *AdminsService {
	return &AdminsService{
		repo:               repo,
		tokenManager:       tokenManager,
		accessTokenTTL:     accessTokenTTL,
		refreshTokenTTL:    refreshTokenTTL,
		firebaseAuthClient: firebaseAuthClient,
	}
}

func (a *AdminsService) SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error) {
	tokenId, err := a.firebaseAuthClient.VerifyIDToken(ctx, input.IdToken)
	if err != nil {
		logger.Errorf("firebase admin user signin error %v", err)
		return nil, err
	}
	admin, err := a.repo.Get(ctx, tokenId.UID)
	if err != nil {
		return nil, err
	}
	return a.createAdminSession(ctx, tokenId.UID, admin.Role)
}

func (a *AdminsService) GetAllAdmins(ctx context.Context) ([]*domain.Admin, error) {
	return a.repo.GetAdmins(ctx)
}

func (a *AdminsService) CreateAdmin(ctx context.Context, input AccountAuthInput) error {
	err := validation.ValidateStruct(&input,
		validation.Field(&input.Password, validation.Required),
		validation.Field(&input.Email, validation.Required),
		validation.Field(&input.Role, validation.Required, validation.Max(4), validation.Min(1)),
	)
	if err != nil {
		logger.Error(err)
		return err
	}
	params := (&auth.UserToCreate{}).
		Email(input.Email).
		Password(input.Password)
	admin, err := a.firebaseAuthClient.Client.CreateUser(ctx, params)
	if err != nil {
		logger.Errorf("error creating user: %v\n", err)
		return err
	}

	err = a.repo.Create(ctx, &domain.Admin{
		Id:      admin.UID,
		Email:   admin.Email,
		Role:    input.Role,
		Created: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *AdminsService) Update(ctx context.Context, admin *domain.Admin) (*domain.Admin, error) {
	user, err := a.repo.Update(ctx, admin)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info(user)
	return user, nil

}

func (a *AdminsService) Delete(ctx context.Context, id string) error {
	return a.firebaseAuthClient.DeleteUser(ctx, id)
}

func (s *AdminsService) createAdminSession(ctx context.Context, adminId string, role int) (*Tokens, error) {
	token, err := s.tokenManager.NewJWT(fmt.Sprint(adminId), s.accessTokenTTL, role)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	refreshtoken, err := s.tokenManager.NewJWT(fmt.Sprint(adminId), s.refreshTokenTTL, role)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &Tokens{
		AccessToken:  token,
		RefreshToken: refreshtoken,
	}, nil
}
