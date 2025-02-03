package quizzes

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type AdminService interface {
	CreateQuizCategory(ctx context.Context, name, description string, image io.Reader) error
	UpdateQuizCategory(ctx context.Context, id uuid.UUID, name, description string, image io.Reader) error
	CreateQuiz(ctx context.Context, name, description, languageCode string, categoryId uuid.UUID, image io.Reader) error
	//DeleteQuizCategory(ctx context.Context, id uuid.UUID) error
}


type QuizzesRepository interface{
	GetCategoryByName(ctx context.Context, name string)(*QuizCatPlainModel, error)
	GetCategoryById(ctx context.Context, id uuid.UUID)(*QuizCatPlainModel, error)
	CreateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)
	UpdateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)

	CreateQuiz(ctx context.Context, quiz *QuizPlainModel) (int, error)

}


const DOMAIN_NAME = "quizzes"
const CATEGORIES = "categories"
const QUIZZES = "quizzess"
const PROFILE = "profile"
