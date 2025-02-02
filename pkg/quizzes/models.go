package quizzes

import (
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/google/uuid"
)

type QuizCatPlainModel struct {
	Id 			uuid.UUID
	Name 		string 
	Description	string 
	CreatedAt 	time.Time
	Active 		bool
	File		*files.FileModel
}


type QuizCatModel struct {
	Id 			uuid.UUID
	Name 		string 
	Description	string 
	CreatedAt 	time.Time
	ImageUrl	string 
}
