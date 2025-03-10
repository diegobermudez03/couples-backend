package auth

import (
	"time"

	"github.com/google/uuid"
)

type UserAuthModel struct {
	Id    			uuid.UUID
	Email 			*string
	Hash  			*string
	OauthProvider 	*string 
	OauthId			*string
	CreatedAt		time.Time
	UserId 			*uuid.UUID
}

type SessionModel struct {
	Id			uuid.UUID
	Token 		string 
	Device 		*string 
	Os 			*string 
	ExpiresAt 	time.Time
	CreatedAt 	time.Time
	LastUsed 	time.Time 
	UserAuthId 	uuid.UUID
}

type AdminSessionModel struct{
	Id			uuid.UUID
	Token 		string 
	CreatedAt 	time.Time
}

type TempCoupleModel struct{
	Code 		int 		`json:"code"`
	StartDate 	time.Time 	`json:"startDate"`
}