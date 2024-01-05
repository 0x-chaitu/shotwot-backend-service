package service

import (
	"shotwot_backend/internal/repository"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/firebase"
	"time"
)

type AdminsService struct {
	repo               repository.Admins
	tokenManager       jwtauth.AdminTokenManager
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
	firebaseAuthClient *firebase.AuthClient
}
