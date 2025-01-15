package domainauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct{
	UserId 		uuid.UUID
	CoupleId 	uuid.UUID
	jwt.RegisteredClaims
}

type AuthService interface {
	RegisterUser(email string, password string, device string, os string) (refreshToken string, err error)
}

type AuthRepository interface {
	CreateUserAuth(id uuid.UUID, email string, hash string) error
	CreateSession(id uuid.UUID, userId uuid.UUID, token string, device string, os string, expiresAt time.Time) error
}