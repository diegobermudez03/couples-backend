package quizzes

import (
	"context"
	"io"

	"github.com/google/uuid"
)


type AdminService interface {
	CreateQuizCategory(ctx context.Context, name, description string, image io.Reader) error
	UpdateQuizCategory(ctx context.Context, id uuid.UUID, name, description string, image io.Reader) error
	DeleteQuizCategory(ctx context.Context, id uuid.UUID) error
}

type UserService interface{
	AuthorizeQuizCreator(ctx context.Context, quizId *uuid.UUID, questionId *uuid.UUID, userId uuid.UUID) error

	GetCategories(ctx context.Context, filters FetchFilters)([]QuizCatModel, error)
	
	CreateQuiz(ctx context.Context, name, description, languageCode string, categoryId, userId *uuid.UUID, image io.Reader) (*uuid.UUID, error)
	UpdateQuiz(ctx context.Context, quizId uuid.UUID, name, description, languageCode string, categoryId *uuid.UUID, image io.Reader) error
	DeleteQuiz(ctx context.Context, quizId uuid.UUID) error

	CreateQuestion(ctx context.Context, quizId uuid.UUID, parameters CreateQuestionRequest, images map[string]io.Reader)(*uuid.UUID, error)
	UpdateQuestion(ctx context.Context, questionId uuid.UUID, parameters UpdateQuestionRequest, images map[string]io.Reader) error
	DeleteQuestion(ctx context.Context, questionId uuid.UUID) error
}

type FetchFilters struct{
	Limit 	*int 
	Page 	*int 
}

type CreateQuestionRequest struct{
	Question 		string 
	QType 			string 
	OptionsJson		map[string]any 
	StrategicAnswerId *uuid.UUID
	StrategicName 		*string 
	StrategicDescription *string 
}

type UpdateQuestionRequest struct{
	Question 			*string 
	OptionsJson			map[string]any 
	StrategicAnswerId 	*uuid.UUID
	StrategicName 		*string 
	StrategicDescription *string 
}


const DOMAIN_NAME = "quizzes"
const CATEGORIES = "categories"
const QUIZZES = "quizzess"
const PROFILE = "profile"


//QUESTION TYPES
const (
	TRUE_FALSE_TYPE = "TRUE_FALSE"
	SLIDER_TYPE = "SLIDER"
	ORDERING_TYPE = "ORDERING"
	OPEN_TYPE = "OPEN"
	MULTIPLE_CH_TYPE = "MULTIPLE_CH"
	MATCHING_TYPE = "MATCHING"
	DRAG_AND_DROP_TYPE = "DRAG_AND_DROP"
)


//PLACEHOLDERS IN QUESTION
const YOU_PLACEHOLDER = "%y%"
const PARTNER_PLACEHOLDER = "%r%"


//sorting types
const LEAST_TO_MOST = "L-M"
const MOST_TO_LEAST = "M-L"