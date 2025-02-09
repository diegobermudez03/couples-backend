package repoquizzes

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diegobermudez03/couples-backend/pkg/infraestructure"
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
		`SELECT id, name, description, created_at, active, image_id
		FROM quiz_categories where name = $1 AND active = TRUE`, 
		name,
	)
	return r.rowToCategory(row)
}

func (r *QuizzesPostgresRepo) GetCategoryById(ctx context.Context, id uuid.UUID)(*quizzes.QuizCatPlainModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, description, created_at, ACTIVE, image_id
		FROM quiz_categories WHERE id = $1 AND active = TRUE`, 
		id,
	)
	return r.rowToCategory(row)
}


func (r *QuizzesPostgresRepo) CreateCategory(ctx context.Context, category *quizzes.QuizCatPlainModel) (int, error){
	executor := infraestructure.GetDBContext(ctx, r.db)
	result, err := executor.ExecContext(
		ctx, 
		`INSERT INTO quiz_categories(id, name, description, created_at, image_id, active)
		VALUES($1, $2, $3, $4, $5, $6)`,
		category.Id, category.Name, category.Description, category.CreatedAt, category.ImageId, category.Active,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}


func (r *QuizzesPostgresRepo) UpdateCategory(ctx context.Context, category *quizzes.QuizCatPlainModel) (int, error){
	executor := infraestructure.GetDBContext(ctx, r.db)
	result, err := executor.ExecContext(
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

func (r *QuizzesPostgresRepo) GetQuizById(ctx context.Context, id uuid.UUID) (*quizzes.QuizPlainModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, description, language_code, image_id, published, active, created_at, category_id, creator_id
		FROM quizzes WHERE id = $1`,
		id,
	)
	m := new(quizzes.QuizPlainModel)
	err := row.Scan(&m.Id, &m.Name, &m.Description, &m.LanguageCode, &m.ImageId, &m.Published, &m.Active, &m.CreatedAt, &m.CategoryId, &m.CreatorId)
	if errors.Is(err, sql.ErrNoRows){
		return nil, nil 
	}else if err != nil{
		return nil, err 
	}
	return m, nil
}

func (r *QuizzesPostgresRepo) CreateQuiz(ctx context.Context, quiz *quizzes.QuizPlainModel) (int, error){
	executor := infraestructure.GetDBContext(ctx, r.db)
	result, err := executor.ExecContext(
		ctx,
		`INSERT INTO quizzes(id, name, description, language_code, image_id, published, active, created_at, category_id, creator_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		quiz.Id, quiz.Name, quiz.Description, quiz.LanguageCode, quiz.ImageId, quiz.Published, quiz.Active,
		quiz.CreatedAt, quiz.CategoryId, quiz.CreatorId,
	)

	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}


func (r *QuizzesPostgresRepo) UpdateQuiz(ctx context.Context, quiz *quizzes.QuizPlainModel) (int, error){
	executor := infraestructure.GetDBContext(ctx, r.db)
	result, err := executor.ExecContext(
		ctx, 
		`UPDATE quizzes SET name = $1, description = $2, category_id = $3, image_id = $4, active = $5, published = $6, language_code = $7
		WHERE id = $8`,
		quiz.Name, quiz.Description, quiz.CategoryId, quiz.ImageId, quiz.Active, quiz.Published, quiz.LanguageCode,quiz.Id,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *QuizzesPostgresRepo) CreateQuestion(ctx context.Context, model *quizzes.QuestionPlainModel) (int, error){
	executor := infraestructure.GetDBContext(ctx, r.db)
	result, err := executor.ExecContext(
		ctx, 
		`INSERT INTO quiz_questions(id, ordering, question, question_type, options_json, quiz_id, strategic_answer_id, active)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		model.Id, model.Ordering, model.Question, model.QuestionType, model.OptionsJson, model.QuizId, model.StrategicAnswerId, model.Active,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *QuizzesPostgresRepo) CreateStrategicTypeAnswer(ctx context.Context, model *quizzes.StrategicAnswerModel) (int, error){
	executor := infraestructure.GetDBContext(ctx, r.db)
	result, err := executor.ExecContext(
		ctx, 
		`INSERT INTO strategic_type_answers(id, name, description)
		VALUES($1, $2, $3)`,
		model.Id, model.Name, model.Description,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *QuizzesPostgresRepo) GetMaxOrderQuestionFromQuiz(ctx context.Context, quizId uuid.UUID) (int, error){
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT COALESCE(max(ordering), 0)
		FROM quiz_questions WHERE quiz_id = $1`,
		quizId,
	)
	var num int 
	if err := row.Scan(&num); err != nil{
		return 0, err  
	}
	return num, nil
}


func (r *QuizzesPostgresRepo) GetStrategicTypeAnswerById(ctx context.Context, id uuid.UUID) (*quizzes.StrategicAnswerModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, description
		FROM strategic_type_answers WHERE id = $1`,
		id,
	)
	model := new(quizzes.StrategicAnswerModel)
	if err := row.Scan(&model.Id, &model.Name, &model.Description); err != nil{
		return nil, err 
	}
	return model, nil
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

func (r *QuizzesPostgresRepo) rowToCategory(row *sql.Row) (*quizzes.QuizCatPlainModel, error){
	model := new(quizzes.QuizCatPlainModel)
	err := row.Scan(&model.Id, &model.Name, &model.Description, &model.CreatedAt, &model.Active, &model.ImageId)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, nil 
		}
		return nil, err 
	}
	return model, nil
}