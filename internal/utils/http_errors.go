package utils

import "errors"

var (
	ErrNoTokenProvided = errors.New("NO_TOKEN_PROVIDED")
	ErrNoPathGiven = errors.New("NO_PATH_GIVEN")
	ErrUnbaleToLoad = errors.New("UNABLE_TO_LOAD")
	ErrFileTooBig = errors.New("FILE_TOO_BIG")
	ErrMissingFields = errors.New("MISSING_FIELDS")
	ErrInvalidId = errors.New("INVALID_ID")
	ErrEmptyCategoryId = errors.New("EMPTY_CATEGORY_ID")
	ErrEmptyQuizId = errors.New("EMPTY_QUIZ_ID")
	ErrEmptyQuestionId = errors.New("EMPTY_QUESTION_ID")
)
