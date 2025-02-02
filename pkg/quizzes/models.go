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
	FileId		*uuid.UUID
}


type QuizCatModel struct {
	Id 			uuid.UUID
	Name 		string 
	Description	string 
	CreatedAt 	time.Time
	ImageUrl	string 
}
