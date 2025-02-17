package quizzes

import (
	"time"

	"github.com/google/uuid"
)

type QuizCatPlainModel struct {
	Id 			uuid.UUID
	Name 		string 
	Description	string 
	CreatedAt 	time.Time
	Active 		bool
	ImageId		uuid.UUID
}


type QuizCatModel struct {
	Id 			uuid.UUID	`json:"id"`
	Name 		string 		`json:"name"`
	Description	string 		`json:"description"`
	CreatedAt 	time.Time	`json:"-"`
	ImageUrl	string 		`json:"imageUrl"`
}


type QuizPlainModel struct{
	Id 				uuid.UUID	
	Name 			string 		
	Description 	string
	LanguageCode 	string
	ImageId 		*uuid.UUID
	Published 		bool 
	Active 			bool 
	CreatedAt 		time.Time
	CategoryId 		*uuid.UUID
	CreatorId 		*uuid.UUID
}

type QuizModel struct{
	Id 			uuid.UUID	`json:"id"`
	Name 		string 		`json:"name"`
	Description string 		`json:"description"`
	ImageUrl 	string 		`json:"imageUrl"`
	Category	*QuizCatModel	`json:"category"`
}

type QuestionPlainModel struct{
	Id 				uuid.UUID
	Ordering 		int 
	Question 		string 
	QuestionType 	string 
	OptionsJson  	string 
	QuizId 			uuid.UUID
	Active 			bool
	StrategicAnswerId 	*uuid.UUID
}


type StrategicAnswerModel struct{
	Id 				uuid.UUID 
	Name 			string 
	Description 	string
}

type UserAnswerPlainModel struct{
	Id 			uuid.UUID
	UserId 		uuid.UUID
	QuestionId 	uuid.UUID
	Answers 	string 
	AnsweredAt 	time.Time
}

type QuizPlayedPlainModel struct{
	Id 		uuid.UUID
	QuizId 	uuid.UUID
	UserId 	uuid.UUID
	Shared 	bool 
	Score 	*int
	StartedAt 	time.Time
	CompletedAt *time.Time 
}