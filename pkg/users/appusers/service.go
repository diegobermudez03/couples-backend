package appusers

import (
	"context"
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
		return nil, users.ErrorInvalidGender
	}
	t := time.Unix(int64(birthDate), 0)

	// if the users is less than 10 years old
	if !time.Now().AddDate(-14, 0, 0).After(t){
		return nil, users.ErrorTooYoung
	}

	//check lang and country
	err1, err2 := s.localizationService.ValidateCountry(countryCode), s.localizationService.ValidateLanguage(languageCode)
	if err1 != nil{
		return nil, users.ErrorInvalidCountryCode 
	}
	if err2 != nil{
		return nil, users.ErrorInvalidLanguageCode
	}

	userId := new(uuid.UUID)
	*userId = uuid.New()
	num, err := s.usersRepo.CreateUser(
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
	if err != nil || num == 0 {
		return nil, users.ErrorUnableCreateUser
	}
	return userId, nil
}

func (s *UsersServiceImpl) DeleteUserById(ctx context.Context, userId uuid.UUID) error{
	couple, err := s.usersRepo.GetCoupleByUserId(ctx, userId)
	if err != nil{
		return users.ErrorDeletingUser
	}else if couple != nil{
		return users.ErrorUserHasActiveCouple 
	}
	if num, err := s.usersRepo.DeleteUserById(ctx, userId); err != nil || num ==0{
		return users.ErrorDeletingUser 
	}
	s.usersRepo.DeleteTempCoupleById(ctx, userId)
	return nil
}

func (s *UsersServiceImpl) CreateTempCouple(ctx context.Context, userId uuid.UUID, startDate int) (int, error){
	couple, _ := s.usersRepo.GetCoupleByUserId(ctx, userId)
	if couple != nil{
		return 0, users.ErrorUserHasActiveCouple
	}
	var code int 
	//unique code creation
	for{
		code = rand.Intn(89999) + 10000
		if couple,err := s.usersRepo.GetTempCoupleByCode(ctx, code); err == nil && couple == nil{
			break
		}
	}
	exists, err := s.usersRepo.CheckTempCoupleById(ctx, userId)
	if err != nil{
		return 0, users.ErrorCreatingTempCouple 
	}

	tempCouple := users.TempCoupleModel{
		UserId: userId,
		Code: code,
		StartDate: time.Unix(int64(startDate), 0),

	}
	var num int
	if exists{
		num, err = s.usersRepo.UpdateTempCouple(ctx, &tempCouple)
	}else{
		num, err = s.usersRepo.CreateTempCouple(ctx, &tempCouple)
	}

	if err != nil || num == 0{
		return 0,  users.ErrorCreatingTempCouple  
	}
	return code, nil
}

func (s *UsersServiceImpl) GetCoupleFromUser(ctx context.Context, userId uuid.UUID) (*users.CoupleModel, error){
	return s.usersRepo.GetCoupleByUserId(ctx, userId)
}


func (s *UsersServiceImpl) ConnectCouple(ctx context.Context, userId uuid.UUID, code int) (*uuid.UUID, *uuid.UUID, error){
	// check that the user doesn't have a couple
	coupleCheck, _ := s.usersRepo.GetCoupleByUserId(ctx, userId)
	if coupleCheck != nil{
		return nil, nil, users.ErrorUserHasActiveCouple
	}
	
	tempCouple, _ := s.usersRepo.GetTempCoupleByCode(ctx, code)
	if tempCouple == nil{
		return nil, nil, users.ErrorInvalidCode
	}
	//check that the user isn't connecting with himself
	if userId == tempCouple.UserId{
		return nil, nil, users.ErrorCantConnectWithYourself
	}
	//create the couple
	user1, err := s.usersRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, nil, users.ErrorConnectingCouple 
	}
	user2, err := s.usersRepo.GetUserById(ctx, tempCouple.UserId)
	if err != nil {
		return nil, nil, users.ErrorConnectingCouple 
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
	if num, err := s.usersRepo.CreateCouple(ctx, couple); err != nil || num == 0{
		return nil, nil, users.ErrorConnectingCouple 
	}

	//delete temp couples
	s.usersRepo.DeleteTempCoupleById(ctx, heId)
	s.usersRepo.DeleteTempCoupleById(ctx, sheId)

	//create first points
	num, err := s.usersRepo.CreateCouplePoints(
		ctx,
		&users.PointsModel{
			Id: uuid.New(),
			Day: time.Now(),
			Points: users.CouplePointsForConnecting,
			CoupleId: &coupleId,
		},
	)
	if err != nil || num == 0{
		return nil, nil, users.ErrorCreatingPoints
	}
	return &coupleId, &tempCouple.UserId, nil
}

func (s *UsersServiceImpl) EditPartnersNickname(ctx context.Context, userId uuid.UUID, coupleId uuid.UUID, nickname string) error{
	couple, err := s.usersRepo.GetCoupleById(ctx, coupleId)
	if err != nil{
		return users.ErrorUpdatingNickname
	}
	var partnerId uuid.UUID
	if couple.HeId == userId{
		partnerId = couple.SheId
	}else{
		partnerId = couple.HeId
	}

	if num, err := s.usersRepo.UpdateUserNicknameById(ctx, partnerId, nickname); err != nil || num == 0{
		return users.ErrorUpdatingNickname 
	}
	return nil
}


func (s *UsersServiceImpl)  CheckPartnerNickname(ctx context.Context, userId uuid.UUID) (hasNickname bool, err error){
	couple, err := s.usersRepo.GetCoupleByUserId(ctx, userId)
	if err != nil{
		return false, users.ErrorUnableToCheckPartnerNickname
	}
	var partnerId uuid.UUID
	if couple.HeId == userId{
		partnerId = couple.SheId
	}else{
		partnerId = couple.HeId
	}
	partner, err := s.usersRepo.GetUserById(ctx, partnerId)
	if err != nil{
		return false, users.ErrorUnableToCheckPartnerNickname
	}
	return partner.NickName != "", nil
}


func(s *UsersServiceImpl)  GetTempCoupleFromUser(ctx context.Context, userId uuid.UUID)(*users.TempCoupleModel, error){
	tempCouple, err := s.usersRepo.GetTempCoupleFromUser(ctx, userId)
	if err != nil{
		return nil, users.ErrorUnableToGetTempCouple
	}else if tempCouple == nil{
		return nil, users.ErrorNoTempCoupleFound
	}
	return tempCouple, nil
}