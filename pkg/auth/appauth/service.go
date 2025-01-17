package appauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	errorInsecurePassword = errors.New("the password isnt secure enough")
	errorHashingPassword = errors.New("the password couldnt be hashed")
	errorJWTToken = errors.New("an error ocurred creating the JWT token")
	errorEmailAlreadyUsed = errors.New("the email already has an account associated")
	errorIncorrectPassword = errors.New("incorrect password")
)

type AuthServiceImpl struct {
	authRepo 			auth.AuthRepository
	accessTokenLife 	int64
	refreshTokenLife 	int64
	jwtSecret 			string
}

func NewAuthService(authRepo auth.AuthRepository, accessTokenLife int64, refreshTokenLife int64) auth.AuthService{
	return &AuthServiceImpl{
		authRepo: authRepo,
		accessTokenLife: accessTokenLife,
		refreshTokenLife : refreshTokenLife,
	}
}


func(s *AuthServiceImpl) RegisterUser(ctx context.Context, email string, password string, device string, os string) (string, error){
	// data verifications
	if num := len(password); num < 6 {
		return "", errorInsecurePassword
	}
	if match, err := regexp.MatchString(`\d`, password); !match || err != nil {
		return "", errorInsecurePassword
	}

	// confirm email uniqueness
	if _, err := s.authRepo.GetUserByEmail(ctx, email); err == nil || !errors.Is(err, auth.ErrorNoUserFoundEmail){
		return "", errorEmailAlreadyUsed
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", errorHashingPassword
	}

	//create user auth
	userId := uuid.New()
	err = s.authRepo.CreateUserAuth(
		ctx,
		userId,
		email,
		string(hashBytes),
	)
	if err != nil{
		return "",  err
	}

	//create the session
	return s.createSession(ctx, userId, &device, &os)
}


func (s *AuthServiceImpl) LoginUser(ctx context.Context, email string, password string, device string, os string) (string, error){
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil{
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Hash), []byte(password)); err != nil{
		return "", errorIncorrectPassword
	}

	return s.createSession(ctx, user.Id, &device, &os)
}

func (s *AuthServiceImpl) GetUserIdFromSession(ctx context.Context, token string) (*uuid.UUID, error){
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return nil, err 
	}
	user, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil{
		return nil, err 
	}
	return user.UserId, nil
}

func (s *AuthServiceImpl) CreateAnonymousSession(ctx context.Context, userId uuid.UUID) (string, error){
	AuthId := uuid.New()
	if err := s.authRepo.CreateEmptyUser(ctx, AuthId, userId); err != nil{
		return "", err 
	}
	return s.createSession(ctx, AuthId, nil, nil)
}

func (s *AuthServiceImpl) VinculateAuthWithUser(ctx context.Context, token string, userId uuid.UUID) error{
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return err 
	}
	if err := s.authRepo.UpdateAuthUserId(ctx, session.UserAuthId, userId); err != nil{
		return err 
	}
	return nil
}


func (s *AuthServiceImpl) CloseSession(ctx context.Context, token string) (userId *uuid.UUID, err error){
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return nil, err
	}
	if err := s.authRepo.DeleteSessionById(ctx, session.Id); err != nil{
		return nil, err
	}
	authUser, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil{
		return nil, err
	}
	//if it's an anonymous account then delete both the account and the user
	if authUser.Email == nil && authUser.OauthProvider == nil{
		if err := s.authRepo.DeleteUserAuthById(ctx, authUser.Id); err != nil{
			return nil, err 
		}
		return authUser.UserId, nil
	}
	return nil, nil
}


/////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////
//								private functions
////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////


func (s *AuthServiceImpl) createSession(ctx context.Context, userId uuid.UUID, device *string, os *string) (string, error){
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes) 
	if err != nil{
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(randomBytes)

	err = s.authRepo.CreateSession(
		ctx, 
		uuid.New(),
		userId,
		token, 
		device, 
		os,
		time.Now().Add(time.Duration(s.refreshTokenLife) * time.Hour),
	)
	if err != nil{
		return "", err 
	}
	return token, nil
}


func (s *AuthServiceImpl) createAccessToken(userId uuid.UUID, coupleId uuid.UUID) (string, error){
	claims := auth.AccessClaims{
		UserId: userId,
		CoupleId: coupleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.accessTokenLife*int64(time.Second)))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil{
		return "", errorJWTToken
	}
	return tokenString, nil
}