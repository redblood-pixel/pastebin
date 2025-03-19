package tokenutil

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	SigningKey      string        `yaml:"signing_key"`
}

type TokenManager struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	signingKey      []byte
}

type JWTCustomClaims struct {
	userID int
	jwt.RegisteredClaims
}

type RefreshToken struct {
	ID        string
	UserID    int
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func New(cfg *Config) *TokenManager {
	return &TokenManager{
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
		signingKey:      []byte(cfg.SigningKey),
	}
}

func (tm *TokenManager) CreateAccessToken(userID int) (string, error) {
	claims := JWTCustomClaims{
		userID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.accessTokenTTL)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	token, err := t.SignedString(tm.signingKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (tm *TokenManager) ParseAccessToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method")
		}
		return tm.signingKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims.userID, nil
	}
	return 0, fmt.Errorf("access token not valid")
}

func (tm *TokenManager) GetRefreshTTL() time.Duration {
	return tm.refreshTokenTTL
}
