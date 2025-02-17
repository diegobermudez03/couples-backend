package repoquizzes

import "github.com/diegobermudez03/couples-backend/pkg/quizzes"

func questionFilter(filter *quizzes.QuestionFilter) map[string]any{
	return map[string]any{
		"id" : filter.Id,
		"quiz_id" : filter.QuizId,
	}
}

func quizzesPlayedFilter(filter *quizzes.QuizPlayedFilter) map[string]any{
	return map[string]any{
		"id" : filter.Id,
		"quiz_id" : filter.QuizId,
	}
}

func userAnswerFilter(filter *quizzes.UserAnswerFilter) map[string]any{
	return map[string]any{
		"id" : filter.Id,
		"question_id" : filter.QuestionId,
	}
}



func quizFilter(filter *quizzes.QuizFilter) map[string]any{
	return map[string]any{
		"id" 	:	filter.Id,
		"category_id" : filter.CategoryId,
		"language_code" : filter.LanguageCode,
	}
}