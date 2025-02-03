package files

import (
	"github.com/google/uuid"
)

type FileModel struct {
	Id        uuid.UUID
	Bucket    string
	Group 	  string
	ObjectKey string
	Url 		*string
	Public 		bool
	Type      string
}