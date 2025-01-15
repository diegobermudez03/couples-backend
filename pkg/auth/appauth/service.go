package appauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth/domainauth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	errorInsecurePassword = errors.New("the password isnt secure enough")
	errorHashingPassword = errors.New("the password couldnt be hashed")
	errorJWTToken = errors.New("an error ocurred creating the JWT token")
)

type AuthServiceImpl struct {
	authRepo 			domainauth.AuthRepository
	accessTokenLife 	int64
	jwtSecret 			string
}

func NewAuthService(authRepo domainauth.AuthRepository, accessTokenLife int64) domainauth.AuthService{
	return &AuthServiceImpl{
		authRepo: authRepo,
		accessTokenLife: accessTokenLife,
	}
}


func(s *AuthServiceImpl) RegisterUser(email string, password string, device string, os string) (string, error){
	// data verifications
	if num := len(password); num < 6 {
		return "", errorInsecurePassword
	}
	if match, err := regexp.MatchString(`\d`, password); !match || err != nil {
		return "", errorInsecurePassword
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", errorHashingPassword
	}

	//create user auth
	userId := uuid.New()
	err = s.authRepo.CreateUserAuth(
		userId,
		email,
		string(hashBytes),
	)
	if err != nil{
		return "",  err
	}

	//create the session
	return s.createSession(userId, device, os)
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////
//								private functions
////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////


func (s *AuthServiceImpl) createSession(userId uuid.UUID, device string, os string) (string, error){
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes) 
	if err != nil{
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(randomBytes)

	err = s.authRepo.CreateSession(
		uuid.New(),
		userId,
		token, 
		device, 
		os,
		time.Now().Add(time.Hour* 8760),
	)
	if err != nil{
		return "", err 
	}
	return token, nil
}


func (s *AuthServiceImpl) createAccessToken(userId uuid.UUID, coupleId uuid.UUID) (string, error){
	claims := domainauth.AccessClaims{
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