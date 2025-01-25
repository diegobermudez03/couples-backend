package auth

import "errors"

var (
	ErrorCreatingSession = errors.New("UNABLE_TO_CREATE_SESSION")
	ErrorCreatingAccessToken = errors.New("UNABLE_TO_CREATE_ACCESS_TOKEN")
	ErrorCheckingStatus = errors.New("UNABLE_TO_GET_STATUS")
	ErrorcreatingTempCouple = errors.New("UNABLE_TO_CREATE_TEMP_COUPLE")
	ErrorNoUserFoundEmail = errors.New("ERROR_NO_USER_FOUND_EMAIL")
	ErrorNoActiveUser = errors.New("ACCOUNT_HAS_NO_ACTIVE_USER")
	ErrorNoActiveCoupleFromUser = errors.New("ACCOUNT_HAS_NO_ACTIVE_COUPLE")
	ErrorExpiredRefreshToken = errors.New("EXPIRED_REFRESH_TOKEN")
	ErrorExpiredAccessToken = errors.New("EXPIRED_ACCESS_TOKEN")
	ErrorMalformedAccessToken = errors.New("MALFORMED_ACCESS_TOKEN")
	ErrorCantLogoutAnonymousAcc = errors.New("UNABLE_TO_LOGOUT_FROM_ANONYMOUS")
	ErrorIncorrectPassword = errors.New("INCORRECT_PASSWORD")
	ErrorInsecurePassword = errors.New("INSECURE_PASSWORD")
	ErrorEmailAlreadyUsed = errors.New("EMAIL_ALREADY_USED")
	ErrorCreatingAccount = errors.New("UNABLE_TO_CREATE_ACCOUNT")
	ErrorVinculatingAccount = errors.New("UNABLE_TO_VINCULATE_ACCOUNT")
	ErrorCreatingUser = errors.New("UNABLE_TO_CREATE_USER")
	ErrorWithLogin = errors.New("UNABLE_TO_LOGIN")
	ErrorInvalidRefreshToken = errors.New("INVALID_REFRESH_TOKEN")
	ErrorWithLogout = errors.New("UNABLE_TO_LOGOUT")
	ErrorNonExistingSession = errors.New("NON_EXISTING_SESSION")
	ErrorUserForAccountAlreadyExists = errors.New("ACCOUNT_ALREADY_HAS_USER")
	ErrorUnableToConnectCouple = errors.New("UNABLE_TO_CONNECT_COUPLE")
)