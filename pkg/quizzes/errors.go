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
	ErrRetrievingQuiz = errors.New("UNABLE_TO_GET_QUIZ")
	ErrInvalidQuestionOptions = errors.New("INVALID_QUESTION_OPTIONS")
	ErrDeletingCategory = errors.New("UNABLE_TO_DELETE_CATEGORY")
	ErrDeletingQuestion  = errors.New("UNABLE_TO_DELETE_QUESTION")
	ErrDeletingQuiz = errors.New("UNABLE_TO_DELETE_QUIZ")
	ErrQuestionNotFound = errors.New("QUESTION_NOT_FOUND")
	ErrUnathorizedToEditQuiz = errors.New("UNATHORIZED_TO_EDIT_QUIZ")
)