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

type AdminAuthService struct {
	adminTokenManager jwtauth.AdminTokenManager
}

func NewAdminAuthService(adminTokenManager jwtauth.AdminTokenManager) *AdminAuthService {
	return &AdminAuthService{
		adminTokenManager: adminTokenManager,
	}
}

func (a *AdminAuthService) AdminIdentity(token string) (*jwtauth.CustomAdminClaims, error) {
	return a.adminTokenManager.Parse(token)
}
