package repoquizzes

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/google/uuid"
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
	row := r.db.QueryRowContext(
		ctx,
		`SELECT q.id, q.name, q.description, q.created_at, q.image_id, q.active, f.bucket, f.grouping, f.object_key, f.created_at, f.type
		FROM quiz_categories q 
		INNER JOIN files f ON f.id = q.image_id 
		AND q.name = $1 AND q.active = TRUE`, 
		name,
	)
	return r.rowToCategory(row)
}

func (r *QuizzesPostgresRepo) GetCategoryById(ctx context.Context, id uuid.UUID)(*quizzes.QuizCatPlainModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT q.id, q.name, q.description, q.created_at, q.image_id, q.active, f.bucket, f.grouping, f.object_key, f.created_at, f.type
		FROM quiz_categories q 
		INNER JOIN files f ON f.id = q.image_id 
		AND q.id = $1 AND q.active = TRUE`, 
		id,
	)
	return r.rowToCategory(row)
}


func (r *QuizzesPostgresRepo) CreateCategory(ctx context.Context, category *quizzes.QuizCatPlainModel) (int, error){
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO quiz_categories(id, name, description, created_at, image_id, active)
		VALUES($1, $2, $3, $4, $5, $6)`,
		category.Id, category.Name, category.Description, category.CreatedAt, category.File.Id, category.Active,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}


func (r *QuizzesPostgresRepo) UpdateCategory(ctx context.Context, category *quizzes.QuizCatPlainModel) (int, error){
	result, err := r.db.ExecContext(
		ctx, 
		`UPDATE quiz_categories SET name = $1, description = $2 WHERE id = $3`,
		category.Name, category.Description, category.Id,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

func (r *QuizzesPostgresRepo) rowToCategory(row *sql.Row) (*quizzes.QuizCatPlainModel, error){
	model := new(quizzes.QuizCatPlainModel)
	model.File = new(files.FileModel)
	err := row.Scan(&model.Id, &model.Name, &model.Description, &model.CreatedAt, &model.File.Id, &model.Active,
		&model.File.Bucket, &model.File.Group, &model.File.ObjectKey, &model.File.CreatedAt, &model.File.Type,
	)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, nil 
		}
		return nil, err 
	}
	return model, nil
}