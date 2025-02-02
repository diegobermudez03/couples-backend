package repoquizzes

import (
	"context"

	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
)

type QuizzesPostgresRepo struct{}

func NewQuizzesPostgresRepo() quizzes.QuizzesRepository{
	return &QuizzesPostgresRepo{}
}


func (r *QuizzesPostgresRepo) GetCategoryByName(ctx context.Context, name string)(*quizzes.QuizCatPlainModel, error){
	return nil, nil
}
func (r *QuizzesPostgresRepo) CreateCategory(ctx context.Context, category *quizzes.QuizCatPlainModel) (int, error){
	return 0, nil
}