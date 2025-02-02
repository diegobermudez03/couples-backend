package quizzes

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("CATEGORY_ALREADY_EXISTS")
	ErrMissingCategoryAttributes = errors.New("MISSING_CATEGORY_ATTRIBUTES")
	ErrCreatingCategory = errors.New("UNABLE_TO_CREATE_CATEGORY")
	ErrInvalidImageType = errors.New("INVALID_IMAGE_TYPE")
)