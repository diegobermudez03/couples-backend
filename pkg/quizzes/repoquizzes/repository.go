package repoquizzes

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx, 
			`INSERT INTO quiz_categories(id, name, description, created_at, image_id, active)
			VALUES($1, $2, $3, $4, $5, $6)`,
			category.Id, category.Name, category.Description, category.CreatedAt, category.ImageId, category.Active,
		)
	})
}


func (r *QuizzesPostgresRepo) UpdateCategory(ctx context.Context, category *quizzes.QuizCatPlainModel) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx, 
			`UPDATE quiz_categories SET name = $1, description = $2 WHERE id = $3`,
			category.Name, category.Description, category.Id,
		)
	})
}

func (r *QuizzesPostgresRepo) GetQuizById(ctx context.Context, id uuid.UUID) (*quizzes.QuizPlainModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, description, language_code, image_id, published, active, created_at, category_id, creator_id
		FROM quizzes WHERE id = $1 AND active=TRUE`,
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
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx,
			`INSERT INTO quizzes(id, name, description, language_code, image_id, published, active, created_at, category_id, creator_id)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			quiz.Id, quiz.Name, quiz.Description, quiz.LanguageCode, quiz.ImageId, quiz.Published, quiz.Active,
			quiz.CreatedAt, quiz.CategoryId, quiz.CreatorId,
		)
	})
}


func (r *QuizzesPostgresRepo) UpdateQuiz(ctx context.Context, quiz *quizzes.QuizPlainModel) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx, 
			`UPDATE quizzes SET name = $1, description = $2, category_id = $3, image_id = $4, active = $5, published = $6, language_code = $7
			WHERE id = $8`,
			quiz.Name, quiz.Description, quiz.CategoryId, quiz.ImageId, quiz.Active, quiz.Published, quiz.LanguageCode,quiz.Id,
		)
	})
}

func (r *QuizzesPostgresRepo) CreateQuestion(ctx context.Context, model *quizzes.QuestionPlainModel) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx, 
			`INSERT INTO quiz_questions(id, ordering, question, question_type, options_json, quiz_id, strategic_answer_id, active)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
			model.Id, model.Ordering, model.Question, model.QuestionType, model.OptionsJson, model.QuizId, model.StrategicAnswerId, model.Active,
		)
	})
}

func (r *QuizzesPostgresRepo) CreateStrategicTypeAnswer(ctx context.Context, model *quizzes.StrategicAnswerModel) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx, 
			`INSERT INTO strategic_type_answers(id, name, description)
			VALUES($1, $2, $3)`,
			model.Id, model.Name, model.Description,
		)
	})
}

func (r *QuizzesPostgresRepo) GetMaxOrderQuestionFromQuiz(ctx context.Context, quizId uuid.UUID) (int, error){
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT COALESCE(max(ordering), 0)
		FROM quiz_questions WHERE quiz_id = $1 AND active=TRUE`,
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

func (r *QuizzesPostgresRepo) DeleteCategoryById(ctx context.Context, id uuid.UUID) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx, 
			`DELETE FROM quiz_categories WHERE id = $1`,
			id,
		)
	})
}

func (r *QuizzesPostgresRepo) SoftDeleteCategoryById(ctx context.Context, id uuid.UUID) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx, 
			`UPDATE quiz_categories SET active = FALSE WHERE id = $1`,
			id,
		)
	})
}

func (r *QuizzesPostgresRepo) DeleteQuestions(ctx context.Context, filter quizzes.QuestionFilter) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		query, args := infraestructure.GetFilteredQuery(
			`DELETE FROM quiz_questions WHERE active=TRUE `,
			questionFilter(&filter),
		)
		log.Print(query)
		return ex.ExecContext(ctx, query,args...)
	})
}


func (r *QuizzesPostgresRepo) SoftDeleteQuestions(ctx context.Context, filter quizzes.QuestionFilter) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		query, args := infraestructure.GetFilteredQuery(
			`UPDATE quiz_questions SET active = FALSE WHERE 1=1 `,
			questionFilter(&filter),
		)
		return ex.ExecContext(ctx, query,args...)
	})
}

func (r *QuizzesPostgresRepo) DeleteQuizById(ctx context.Context, quizId uuid.UUID)(int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx,
			`DELETE FROM quizzes WHERE id=$1`,
			quizId,
		)
	})
}


func (r *QuizzesPostgresRepo) SoftDeleteQuizById(ctx context.Context, quizId uuid.UUID) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx,
			`UPDATE quizzes SET active = FALSE WHERE id=$1`,
			quizId,
		)
	})
}


func (r *QuizzesPostgresRepo) DeleteQuizzesPlayed(ctx context.Context, filter quizzes.QuizPlayedFilter) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		query, args := infraestructure.GetFilteredQuery(
			`DELETE FROM quizzes_played WHERE 1=1 `,
			quizzesPlayedFilter(&filter),
		)
		return ex.ExecContext(ctx, query, args...)
	})
}

func (r *QuizzesPostgresRepo) DeleteUsersAnswers(ctx context.Context, filter quizzes.UserAnswerFilter)(int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		query, args := infraestructure.GetFilteredQuery(
			`DELETE FROM user_answers WHERE 1=1 `,
			userAnswerFilter(&filter),
		)
		return ex.ExecContext(ctx, query, args...)
	})
}

func (r *QuizzesPostgresRepo) GetQuestions(ctx context.Context, filter quizzes.QuestionFilter) ([]quizzes.QuestionPlainModel, error){
	query, args := infraestructure.GetFilteredQuery(
		`SELECT id, ordering, question, question_type, options_json, active, quiz_id, strategic_answer_id
		FROM quiz_questions WHERE active=TRUE `,
		questionFilter(&filter),
	)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil{
		return nil, err 
	}
	var questions[]quizzes.QuestionPlainModel
	for rows.Next(){
		model, err := r.rowToQuestion(rows)
		if err != nil{
			return nil, err 
		}
		questions = append(questions, *model)
	}
	return questions, nil
}


func (r *QuizzesPostgresRepo) GetQuizzes(ctx context.Context, filter quizzes.QuizFilter) ([]quizzes.QuizPlainModel, error){
	query, args := infraestructure.GetFilteredQuery(
		`SELECT id, name, description, language_code, image_id, published, active, created_at, category_id, creator_id
		FROM quizzes WHERE active=TRUE `,
		quizFilter(&filter),
	)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil{
		return nil, err 
	}
	var quizes []quizzes.QuizPlainModel
	for rows.Next(){
		model, err := r.rowToQuiz(rows)
		if err != nil{
			return nil, err
		}
		quizes = append(quizes, *model)
	}
	return quizes, nil
}

func (r *QuizzesPostgresRepo) GetQuizzesPlayedCount(ctx context.Context, filter quizzes.QuizPlayedFilter) (int, error){
	query, args := infraestructure.GetFilteredQuery(
		`SELECT count(id) FROM quizzes_played WHERE 1=1 `,
		quizzesPlayedFilter(&filter),
	)
	row := r.db.QueryRowContext(ctx, query, args...)
	var count int 
	if err := row.Scan(&count); err != nil{
		return 0, err 
	}
	return count, nil
}

func (r *QuizzesPostgresRepo) GetUsersAnswersCount(ctx context.Context, filter quizzes.UserAnswerFilter) (int, error){
	query, args := infraestructure.GetFilteredQuery(
		`SELECT count(id) FROM user_answers WHERE 1=1 `,
		userAnswerFilter(&filter),
	)
	row := r.db.QueryRowContext(ctx, query, args...)
	var count int 
	if err := row.Scan(&count); err != nil{
		return 0, err 
	}
	return count, nil
}

func (r *QuizzesPostgresRepo) GetQuestionById(ctx context.Context, questionId uuid.UUID) (*quizzes.QuestionPlainModel, error){
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT id, ordering, question, question_type, options_json, active, quiz_id, strategic_answer_id
		FROM quiz_questions WHERE active=TRUE AND id=$1`,
		questionId,
	)
	return r.rowToQuestion(row)
}

func (r *QuizzesPostgresRepo) UpdateQuestion(ctx context.Context, model *quizzes.QuestionPlainModel) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(ctx, 
			`UPDATE quiz_questions 
			SET question = $1, options_json = $2, strategic_answer_id = $3
			WHERE id=$4`,
			model.Question, model.OptionsJson, model.StrategicAnswerId, model.Id,
		)
	})
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

func (r *QuizzesPostgresRepo) rowToCategory(row infraestructure.Scanable) (*quizzes.QuizCatPlainModel, error){
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

func (r *QuizzesPostgresRepo) rowToQuestion(row infraestructure.Scanable) (*quizzes.QuestionPlainModel, error){
	model := new(quizzes.QuestionPlainModel)
	err := row.Scan(&model.Id, &model.Ordering, &model.Question, &model.QuestionType, &model.OptionsJson, &model.Active, &model.QuizId, &model.StrategicAnswerId)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, nil 
		}
		return nil, err 
	}
	return model, nil
}

func (r *QuizzesPostgresRepo) rowToQuiz(row infraestructure.Scanable) (*quizzes.QuizPlainModel, error){
	model := new(quizzes.QuizPlainModel)
	err := row.Scan(&model.Id, &model.Name, &model.Description, &model.LanguageCode, &model.ImageId, 
	&model.Published, &model.Active, &model.CreatedAt, &model.CategoryId, &model.CreatorId)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, nil 
		}
		return nil, err 
	}
	return model, nil
}