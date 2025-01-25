package users

import "errors"

var (
	ErrorInvalidCreateUser = errors.New("INVALID_TO_CREATE_USER")
	ErrorInvalidCountryCode = errors.New("INVALID_COUNTRY_CODE")
	ErrorInvalidLanguageCode = errors.New("INVALID_LANG_CODE")
	ErrorInvalidGender = errors.New("INVALID_GENDER")
	ErrorTooYoung = errors.New("TOO_YOUNG")
	ErrorNoCoupleFound = errors.New("no couple was found")
	ErrorUserHasActiveCouple = errors.New("USER_HAS_ACTIVE_COUPLE")
	ErrorDeletingUser = errors.New("UNABLE_TO_DELETE_USER")
	ErrorCreatingTempCouple = errors.New("UNABLE_TO_CREATE_CODE")
	ErrorInvalidCode = errors.New("INVALID_CODE")
	ErrorCantConnectWithYourself =  errors.New("CANT_CONNECT_WITH_YOURSELF")
	ErrorConnectingCouple = errors.New("UNABLE_TO_CONNECT_COUPLE")
	ErrorCreatingPoints = errors.New("UNABLE_TO_ADD_POINTS")
	ErrorUpdatingNickname = errors.New("UNABLE_TO_EDIT_NICKNAME")
)