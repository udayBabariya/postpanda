package auth

import (
	"errors"
	"postpanda/backend-go/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID         string `json:"userId"`
	OrganizationID string `json:"organizationId"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, organizationID string) (string, error) {
	claims := Claims{
		UserID:         userID,
		OrganizationID: organizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.App.JWTSecret))
}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.App.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
