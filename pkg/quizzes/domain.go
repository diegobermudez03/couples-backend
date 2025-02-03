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
	UpdateQuiz(ctx context.Context, quizId uuid.UUID, name, description string, categoryId *uuid.UUID, image io.Reader) error
	//DeleteQuizCategory(ctx context.Context, id uuid.UUID) error
}


type QuizzesRepository interface{
	GetCategoryByName(ctx context.Context, name string)(*QuizCatPlainModel, error)
	GetCategoryById(ctx context.Context, id uuid.UUID)(*QuizCatPlainModel, error)
	CreateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)
	UpdateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)

	GetQuizById(ctx context.Context, id uuid.UUID) (*QuizPlainModel, error)
	CreateQuiz(ctx context.Context, quiz *QuizPlainModel) (int, error)
	UpdateQuiz(ctx context.Context, quiz *QuizPlainModel) (int, error)

}


const DOMAIN_NAME = "quizzes"
const CATEGORIES = "categories"
const QUIZZES = "quizzess"
const PROFILE = "profile"
