package auth

import "errors"

var (
	ErrorNoUserFoundEmail = errors.New("no user found with that email")
	ErrorNoActiveUser = errors.New("the user has no active user")
)