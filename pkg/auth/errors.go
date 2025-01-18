package auth

import "errors"

var (
	ErrorNoUserFoundEmail = errors.New("no user found with that email")
	ErrorNoActiveUser = errors.New("the user has no active user")
	ErrorNoActiveCoupleFromUser = errors.New("the user has no active couple")
	ErrorExpiredRefreshToken = errors.New("the refresh token is expired")
	ErrorExpiredAccessToken = errors.New("the access token is expired")
	ErrorMalformedAccessToken = errors.New("malformed access token")
	ErrorCantLogoutAnonymousAcc = errors.New("cant logout from an anonymous account")
)