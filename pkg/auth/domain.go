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
	SessionId 	uuid.UUID
	jwt.RegisteredClaims
}

type AuthService interface {
	RegisterUserAuth(ctx context.Context, email, password, device, os, token string) (refreshToken string, err error)
	LoginUserAuth(ctx context.Context, email string, password string, device string, os string) (refreshToken string, err error)
	CloseUsersSession(ctx context.Context, token string) (error)
	CreateTempCouple(ctx context.Context, token string, startDate int) (int, error)
	CreateUser(ctx context.Context, token, firstName, lastName, gender, countryCode, languageCode string,birthDate int,) (refrshToken string, err error)
	ConnectCouple(ctx context.Context, token string, code int) (accessToken string, err error)
	CheckUserAuthStatus(ctx context.Context, token string) (string, error)
	CreateAccessToken(ctx context.Context, token string)(string, *string, error)
	ValidateAccessToken(ctx context.Context, accessTokenString string) (*AccessClaims, error)
	GetTempCoupleOfUser(ctx context.Context, token string)(*TempCoupleModel, error)
	LogoutSession(ctx context.Context, sessionId uuid.UUID) error
	SuscribeTempCoupleNot(ctx context.Context, token string)(chan string, *uuid.UUID, error)
	RemoveCodeSuscriber(userId uuid.UUID)
}

type AuthRepository interface {
	CreateUserAuth(ctx context.Context, id uuid.UUID, email string, hash string) (int, error)
	CreateEmptyUser(ctx context.Context, id uuid.UUID, userId uuid.UUID) (int, error)
	CreateSession(ctx context.Context, sessionModel *SessionModel) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*UserAuthModel, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*UserAuthModel, error)
	GetSessionByToken(ctx context.Context, token string) (*SessionModel, error)
	GetSessionById(ctx context.Context, id uuid.UUID) (*SessionModel, error)
	UpdateAuthUserId(ctx context.Context, authId uuid.UUID, userId uuid.UUID) (int, error)
	DeleteSessionById(ctx context.Context, sessionId uuid.UUID) (int, error)
	DeleteUserAuthById(ctx context.Context, authId uuid.UUID) (int, error)
	UpdateAuthUserById(ctx context.Context, authId uuid.UUID, authModel *UserAuthModel) (int, error)
	UpdateSessionLastUsed(ctx context.Context, sessionId uuid.UUID, lastTime time.Time) (int, error)
}


///// messages
const (
	StatusNoUserCreated = "NO_USER_ASSOCIATED"
	StatusUserCreated = "USER_ASSOCIATED"
	StatusCoupleCreated = "COUPLE_CREATED"
	StatusPartnerWithoutNickname = "PARTNER_WITHOUT_NICKNAME"
	StatusVinculated = "PARTNER_VINCULATED"
)