package quizzes

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("CATEGORY_ALREADY_EXISTS")
	ErrMissingCategoryAttributes = errors.New("MISSING_CATEGORY_ATTRIBUTES")
	ErrCreatingCategory = errors.New("UNABLE_TO_CREATE_CATEGORY")
	ErrInvalidImageType = errors.New("INVALID_IMAGE_TYPE")
	ErrUpdatingCategory = errors.New("UNABLE_TO_UPDATE_CATEGORY")
	ErrNonExistingCategory = errors.New("NON_EXISTING_CATEGORY")
	ErrCreatingQuiz = errors.New("UNABLE_TO_CREATE_QUIZ")
	ErrCategoryDontExists = errors.New("CATEGORY_DOESNT_EXIST")
	ErrEmptyQuizName = errors.New("EMPTY_QUIZ_NAME")
	ErrInvalidLanguage = errors.New("INVALID_LANGUAGE")
	ErrUpdatingQuiz = errors.New("UNABLE_TO_UPDATE_QUIZ")
	ErrQuizNotFound = errors.New("QUIZ_NOT_FOUND")
	ErrInvalidQuestionType = errors.New("INVALID_QUESTION_TYPE")
	ErrCreatingQuestion = errors.New("UNABLE_TO_CREATE_QUESTION")
	ErrInvalidQuestionOptions = errors.New("INVALID_QUESTION_OPTIONS")
)