package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
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
	repo            repository.Users
	tokenManager    jwtauth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	authClient      *firebase.AuthClient
}

func (a AccountAuthInput) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required, validation.Length(6, 64)),
	)
}

func NewUsersService(repo repository.Users, tokenManager jwtauth.TokenManager, accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration, authClient *firebase.AuthClient) *UsersService {
	return &UsersService{
		repo:            repo,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		authClient:      authClient,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input AccountAuthInput) (*Tokens, error) {
	if err := input.Validate(); err != nil {
		return nil, domain.ErrEmailPasswordInvalid
	}
	params := (&fireauth.UserToCreate{}).
		Email(input.Email).
		EmailVerified(false).
		Password(input.Password)
	usr, err := s.authClient.CreateUser(ctx, params)
	if err != nil {
		if fireauth.IsEmailAlreadyExists(err) {
			return nil, domain.ErrAccountAlreadyExists
		}
		logger.Errorf("firebase user creation error %v", err)
		return nil, err
	}
	u, err := s.authClient.GetUser(ctx, usr.UID)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully fetched user data: %v\n", u)
	// email, err := s.authClient.EmailVerificationLink(ctx, input.Email)
	if err != nil {
		logger.Errorf("error generating email link: %v\n", err)
		return nil, err
	}
	account := &domain.User{
		Email: input.Email,
		Id:    usr.UID,
	}
	_, err = s.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}
	return s.createSession(ctx, usr.UID)

}

func (s *UsersService) SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error) {
	if err := input.Validate(); err != nil {
		return nil, domain.ErrEmailPasswordInvalid
	}
	jsonValue, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	id, err := s.firebaseSignIn(jsonValue)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return s.createSession(ctx, id)
}

func (s *UsersService) Update(ctx context.Context, input *domain.User) (*domain.User, error) {
	if err := input.Validate(); err != nil {
		return nil, domain.ErrInvalidInput
	}
	user, err := s.repo.Update(ctx, input)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info(user)
	return user, nil
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

func (s *UsersService) firebaseSignIn(body []byte) (string, error) {
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=AIzaSyD8h8LaQf_FaPQvbbn4eaU6hRLfjKGJEw0"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error(err)
		return "", err
	} else if resp.StatusCode != 200 {
		return "", errors.New("firebase error")
	}

	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var id struct {
		LocalId string
	}
	err = json.Unmarshal(respBody, &id)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return id.LocalId, nil
}
