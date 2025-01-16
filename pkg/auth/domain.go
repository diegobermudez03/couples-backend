package auth

import (
	"context"
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
	RegisterUser(ctx context.Context, email string, password string, device string, os string) (refreshToken string, err error)
	LoginUser(ctx context.Context, email string, password string, device string, os string) (refreshToken string, err error)
	GetUserIdFromSession(ctx context.Context, token string) (*uuid.UUID, error)
	CreateAnonymousSession(ctx context.Context, userId uuid.UUID) (refreshToken string, err error)
	VinculateAuthWithUser(ctx context.Context, token string, userId uuid.UUID) error
}

type AuthRepository interface {
	CreateUserAuth(ctx context.Context, id uuid.UUID, email string, hash string) error
	CreateEmptyUser(ctx context.Context, id uuid.UUID, userId uuid.UUID) error
	CreateSession(ctx context.Context, id uuid.UUID, userId uuid.UUID, token string, device *string, os *string, expiresAt time.Time) error
	GetUserByEmail(ctx context.Context, email string) (*UserAuthModel, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*UserAuthModel, error)
	GetSessionByToken(ctx context.Context, token string) (*SessionModel, error)
	UpdateAuthUserId(ctx context.Context, authId uuid.UUID, userId uuid.UUID) error
}