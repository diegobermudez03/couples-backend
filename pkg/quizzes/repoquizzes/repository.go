package repoquizzes

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
)

type QuizzesPostgresRepo struct{
	db 		*sql.DB
}

func NewQuizzesPostgresRepo(db *sql.DB) quizzes.QuizzesRepository{
	return &QuizzesPostgresRepo{
		db: db,
	}
}


func (r *QuizzesPostgresRepo) GetCategoryByName(ctx context.Context, name string)(*quizzes.QuizCatPlainModel, error){
	result := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, description, created_at, image_id
		FROM quiz_categories WHERE name = $1`, 
		name,
	)
	model := new(quizzes.QuizCatPlainModel)
	if err := result.Scan(&model.Id, &model.Name, &model.Description, &model.CreatedAt, &model.ImageId); err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, nil 
		}
		return nil, err 
	}
	return model, nil
}
func (r *QuizzesPostgresRepo) CreateCategory(ctx context.Context, category *quizzes.QuizCatPlainModel) (int, error){
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO quiz_categories(id, name, description, created_at, image_id)
		VALUES($1, $2, $3, $4, $5)`,
		category.Id, category.Name, category.Description, category.CreatedAt, category.ImageId,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}