package appauth

import (
	"errors"

	"github.com/diegobermudez03/couples-backend/pkg/auth/domainauth"
)

type AuthServiceImpl struct {
}

func NewAuthService() domainauth.AuthService{
	return &AuthServiceImpl{}
}


func(s *AuthServiceImpl) RegisterUser(email string, password string, device string, os string) (string, string, error){
	return "sdfgdfghfghfghgfhf", "dfgdgdhgfghf", errors.New("unable to create the user")
}