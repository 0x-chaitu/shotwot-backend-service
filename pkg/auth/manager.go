package jwtauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenManager provides logic for JWT & Refresh tokens generation and parsing.
type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	Parse(accessToken string) (*CustomClaims, error)
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
		return nil, err
	}
	if err := json.Unmarshal(jsonbody, &customClaims); err != nil {
		return nil, err
	}
	return &customClaims, nil
}
