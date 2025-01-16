package domainauth

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