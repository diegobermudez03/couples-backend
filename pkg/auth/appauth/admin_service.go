package appauth

import (
	"context"
	"errors"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AdminAuthServiceImpl struct {
	repo auth.AuthRepository
	jwtSecret 	string
	accessTokenLife  int64
}

func NewAdminAuthService(repo auth.AuthRepository, jwtSecret string, accessTokenLife  int64) auth.AuthAdminService{
	return &AdminAuthServiceImpl{
		repo: repo,
		jwtSecret: jwtSecret,
		accessTokenLife: accessTokenLife,
	}	
}


func (s *AdminAuthServiceImpl) ValidateAccessToken(ctx context.Context, accessTokenString string) (*auth.AdminAccessClaims, error){
	accessToken, err := jwt.ParseWithClaims(accessTokenString, &auth.AdminAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil{
		if errors.Is(err, jwt.ErrTokenExpired){
			return nil, auth.ErrorExpiredAccessToken
		}
		return nil, auth.ErrorMalformedAccessToken
	}
	if !accessToken.Valid{
		return nil, auth.ErrorExpiredAccessToken
	}
	claims := accessToken.Claims.(*auth.AdminAccessClaims)
	return claims, nil
}


func (s *AdminAuthServiceImpl) CreateAccessToken(ctx context.Context, token string)(string, error){
	session, err := s.repo.GetAdminSessionByToken(ctx, token)
	if err != nil || session == nil{
		return "", auth.ErrorNonExistingSession 
	}
	accessToken, err := s.createAccessToken(session.Id)
	return accessToken, err
}

func (s *AdminAuthServiceImpl) createAccessToken(sessionId uuid.UUID)(string, error){
	claims := auth.AdminAccessClaims{
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.accessTokenLife*int64(time.Second)))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil{
		return "", auth.ErrorCreatingAccessToken
	}
	return tokenString, nil
}
