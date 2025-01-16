package users

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	Id uuid.UUID
	FirstName 	string 
	LastName 	string 
	Gender 		string 
	BirthDate 	time.Time
	CreatedAt 	time.Time
	Active 		bool 
	CountryCode string 
	LanguageCode string  
}