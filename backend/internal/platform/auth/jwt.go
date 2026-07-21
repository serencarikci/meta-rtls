package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const TOKEN_ISSUER = "metartls"

type Claims struct {
	UserID   string `json:"uid"`
	TenantID string `json:"tid"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(secret string, ttlMinutes int) *TokenService {
	return &TokenService{
		secret: []byte(secret),
		ttl:    time.Duration(ttlMinutes) * time.Minute,
	}
}

func (s *TokenService) Issue(userID, tenantID, email, role string) (string, time.Time, error) {
	exp := time.Now().UTC().Add(s.ttl)
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        uuid.NewString(),
			Issuer:    TOKEN_ISSUER,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	return signed, exp, err
}

func (s *TokenService) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
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
