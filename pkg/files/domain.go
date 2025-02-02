package files

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type Service interface {
	UploadImage(ctx context.Context, image io.Reader, maxSize int64, path ...string) (*uuid.UUID,error)
}

type Repository interface{
	CreateFile(ctx context.Context, file *FileModel) (int,error)
}

type FileRepository interface{
	StoreFile(ctx  context.Context, bucket, group, objectKey string, image io.Reader) error
}


const JPG_TYPE = "image/jpeg"
const PNG_TYPE = "image/png"


const MAX_SIZE_PROFILE_PICTURE = 2073600 		//	1920x1080