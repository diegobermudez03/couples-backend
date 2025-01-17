package users

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	Id uuid.UUID
	FirstName 	string 
	LastName 	string 
	NickName	string
	Gender 		string 
	BirthDate 	time.Time
	CreatedAt 	time.Time
	Active 		bool 
	CountryCode string 
	LanguageCode string  
}

type TempCoupleModel struct{
	UserId 		uuid.UUID
	Code 		int 
	StartDate 	time.Time 
	CreatedAt 	time.Time 
	UpdatedAt 	time.Time
}