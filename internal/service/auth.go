package service

import (
	jwtauth "shotwot_backend/pkg/auth"
)

type AuthService struct {
	tokenManager jwtauth.TokenManager
}

func NewAuthService(tokenManager jwtauth.TokenManager) *AuthService {
	return &AuthService{
		tokenManager: tokenManager,
	}
}

func (a *AuthService) UserIdentity(token string) (*jwtauth.CustomClaims, error) {
	return a.tokenManager.Parse(token)

}
