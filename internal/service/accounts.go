package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/firebase"
	"shotwot_backend/pkg/logger"
	"time"

	fireauth "firebase.google.com/go/v4/auth"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type AccountsService struct {
	repo            repository.Accounts
	tokenManager    auth.TokenManager
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

func (s *AccountsService) SignUp(ctx context.Context, input AccountAuthInput) (*Tokens, error) {
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
		return nil, err
	}

	account := &domain.Account{
		Email: input.Email,
		Id:    usr.UID,
	}
	_, err = s.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}
	return s.createSession(ctx, usr.UID)

}

func (s *AccountsService) SignIn(ctx context.Context, input AccountAuthInput) (*Tokens, error) {
	if err := input.Validate(); err != nil {
		return nil, domain.ErrEmailPasswordInvalid
	}
	jsonValue, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	id, err := s.firebaseSignIn(jsonValue)
	if err != nil {
		return nil, err
	}
	return s.createSession(ctx, id)
}

func (s *AccountsService) createSession(ctx context.Context, userId string) (*Tokens, error) {
	logger.Debug(userId)
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

func (s *AccountsService) firebaseSignIn(body []byte) (string, error) {
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=AIzaSyD8h8LaQf_FaPQvbbn4eaU6hRLfjKGJEw0"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var id struct {
		// kind        string
		LocalId string
		// email       string
		// displayName string
		// idToken     string
		// registered  bool
	}
	err = json.Unmarshal(respBody, &id)
	if err != nil {
		return "", err
	}
	return id.LocalId, nil
}
