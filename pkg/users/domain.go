package users

import (
	"context"

	"github.com/google/uuid"
)

type UsersService interface {
	CreateUser(
		ctx context.Context, 
		firstName, lastName, gender, countryCode, languageCode string,
		birthDate int,
	) (*uuid.UUID, error) 
	CreateTempCouple(ctx context.Context, userId uuid.UUID, startDate int) (int, error)
	DeleteUserById(ctx context.Context, userId uuid.UUID) error
	GetCoupleFromUser(ctx context.Context, userId uuid.UUID) (*CoupleModel, error)
	ConnectCouple(ctx context.Context, userId uuid.UUID, code int)(*uuid.UUID, error)
	EditPartnersNickname(ctx context.Context, userId uuid.UUID, coupleId uuid.UUID, nickname string) error
}

type UsersRepo interface{
	CreateUser(ctx context.Context, user *UserModel) error
	DeleteUserById(ctx context.Context, userId uuid.UUID) error
	GetTempCoupleByCode(ctx context.Context, code int) (*TempCoupleModel, error)
	CheckTempCoupleById(ctx context.Context, userId uuid.UUID) (exists bool, err error)
	UpdateTempCouple(ctx context.Context, tempCouple *TempCoupleModel) error
	CreateTempCouple(ctx context.Context, tempCouple *TempCoupleModel) error
	GetCoupleByUserId(ctx context.Context, userId uuid.UUID) (*CoupleModel, error)
	DeleteTempCoupleById(ctx context.Context, id uuid.UUID) error
	GetUserById(ctx context.Context, userId uuid.UUID) (*UserModel, error)
	CreateCouple(ctx context.Context, couple *CoupleModel) error
	CreateCouplePoints(ctx context.Context, points *PointsModel) error
	GetCoupleById(ctx context.Context, coupleId uuid.UUID) (*CoupleModel, error)
	UpdateUserNicknameById(ctx context.Context, userId uuid.UUID, nickname string) error
}



///////////////////////// POINTS
const CouplePointsForConnecting = 50