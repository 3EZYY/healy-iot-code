package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rafif/healy-backend/internal/domain"
	"github.com/rafif/healy-backend/pkg/config"
)

type TokenGenerator interface {
	GenerateToken(user domain.User) (string, int64, error)
	ValidateToken(tokenStr string) (*jwt.MapClaims, error)
}

type jwtGenerator struct {
	cfg *config.Config
}

func NewJWTGenerator(cfg *config.Config) TokenGenerator {
	return &jwtGenerator{
		cfg: cfg,
	}
}

func (g *jwtGenerator) GenerateToken(user domain.User) (string, int64, error) {
	expirationTime := time.Now().Add(time.Duration(g.cfg.JWTExpiryHours) * time.Hour)
	
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.cfg.JWTSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

func (g *jwtGenerator) ValidateToken(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(g.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}
