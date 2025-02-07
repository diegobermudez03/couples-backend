package quizzes

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type AdminService interface {
	CreateQuizCategory(ctx context.Context, name, description string, image io.Reader) error
	UpdateQuizCategory(ctx context.Context, id uuid.UUID, name, description string, image io.Reader) error
	//DeleteQuizCategory(ctx context.Context, id uuid.UUID) error
}


type CreateQuestionRequest struct{
	Question 		string 
	QType 			string 
	OptionsJson		map[string]any 
	StrategicAnswerId *uuid.UUID
	StrategicName 		*string 
	StrategicDescription *string 
}

type UserService interface{
	CreateQuiz(ctx context.Context, name, description, languageCode string, categoryId, userId *uuid.UUID, image io.Reader) error
	UpdateQuiz(ctx context.Context, quizId uuid.UUID, name, description, languageCode string, categoryId *uuid.UUID, image io.Reader) error
	CreateQuestion(ctx context.Context, quizId uuid.UUID, parameters CreateQuestionRequest, images map[string]io.Reader) error
	GetQuizById(ctx context.Context, quizId uuid.UUID)(*QuizPlainModel, error)
}

type QuizzesRepository interface{
	GetCategoryByName(ctx context.Context, name string)(*QuizCatPlainModel, error)
	GetCategoryById(ctx context.Context, id uuid.UUID)(*QuizCatPlainModel, error)
	CreateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)
	UpdateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)

	GetQuizById(ctx context.Context, id uuid.UUID) (*QuizPlainModel, error)
	CreateQuiz(ctx context.Context, quiz *QuizPlainModel) (int, error)
	UpdateQuiz(ctx context.Context, quiz *QuizPlainModel) (int, error)


	//GetQuestionsFromQuizId(ctx context.Context, quizId uuid.UUID) ([]QuestionPlainModel, error)
	GetMaxOrderQuestionFromQuiz(ctx context.Context, quizId uuid.UUID) (int, error)
	CreateQuestion(ctx context.Context, model *QuestionPlainModel) (int, error)

	GetStrategicTypeAnswerById(ctx context.Context, id uuid.UUID) (*StrategicAnswerModel, error)
	CreateStrategicTypeAnswer(ctx context.Context, model *StrategicAnswerModel) (int, error)

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