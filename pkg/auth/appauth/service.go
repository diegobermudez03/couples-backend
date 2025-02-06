package appauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/diegobermudez03/couples-backend/pkg/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


type AuthServiceImpl struct {
	authRepo 			auth.AuthRepository
	usersService 		users.UsersService
	accessTokenLife 	int64
	refreshTokenLife 	int64
	jwtSecret 			string
	codeSuscribers 		map[uuid.UUID] chan string
	suscribersMutex 	sync.RWMutex
}


func NewAuthService(authRepo auth.AuthRepository, usersService users.UsersService, accessTokenLife int64, refreshTokenLife int64, jwtSecret string) auth.AuthService{
	return &AuthServiceImpl{
		authRepo: authRepo,
		usersService: usersService,
		accessTokenLife: accessTokenLife,
		refreshTokenLife : refreshTokenLife,
		jwtSecret: jwtSecret,
		codeSuscribers: make(map[uuid.UUID] chan string),
		suscribersMutex : sync.RWMutex{},
	}
}


func(s *AuthServiceImpl) RegisterUserAuth(ctx context.Context, email, password, device, os, token string) (string, error){
	// data verifications
	if num := len(password); num < 6 {
		return "", auth.ErrorInsecurePassword
	}
	if match, err := regexp.MatchString(`\d`, password); !match || err != nil {
		return "", auth.ErrorInsecurePassword
	}
	email = strings.ToLower(email)
	// confirm email uniqueness
	if acc, err := s.authRepo.GetUserByEmail(ctx, email); err != nil{
		return "", auth.ErrorCreatingAccount
	}else if acc != nil{
		return "", auth.ErrorEmailAlreadyUsed
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		log.Print("error hashing password: ", err.Error())
		return "", auth.ErrorCreatingAccount
	}
	hashString := string(hashBytes)

	var userId uuid.UUID
	
	// check token 
	if token != ""{
		session, _ := s.authRepo.GetSessionByToken(ctx, token)
		//if the token do have a session associated
		if session != nil{
			userAuth, _ := s.authRepo.GetUserById(ctx, session.UserAuthId)
			// if the session has an anonymous account associated
			if userAuth != nil && s.checkIfAnonymousAuth(userAuth){
				userId = userAuth.Id
				userAuth.Email = &email
				userAuth.Hash = &hashString
				if num, err := s.authRepo.UpdateAuthUserById(
					ctx,
					userId,
					userAuth,
				); err != nil || num == 0{
					return "", auth.ErrorVinculatingAccount 
				}
				return token, nil
			}
		}
	}
	//create user auth
	userId = uuid.New()
	num, err := s.authRepo.CreateUserAuth(
		ctx,
		userId,
		email,
		hashString,
	)
	if err != nil || num == 0{
		return "",  auth.ErrorCreatingAccount
	}

	//create the session
	return s.createSession(ctx, userId, &device, &os)
}


func (s *AuthServiceImpl) LoginUserAuth(ctx context.Context, email string, password string, device string, os string) (string, error){
	email = strings.ToLower(email)
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil{
		return "", auth.ErrorWithLogin
	} else if user == nil {
		return "", auth.ErrorNoUserFoundEmail
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Hash), []byte(password)); err != nil{
		return "", auth.ErrorIncorrectPassword
	}

	return s.createSession(ctx, user.Id, &device, &os)
}

func (s *AuthServiceImpl) CloseUsersSession(ctx context.Context, token string) error{
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return auth.ErrorWithLogout
	} else if session == nil{
		return auth.ErrorNonExistingSession
	}
	if num, err := s.authRepo.DeleteSessionById(ctx, session.Id); err != nil || num == 0{
		return auth.ErrorWithLogout
	} 
	authUser, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil || authUser == nil{
		return nil
	} 
	//if it's an anonymous account then delete both the account and the user
	if s.checkIfAnonymousAuth(authUser){
		if _, err := s.authRepo.DeleteUserAuthById(ctx, authUser.Id); err != nil{
			return auth.ErrorWithLogout 
		}
		if authUser.UserId != nil{
			if err := s.usersService.DeleteUserById(ctx, *authUser.UserId); err != nil{
				return auth.ErrorWithLogout 
			}
		}
	}
	return nil
}

func (s *AuthServiceImpl) CreateTempCouple(ctx context.Context, token string, startDate int) (int, error){
	userId, err := s.getUserIdFromSession(ctx, token)
	if err != nil{
		return 0, auth.ErrorcreatingTempCouple  
	}
	if userId == nil{
		return 0, auth.ErrorNoActiveUser
	}
	code, err := s.usersService.CreateTempCouple(ctx, *userId, startDate)
	if errors.Is(err, users.ErrorUserHasActiveCouple){
		return 0, auth.ErrCantCreateNewCouple
	}
	return code, err
}

func (s *AuthServiceImpl) CreateUser(ctx context.Context, token, firstName, lastName, gender, countryCode, languageCode string,birthDate int,) (string, error){
	//check token if its validate
	session, _ := s.authRepo.GetSessionByToken(ctx, token)
	if session != nil{
		userAuth, _ := s.authRepo.GetUserById(ctx, session.UserAuthId)
		if userAuth != nil && userAuth.UserId != nil{
			if s.checkIfAnonymousAuth(userAuth) {
				_, err := s.usersService.GetCoupleFromUser(ctx, *userAuth.UserId)
				if errors.Is(err, users.ErrorNoCoupleFound){
					//at this point, if there was an anoynmous account, and it had no couple, then we delete it
					go func (){ 
						s.usersService.DeleteUserById(context.Background(), *userAuth.UserId) 
						s.authRepo.DeleteSessionById(context.Background(), session.Id)
						s.usersService.DeleteUserById(context.Background(), userAuth.Id)
					}()
				}else{
					return "", auth.ErrorUserForAccountAlreadyExists
				}
			}
		}
	}

	// create user with users service (receives the userId)
	userId,  err := s.usersService.CreateUser(ctx, firstName, lastName, gender, countryCode, languageCode, birthDate)
	if err != nil{
		return "", err //sending other's domain ERROR 
	}
	//	if no token, create anonymous auth user and connect with user
	if session != nil{
		if num, err := s.authRepo.UpdateAuthUserId(ctx, session.UserAuthId, *userId); err != nil || num == 0{
			return "", auth.ErrorCreatingUser  
		}
	} else{
		authId := uuid.New()
		if num, err := s.authRepo.CreateEmptyUser(ctx, authId, *userId); err != nil || num == 0{
			return "", auth.ErrorCreatingUser 
		}
		return s.createSession(ctx, authId, nil, nil)
	}
	// if token, get the user auth and connect with userId
	return token, nil
}


func (s *AuthServiceImpl) CheckUserAuthStatus(ctx context.Context, token string) (string, error){
	userId, err := s.getUserIdFromSession(ctx, token)
	if err != nil{
		return "", auth.ErrorCheckingStatus
	}
	if userId == nil{
		return auth.StatusNoUserCreated, nil 
	}
	couple, _ := s.usersService.GetCoupleFromUser(ctx, *userId)
	if couple == nil{
		return auth.StatusUserCreated, nil 
	}
	partnerHasNickname, err := s.usersService.CheckPartnerNickname(ctx, *userId)
	if err == nil && !partnerHasNickname{
		return auth.StatusPartnerWithoutNickname, nil
	}
	return auth.StatusCoupleCreated, nil

}

func (s *AuthServiceImpl) ConnectCouple(ctx context.Context, token string, code int) (string, error) {
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return "", auth.ErrorUnableToConnectCouple 
	}else if session == nil{
		return "", auth.ErrorNonExistingSession
	}
	authUser, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil || authUser == nil{
		return "", auth.ErrorUnableToConnectCouple  
	}
	coupleId, partnerId, err := s.usersService.ConnectCouple(ctx, *authUser.UserId, code)
	if errors.Is(err, users.ErrorCantConnectWithYourself){
		return "", auth.ErrCantConnectWithYourself
	}else if errors.Is(err, users.ErrorUserHasActiveCouple){
		return "", auth.ErrCantCreateNewCouple
	}else if errors.Is(err, users.ErrorInvalidCode){
		return "", auth.ErrNonExistingCode
	}else if err != nil{
		return "",  auth.ErrorUnableToConnectCouple  
	}

	s.suscribersMutex.RLock()
	channel, ok := s.codeSuscribers[*partnerId]
	s.suscribersMutex.RUnlock()
	
	if ok{
		channel <- auth.StatusVinculated
	}
	return s.createAccessToken(*authUser.UserId, *coupleId, session.Id)
}


func (s *AuthServiceImpl) CreateAccessToken(ctx context.Context, token string)(string, *string, error){
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return "", nil, auth.ErrorNonExistingSession 
	}
	var newRefreshToken *string = nil
	//if expired, we create a new one
	if(session.ExpiresAt.Before(time.Now())){
		rToken, err := s.createSession(ctx, session.UserAuthId, session.Device, session.Os)
		if err != nil{
			go func(){s.authRepo.DeleteSessionById(ctx, session.Id)}()
			newRefreshToken = &rToken
		}
	}else{
		//in a separated go routine we update the last time used of the session
		go func(){
			s.authRepo.UpdateSessionLastUsed(context.Background(), session.Id, time.Now())
		}()
	}
	user, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil{
		return "", nil, auth.ErrorCreatingAccessToken 
	}
	if user.UserId == nil{
		return "", nil, auth.ErrorNoActiveUser
	}
	couple, err := s.usersService.GetCoupleFromUser(ctx, *user.UserId)
	if err != nil{
		return "", nil, auth.ErrorCreatingAccessToken 
	}
	if couple == nil {
		return "", nil, auth.ErrorNoActiveCoupleFromUser
	}
	accessToken, err := s.createAccessToken(*user.UserId, couple.Id, session.Id)
	return accessToken, newRefreshToken, err
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
		return auth.ErrorWithLogout 
	}
	authUser, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil || authUser == nil{
		return auth.ErrorWithLogout  
	}
	if s.checkIfAnonymousAuth(authUser){
		return auth.ErrorCantLogoutAnonymousAcc
	}
	if num, err := s.authRepo.DeleteSessionById(ctx, session.Id); err != nil || num == 0{
		return auth.ErrorWithLogout  
	}
	return nil
}


func (s *AuthServiceImpl) GetTempCoupleOfUser(ctx context.Context, token string)(*auth.TempCoupleModel, error){
	userId, err := s.getUserIdFromSession(ctx, token)
	if err != nil{
		return nil, auth.ErrorGettingTempCouple
	}
	tempCouple, err:= s.usersService.GetTempCoupleFromUser(ctx, *userId)
	if errors.Is(err, users.ErrorNoTempCoupleFound){
		return nil, auth.ErrTempCoupleNotFound
	}else if err != nil{
		return nil, auth.ErrorGettingTempCouple
	}
	tempCoupleAuth := new(auth.TempCoupleModel)
	*tempCoupleAuth = auth.TempCoupleModel{
		Code: tempCouple.Code,
		StartDate: tempCouple.StartDate,
	}
	return tempCoupleAuth, nil
}

func (s *AuthServiceImpl) RemoveCodeSuscriber(userId uuid.UUID){
	s.suscribersMutex.Lock()
	channel, ok := s.codeSuscribers[userId]
	if ok{
		delete(s.codeSuscribers, userId)
		close(channel)
	}
	s.suscribersMutex.Unlock()
}

func (s *AuthServiceImpl) SuscribeTempCoupleNot(ctx context.Context, token string)(chan string, *uuid.UUID, error){
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil{
		return nil, nil, auth.ErrUnableToSuscribe 
	}else if session == nil{
		return nil, nil, auth.ErrUnableToSuscribe
	}
	authUser, err := s.authRepo.GetUserById(ctx, session.UserAuthId)
	if err != nil || authUser == nil{
		return nil, nil, auth.ErrUnableToSuscribe  
	}
	_, err = s.usersService.GetTempCoupleFromUser(ctx, *authUser.UserId)
	if errors.Is(err, users.ErrorNoTempCoupleFound){
		return nil, nil, auth.ErrNoCodeToSuscribe
	}else if err != nil{
		return nil, nil, auth.ErrUnableToSuscribe
	}
	s.RemoveCodeSuscriber(*authUser.UserId)
	s.suscribersMutex.Lock()
	newChannel := make(chan string)
	s.codeSuscribers[*authUser.UserId] = newChannel
	s.suscribersMutex.Unlock()
	return newChannel, authUser.UserId,nil
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
			return "", auth.ErrorCreatingSession 
		}
		token = base64.URLEncoding.EncodeToString(randomBytes)
		session, _ := s.authRepo.GetSessionByToken(ctx, token)
		if session == nil{
			break
		}
	}
	//expires in 1 month
	session := auth.SessionModel{
		Id: uuid.New(),
		Token: token,
		Device: device,
		Os: os,
		ExpiresAt: time.Now().AddDate(0,0,30),
		CreatedAt: time.Now(),
		LastUsed: time.Now(),
		UserAuthId: authId,
	}
	num, err := s.authRepo.CreateSession(
		ctx, 
		&session,
	)
	if err != nil || num == 0{
		return "", auth.ErrorCreatingSession 
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
		return "", auth.ErrorCreatingAccessToken
	}
	return tokenString, nil
}

func (s *AuthServiceImpl) getUserIdFromSession(ctx context.Context, token string) (*uuid.UUID, error){
	session, err := s.authRepo.GetSessionByToken(ctx, token)
	if err != nil || session == nil{
		return nil, auth.ErrorNonExistingSession 
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