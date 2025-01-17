package appusers

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/diegobermudez03/couples-backend/pkg/localization"
	"github.com/diegobermudez03/couples-backend/pkg/users"
	"github.com/google/uuid"
)

/////////  HELPERS
var genders = map[string]bool{
	"male" : true,
	"female" : true,
}

var (
	errorInvalidGender = errors.New("invalid gender")
	errorTooYoung = errors.New("the user must be at least 12 years old")
)
///////////////////////////////////////////

type UsersServiceImpl struct {
	authService 		auth.AuthService
	localizationService localization.LocalizationService
	usersRepo 			users.UsersRepo
}

func NewUsersServiceImpl(authService auth.AuthService,localizationService localization.LocalizationService, usersRepo users.UsersRepo) users.UsersService {
	return &UsersServiceImpl{
		authService: authService,
		usersRepo: usersRepo,
		localizationService: localizationService,
	}
}

func (s *UsersServiceImpl) CreateUser(
	ctx context.Context,
	firstName, lastName, gender, countryCode, languageCode string,
	birthDate int,
	token string,
) (string, error) {
	lowerGender := strings.ToLower(gender)
	if _, ok := genders[lowerGender]; !ok{
		return "", errorInvalidGender
	}
	t := time.Unix(int64(birthDate), 0)

	// if the users is less than 10 years old
	if !time.Now().AddDate(-12, 0, 0).After(t){
		return "", errorTooYoung
	}

	//check lang and country
	err1, err2 := s.localizationService.ValidateCountry(countryCode), s.localizationService.ValidateLanguage(languageCode)
	if err1 != nil{
		return "", err1 
	}
	if err2 != nil{
		return "", err2
	}


	userId := uuid.New()
	err := s.usersRepo.CreateUser(
		ctx, 
		&users.UserModel{
			Id: userId,
			FirstName: firstName,
			LastName: lastName,
			Gender: lowerGender,
			BirthDate: t,
			CreatedAt: time.Now(),
			Active: true,
			CountryCode: countryCode,
			LanguageCode: languageCode,
		},
	)
	if err != nil{
		return "", err
	}

	//if the user had already a session we end here
	if token != ""{
		if err := s.authService.VinculateAuthWithUser(ctx, token, userId); err != nil{
			return "", err
		}
		return token, nil
	}
	
	//if the user doesn't have a session, then we create it
	return s.authService.CreateAnonymousSession(ctx, userId)
}

func (s *UsersServiceImpl) CheckUserExistance(ctx context.Context, token string) error {
	_, err := s.authService.GetUserIdFromSession(ctx, token)
	return err
}

func (s *UsersServiceImpl) CloseOngoingSession(ctx context.Context, token string) error{
	userId, err := s.authService.CloseSession(ctx, token)
	if err != nil{
		return err
	}
	if userId == nil{
		return nil 
	}

	if err := s.usersRepo.DeleteUser(ctx, *userId); err != nil{
		return err 
	}
	return nil
}