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
		token string,
	) (string, error)

	CheckUserExistance(ctx context.Context, token string) error
	CloseOngoingSession(ctx context.Context, token string) error
	CreateTempCouple(ctx context.Context, token string, startDate int) (int, error)
}

type UsersRepo interface{
	CreateUser(ctx context.Context, user *UserModel) error
	DeleteUser(ctx context.Context, userId uuid.UUID) error
	GetTempCoupleByCode(ctx context.Context, code int) (*TempCoupleModel, error)
	CheckTempCoupleById(ctx context.Context, userId uuid.UUID) (exists bool, err error)
	UpdateTempCouple(ctx context.Context, tempCouple *TempCoupleModel) error
	CreateTempCouple(ctx context.Context, tempCouple *TempCoupleModel) error
}