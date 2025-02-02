package files

import (
	"time"

	"github.com/google/uuid"
)

type FileModel struct {
	Id        uuid.UUID
	Bucket    string
	Group 	  string
	ObjectKey string
	Type      string
	CreatedAt time.Time
}