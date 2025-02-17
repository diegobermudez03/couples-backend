package quizzes

import "github.com/google/uuid"

type QuestionFilter struct {
	Id     *uuid.UUID
	QuizId *uuid.UUID
}

type QuizPlayedFilter struct {
	Id     *uuid.UUID
	QuizId *uuid.UUID
}

type UserAnswerFilter struct {
	Id         *uuid.UUID
	QuestionId *uuid.UUID
}

type QuizFilter struct{
	Id 		*uuid.UUID
	CategoryId 	*uuid.UUID
	Text 		*string
	OrderBy 	*string 
	PlayerId 	*uuid.UUID
	CreatorId 	*uuid.UUID
	LanguageCode *string
	FetchFilters
}
