package jwtauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"shotwot_backend/pkg/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	SuperAdmin = iota + 1
	Admin
	BriefManager
	Curator
)

// TokenManager provides logic for JWT & Refresh tokens generation and parsing.
type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	Parse(accessToken string) (*CustomClaims, error)
}

type AdminTokenManager interface {
	NewJWT(userId string, ttl time.Duration, role int) (string, error)
	Parse(accessToken string) (*CustomAdminClaims, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

type CustomClaims struct {
	IsPro   bool
	Subject string
}

func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(ttl).Unix(),
		"data": CustomClaims{
			IsPro:   false,
			Subject: userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (*CustomClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from token")
	}
	data := claims["data"].(map[string]interface{})
	customClaims := CustomClaims{}
	jsonbody, err := json.Marshal(data)
	if err != nil {
		logger.Info(err)
		return nil, err
	}
	if err := json.Unmarshal(jsonbody, &customClaims); err != nil {
		logger.Info(err)
		return nil, err
	}
	return &customClaims, nil
}

type AdminManager struct {
	signingKey string
}

func NewAdminManager(signingKey string) (*AdminManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &AdminManager{signingKey: signingKey}, nil
}

type CustomAdminClaims struct {
	AdminRole int
	Subject   string
}

func (m *AdminManager) NewJWT(userId string, ttl time.Duration, role int) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(ttl).Unix(),
		"data": CustomAdminClaims{
			AdminRole: role,
			Subject:   userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *AdminManager) Parse(accessToken string) (*CustomAdminClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from token")
	}
	data := claims["data"].(map[string]interface{})
	customClaims := CustomAdminClaims{}
	jsonbody, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if err := json.Unmarshal(jsonbody, &customClaims); err != nil {
		logger.Error(err)
		return nil, err
	}
	return &customClaims, nil
}
