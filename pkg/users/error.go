package users

import "errors"

var (
	ErrorNoCoupleFound = errors.New("no couple was found")
	ErrorUserHasActiveCouple = errors.New("the user has an active couple, cant be deleted")
)