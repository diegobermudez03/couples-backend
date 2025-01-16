package appusers

import (
	"context"

	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/diegobermudez03/couples-backend/pkg/users"
)

type UsersServiceImpl struct {
	authService 	auth.AuthService
}

func NewUsersServiceImpl(authService auth.AuthService) users.UsersService {
	return &UsersServiceImpl{
		authService: authService,
	}
}

func (s *UsersServiceImpl) CreateUser(
	ctx context.Context,
	firstName, lastName, gender, countryCode, languageCode string,
	birthDate int,
	token string,
) (string, error) {
	return "", nil
}

func (s *UsersServiceImpl) CheckUserExistance(ctx context.Context, token string) error {
	_, err := s.authService.GetUserIdFromSession(ctx, token)
	return err
}