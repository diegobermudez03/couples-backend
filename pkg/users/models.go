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

type CoupleModel struct{
	Id 		uuid.UUID
	HeId 	uuid.UUID
	SheId 	uuid.UUID
	RelationStart 	time.Time 
	EndDate 		*time.Time
}


type PointsModel struct{
	Id 			uuid.UUID
	Day 		time.Time
	Points 		int 
	UserId 		*uuid.UUID
}