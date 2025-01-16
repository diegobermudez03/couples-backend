package users

import "context"

type UsersService interface {
	CreateUser(
		ctx context.Context, 
		firstName, lastName, gender, countryCode, languageCode string,
		birthDate int,
		token string,
	) (string, error)

	CheckUserExistance(ctx context.Context, token string) error
}