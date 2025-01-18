package appauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/diegobermudez03/couples-backend/pkg/users"
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
	errorUserAlreadyHasUser = errors.New("there's already an user created")
)


type AuthServiceImpl struct {
	authRepo 			auth.AuthRepository
	usersService 		users.UsersService
	accessTokenLife 	int64
	refreshTokenLife 	int64
	jwtSecret 			string
}

func NewAuthService(authRepo auth.AuthRepository, usersService users.UsersService, accessTokenLife int64, refreshTokenLife int64) auth.AuthService{
	return &AuthServiceImpl{
		authRepo: authRepo,
		usersService: usersService,
		accessTokenLife: accessTokenLife,
		refreshTokenLife : refreshTokenLife,
	}
}


func(s *AuthServiceImpl) RegisterUserAuth(ctx context.Context, email, password, device, os, token string) (string, error){
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
	hashString := string(hashBytes)

	var userId uuid.UUID
	log.Print(token)
	// check token 
	if token != ""{
		session, _ := s.authRepo.GetSessionByToken(ctx, token)
		log.Print(session)
		if session != nil{
			log.Print("session not null")
			userAuth, _ := s.authRepo.GetUserById(ctx, session.UserAuthId)
			log.Print(userAuth)
			if userAuth != nil && s.checkIfAnonymousAuth(userAuth){
				log.Print("user anonymous")
				userId = userAuth.Id
				userAuth.Email = &email
				userAuth.Hash = &hashString
				if err := s.authRepo.UpdateAuthUserById(
					ctx,
					userId,
					userAuth,
				); err != nil{
					return "", err 
				}
				return token, nil
			}
		}
	}
	//create user auth
	userId = uuid.New()
	err = s.authRepo.CreateUserAuth(
		ctx,
		userId,
		email,
		hashString,
	)
	if err != nil{
		return "",  err
	}

	//create the session
	return s.createSession(ctx, userId, &device, &os)
}


func (s *AuthServiceImpl) LoginUserAuth(ctx context.Context, email string, password string, device string, os string) (string, error){
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil{
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Hash), []byte(password)); err != nil{
		return "", errorIncorrectPassword
	}

	return s.createSession(ctx, user.Id, &device, &os)
}

func (s *AuthServiceImpl) CloseUsersSession(ctx context.Context, token string) error{
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return err
	}
	if err := s.authRepo.DeleteSessionById(ctx, session.Id); err != nil{
		return err
	}
	authUser, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil{
		return err
	}
	//if it's an anonymous account then delete both the account and the user
	if s.checkIfAnonymousAuth(authUser){
		if err := s.authRepo.DeleteUserAuthById(ctx, authUser.Id); err != nil{
			return err 
		}
		if authUser.UserId != nil{
			if err := s.usersService.DeleteUserById(ctx, *authUser.UserId); err != nil{
				return err
			}
		}
	}
	return nil
}

func (s *AuthServiceImpl) CreateTempCouple(ctx context.Context, token string, startDate int) (int, error){
	userId, err := s.getUserIdFromSession(ctx, token)
	if err != nil{
		return 0, err  
	}
	if userId == nil{
		return 0, auth.ErrorNoActiveUser
	}
	return s.usersService.CreateTempCouple(ctx, *userId, startDate)
}

func (s *AuthServiceImpl) CreateUser(ctx context.Context, token, firstName, lastName, gender, countryCode, languageCode string,birthDate int,) (string, error){
	//check token if its validate
	session, _ := s.authRepo.GetSessionByToken(ctx, token)
	if session != nil{
		userAuth, _ := s.authRepo.GetUserById(ctx, session.UserAuthId)
		if userAuth != nil && userAuth.UserId != nil{
			return "", errorUserAlreadyHasUser
		}
	}

	// create user with users service (receives the userId)
	userId,  err := s.usersService.CreateUser(ctx, firstName, lastName, gender, countryCode, languageCode, birthDate)
	if err != nil{
		return "", err 
	}
	//	if no token, create anonymous auth user and connect with user
	if session != nil{
		if err := s.authRepo.UpdateAuthUserId(ctx, session.UserAuthId, *userId); err != nil{
			return "", err 
		}
	} else{
		authId := uuid.New()
		if err := s.authRepo.CreateEmptyUser(ctx, authId, *userId); err != nil{
			return "", err 
		}
		return s.createSession(ctx, authId, nil, nil)
	}
	// if token, get the user auth and connect with userId
	return token, nil
}


func (s *AuthServiceImpl) CheckUserAuthStatus(ctx context.Context, token string) (string, error){
	userId, err := s.getUserIdFromSession(ctx, token)
	if err != nil{
		return "", err
	}
	if userId == nil{
		return auth.StatusNoUserCreated, nil 
	}
	couple, _ := s.usersService.GetCoupleFromUser(ctx, *userId)
	if couple == nil{
		return auth.StatusUserCreated, nil 
	}
	return auth.StatusCoupleCreated, nil

}

func (s *AuthServiceImpl) ConnectCouple(ctx context.Context, token string, code int) (string, error) {
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return "", err 
	}
	auth, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil{
		return "", err 
	}
	coupleId, err := s.usersService.ConnectCouple(ctx, *auth.UserId, code)
	if err != nil{
		return "", err
	}
	return s.createAccessToken(*auth.UserId, *coupleId, session.Id)
}


func (s *AuthServiceImpl) CreateAccessToken(ctx context.Context, token string)(string, error){
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return "", err 
	}
	if session.ExpiresAt.Before(time.Now()){
		s.authRepo.DeleteSessionById(ctx, session.Id)
		return "", auth.ErrorExpiredRefreshToken
	}
	user, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil{
		return "", err 
	}
	if user.UserId == nil{
		return "", auth.ErrorNoActiveUser
	}
	couple, err := s.usersService.GetCoupleFromUser(ctx, *user.UserId)
	if err != nil{
		return "", err 
	}
	if couple == nil {
		return "", auth.ErrorNoActiveCoupleFromUser
	}
	return s.createAccessToken(*user.UserId, couple.Id, session.Id)
}

func (s *AuthServiceImpl) ValidateAccessToken(ctx context.Context, accessTokenString string) (*auth.AccessClaims, error){
	accessToken, err := jwt.ParseWithClaims(accessTokenString, &auth.AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
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
	claims := accessToken.Claims.(*auth.AccessClaims)
	return claims, nil
}

func (s *AuthServiceImpl) LogoutSession(ctx context.Context, sessionId uuid.UUID) error{
	session, err := s.authRepo.GetSessionById(ctx, sessionId)
	if err != nil{
		return err 
	}
	authUser, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil{
		return err 
	}
	if s.checkIfAnonymousAuth(authUser){
		return auth.ErrorCantLogoutAnonymousAcc
	}
	if err := s.authRepo.DeleteSessionById(ctx, session.Id); err != nil{
		return err 
	}
	return nil
}


/////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////
//								private functions
////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////


func (s *AuthServiceImpl) createSession(ctx context.Context, authId uuid.UUID, device *string, os *string) (string, error){
	var token string 
	for{
		randomBytes := make([]byte, 32)
		_, err := rand.Read(randomBytes) 
		if err != nil{
			return "", err
		}
		token = base64.URLEncoding.EncodeToString(randomBytes)
		session, _ := s.authRepo.GetSessionByToken(ctx, token)
		if session == nil{
			break
		}
	}
	
	err := s.authRepo.CreateSession(
		ctx, 
		uuid.New(),
		authId,
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


func (s *AuthServiceImpl) createAccessToken(userId uuid.UUID, coupleId uuid.UUID, sessionId uuid.UUID) (string, error){
	claims := auth.AccessClaims{
		UserId: userId,
		CoupleId: coupleId,
		SessionId : sessionId,
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

func (s *AuthServiceImpl) getUserIdFromSession(ctx context.Context, token string) (*uuid.UUID, error){
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

func (s *AuthServiceImpl) checkIfAnonymousAuth(auth *auth.UserAuthModel) bool{
	return auth.Email == nil && auth.OauthProvider == nil
}