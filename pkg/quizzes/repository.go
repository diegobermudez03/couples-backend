package quizzes

import (
	"context"

	"github.com/google/uuid"
)

type QuizzesRepository interface {
	GetCategoryByName(ctx context.Context, name string) (*QuizCatPlainModel, error)
	GetCategoryById(ctx context.Context, id uuid.UUID) (*QuizCatPlainModel, error)
	CreateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)
	UpdateCategory(ctx context.Context, category *QuizCatPlainModel) (int, error)
	DeleteCategoryById(ctx context.Context, id uuid.UUID) (int, error)
	SoftDeleteCategoryById(ctx context.Context, id uuid.UUID) (int, error)

	GetQuizzes(ctx context.Context, filter QuizFilter) ([]QuizPlainModel, error)
	GetQuizById(ctx context.Context, id uuid.UUID) (*QuizPlainModel, error)
	CreateQuiz(ctx context.Context, quiz *QuizPlainModel) (int, error)
	UpdateQuiz(ctx context.Context, quiz *QuizPlainModel) (int, error)
	SoftDeleteQuizById(ctx context.Context, quizId uuid.UUID) (int, error)
	DeleteQuizById(ctx context.Context, quizId uuid.UUID)(int, error)

	GetMaxOrderQuestionFromQuiz(ctx context.Context, quizId uuid.UUID) (int, error)
	GetQuestions(ctx context.Context, filter QuestionFilter) ([]QuestionPlainModel, error)
	GetQuestionById(ctx context.Context, questionId uuid.UUID) (*QuestionPlainModel, error)
	CreateQuestion(ctx context.Context, model *QuestionPlainModel) (int, error)
	UpdateQuestion(ctx context.Context, model *QuestionPlainModel) (int, error)
	DeleteQuestions(ctx context.Context, filter QuestionFilter) (int, error)
	SoftDeleteQuestions(ctx context.Context, filter QuestionFilter) (int, error)

	GetQuizzesPlayedCount(ctx context.Context, filter QuizPlayedFilter) (int, error)
	DeleteQuizzesPlayed(ctx context.Context, filter QuizPlayedFilter) (int, error)

	GetStrategicTypeAnswerById(ctx context.Context, id uuid.UUID) (*StrategicAnswerModel, error)
	CreateStrategicTypeAnswer(ctx context.Context, model *StrategicAnswerModel) (int, error)

	GetUsersAnswersCount(ctx context.Context, filter UserAnswerFilter) (int, error)
	DeleteUsersAnswers(ctx context.Context, filter UserAnswerFilter)(int, error)
}
