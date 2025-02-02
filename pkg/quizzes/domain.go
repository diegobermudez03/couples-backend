package quizzes

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type AdminService interface {
	CreateQuizCategory(ctx context.Context, name, description string, image io.Reader) error
	UpdateQuizCategory(ctx context.Context, id uuid.UUID, name, description string, image io.Reader) error
}


type QuizzesRepository interface{
	GetCategoryByName(ctx context.Context, name string)(*QuizCatPlainModel, error)
	GetCategoryById(ctx context.Context, id uuid.UUID)(*QuizCatPlainModel, error)
	CreateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)
	UpdateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)
}


const DOMAIN_NAME = "quizzes"
const CATEGORIES = "categories"
