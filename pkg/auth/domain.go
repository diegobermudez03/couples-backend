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
	RegisterUserAuth(ctx context.Context, email string, password string, device string, os string) (refreshToken string, err error)
	LoginUserAuth(ctx context.Context, email string, password string, device string, os string) (refreshToken string, err error)
	CloseSession(ctx context.Context, token string) (error)
	CreateTempCouple(ctx context.Context, token string, startDate int) (int, error)
	CreateUser(ctx context.Context, token, firstName, lastName, gender, countryCode, languageCode string,birthDate int,) (refrshToken string, err error)
	ConnectCouple(ctx context.Context, token string, code int) (accessToken string, err error)
	CheckUserAuthStatus(ctx context.Context, token string) (string, error)
	CreateAccessToken(ctx context.Context, token string)(string, error)
	ValidateAccessToken(ctx context.Context, accessTokenString string) (*AccessClaims, error)
}

type AuthRepository interface {
	CreateUserAuth(ctx context.Context, id uuid.UUID, email string, hash string) error
	CreateEmptyUser(ctx context.Context, id uuid.UUID, userId uuid.UUID) error
	CreateSession(ctx context.Context, id uuid.UUID, userId uuid.UUID, token string, device *string, os *string, expiresAt time.Time) error
	GetUserByEmail(ctx context.Context, email string) (*UserAuthModel, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*UserAuthModel, error)
	GetSessionByToken(ctx context.Context, token string) (*SessionModel, error)
	UpdateAuthUserId(ctx context.Context, authId uuid.UUID, userId uuid.UUID) error
	DeleteSessionById(ctx context.Context, sessionId uuid.UUID) error
	DeleteUserAuthById(ctx context.Context, authId uuid.UUID) error
}


///// messages
const (
	StatusNoUserCreated = "there's no user associated"
	StatusUserCreated = "user has an user associated"
	StatusCoupleCreated = "user has a couple associated"
)