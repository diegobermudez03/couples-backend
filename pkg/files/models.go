package files

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id        uuid.UUID
	Bucket    string
	ObjectKey string
	Type      string
	CreatedAt time.Time
}