package quizzes

import (
	"context"
	"io"
)

type AdminService interface {
	CreateQuizCategory(ctx context.Context, name, description string, image io.Reader) error
}


type QuizzesRepository interface{
	GetCategoryByName(ctx context.Context, name string)(*QuizCatPlainModel, error)
	CreateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)
}


const DOMAIN_NAME = "quizzes"
const CATEGORIES = "categories"
