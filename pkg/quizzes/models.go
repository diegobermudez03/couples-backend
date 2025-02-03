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
	Id 			uuid.UUID
	Name 		string 
	Description	string 
	CreatedAt 	time.Time
	ImageUrl	string 
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
	CategoryId 		uuid.UUID
	CreatorId 		*uuid.UUID
}