package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT returns a signed token with subject=userID and claim role
func GenerateJWT(userID string, role string, secret string, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims.Subject,
		"exp":   claims.ExpiresAt.Unix(),
		"iat":   claims.IssuedAt.Unix(),
		"role":  role,
	})
	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenStr string, secret string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
