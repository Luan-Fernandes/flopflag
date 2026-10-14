package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token inválido")
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret  []byte
	expires time.Duration
}

func NewTokenManager(secret string, expires time.Duration) *TokenManager {
	return &TokenManager{
		secret:  []byte(secret),
		expires: expires,
	}
}

func (m *TokenManager) Sign(userID uuid.UUID, email, role string) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expires)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}

func (m *TokenManager) Parse(tokenString string) (uuid.UUID, string, error) {
	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		return uuid.Nil, "", ErrInvalidToken
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return uuid.Nil, "", ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, "", ErrInvalidToken
	}

	return userID, claims.Role, nil
}
