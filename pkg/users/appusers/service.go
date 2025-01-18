package appusers

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/localization"
	"github.com/diegobermudez03/couples-backend/pkg/users"
	"github.com/google/uuid"
)

const (
	maleGender = "male"
	femaleGender = "female"
)

/////////  HELPERS
var genders = map[string]bool{
	maleGender : true,
	femaleGender : true,
}

var (
	errorInvalidGender = errors.New("invalid gender")
	errorTooYoung = errors.New("the user must be at least 12 years old")
	errorUserAlreadyHasCouple = errors.New("the user already has a couple")
	errorInvalidCode = errors.New("the couple code is invalid")
	errorCantConnectWithYourself = errors.New("you cant connect with yourself")
)
///////////////////////////////////////////

type UsersServiceImpl struct {
	localizationService localization.LocalizationService
	usersRepo 			users.UsersRepo
}

func NewUsersServiceImpl(localizationService localization.LocalizationService, usersRepo users.UsersRepo) users.UsersService {
	return &UsersServiceImpl{
		usersRepo: usersRepo,
		localizationService: localizationService,
	}
}

func (s *UsersServiceImpl) CreateUser(
	ctx context.Context,
	firstName, lastName, gender, countryCode, languageCode string,
	birthDate int,
) (*uuid.UUID, error) {
	lowerGender := strings.ToLower(gender)
	if _, ok := genders[lowerGender]; !ok{
		return nil, errorInvalidGender
	}
	t := time.Unix(int64(birthDate), 0)

	// if the users is less than 10 years old
	if !time.Now().AddDate(-12, 0, 0).After(t){
		return nil, errorTooYoung
	}

	//check lang and country
	err1, err2 := s.localizationService.ValidateCountry(countryCode), s.localizationService.ValidateLanguage(languageCode)
	if err1 != nil{
		return nil, err1 
	}
	if err2 != nil{
		return nil, err2
	}

	userId := new(uuid.UUID)
	*userId = uuid.New()
	err := s.usersRepo.CreateUser(
		ctx, 
		&users.UserModel{
			Id: *userId,
			FirstName: firstName,
			LastName: lastName,
			Gender: lowerGender,
			BirthDate: t,
			CreatedAt: time.Now(),
			Active: true,
			CountryCode: countryCode,
			LanguageCode: languageCode,
			NickName: firstName,
		},
	)
	if err != nil{
		return nil, err
	}
	return userId, nil
}

func (s *UsersServiceImpl) DeleteUserById(ctx context.Context, userId uuid.UUID) error{
	couple, err := s.usersRepo.GetCoupleByUserId(ctx, userId)
	if !errors.Is(err, users.ErrorNoCoupleFound){
		if couple != nil{
			return users.ErrorUserHasActiveCouple 
		}else{
			return err
		}
	}
	if err := s.usersRepo.DeleteUserById(ctx, userId); err != nil{
		return err 
	}
	s.usersRepo.DeleteTempCoupleById(ctx, userId)
	return nil
}

func (s *UsersServiceImpl) CreateTempCouple(ctx context.Context, userId uuid.UUID, startDate int) (int, error){
	couple, _ := s.usersRepo.GetCoupleByUserId(ctx, userId)
	if couple != nil{
		return 0, errorUserAlreadyHasCouple
	}
	var code int 
	//unique code creation
	for{
		code = rand.Intn(89999) + 10000
		if _,err := s.usersRepo.GetTempCoupleByCode(ctx, code); err != nil{
			break
		}
	}
	exists, err:= s.usersRepo.CheckTempCoupleById(ctx, userId)
	if err != nil{
		return 0, err 
	}

	tempCouple := users.TempCoupleModel{
		UserId: userId,
		Code: code,
		StartDate: time.Unix(int64(startDate), 0),

	}
	if exists{
		err = s.usersRepo.UpdateTempCouple(ctx, &tempCouple)
	}else{
		err = s.usersRepo.CreateTempCouple(ctx, &tempCouple)
	}

	if err != nil{
		return 0, err 
	}
	return code, nil
}

func (s *UsersServiceImpl) GetCoupleFromUser(ctx context.Context, userId uuid.UUID) (*users.CoupleModel, error){
	return s.usersRepo.GetCoupleByUserId(ctx, userId)
}


func (s *UsersServiceImpl) ConnectCouple(ctx context.Context, userId uuid.UUID, code int) (*uuid.UUID, error){
	// check that the user doesn't have a couple
	coupleCheck, _ := s.usersRepo.GetCoupleByUserId(ctx, userId)
	if coupleCheck != nil{
		return nil, errorUserAlreadyHasCouple
	}
	
	tempCouple, _ := s.usersRepo.GetTempCoupleByCode(ctx, code)
	if tempCouple == nil{
		return nil, errorInvalidCode
	}
	//check that the user isn't connecting with himself
	if userId == tempCouple.UserId{
		return nil, errorCantConnectWithYourself
	}
	//create the couple
	user1, err := s.usersRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err 
	}
	user2, err := s.usersRepo.GetUserById(ctx, tempCouple.UserId)
	if err != nil {
		return nil, err 
	}

	var heId uuid.UUID
	var sheId uuid.UUID
	if user1.Gender == maleGender{
		heId = user1.Id
		sheId = user2.Id
	}else{
		sheId = user1.Id
		heId = user2.Id
	}

	coupleId := uuid.New()
	couple := &users.CoupleModel{
		Id: coupleId,
		RelationStart: tempCouple.StartDate,
		HeId: heId,
		SheId: sheId,
	}
	if err := s.usersRepo.CreateCouple(ctx, couple); err != nil{
		return nil, err 
	}

	//delete temp couples
	s.usersRepo.DeleteTempCoupleById(ctx, heId)
	s.usersRepo.DeleteTempCoupleById(ctx, sheId)

	//create first points
	err = s.usersRepo.CreateCouplePoints(
		ctx,
		&users.PointsModel{
			Id: uuid.New(),
			Day: time.Now(),
			Points: users.CouplePointsForConnecting,
			CoupleId: &coupleId,
		},
	)
	if err != nil{
		return nil, err
	}
	return &coupleId, nil
}

func (s *UsersServiceImpl) EditPartnersNickname(ctx context.Context, userId uuid.UUID, coupleId uuid.UUID, nickname string) error{
	couple, err := s.usersRepo.GetCoupleById(ctx, coupleId)
	if err != nil{
		return err 
	}
	var partnerId uuid.UUID
	if couple.HeId == userId{
		partnerId = couple.SheId
	}else{
		partnerId = couple.HeId
	}

	if err := s.usersRepo.UpdateUserNicknameById(ctx, partnerId, nickname); err != nil{
		return err 
	}
	return nil
}