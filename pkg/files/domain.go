package files

import (
	"context"
	"io"
	"os"

	"github.com/google/uuid"
)

type Service interface {
	UploadImage(ctx context.Context, image io.Reader, maxSize int64, public bool, path ...string) (*uuid.UUID,error)
	UpdateImage(ctx context.Context, image io.Reader, maxSize int64, id uuid.UUID) (error)
	GetImage(ctx context.Context, path string) (*os.File, string, error)
	//DeleteImage(ctx context.Context, imageId uuid.UUID) error
}

type Repository interface{
	CreateFile(ctx context.Context, file *FileModel) (int,error)
	GetFileById(ctx context.Context, id uuid.UUID) (*FileModel, error)
}

type FileRepository interface{
	StoreFile(ctx  context.Context, bucket, group, objectKey string, image io.Reader) error
	GetFile(ctx context.Context, path string) (*os.File, error)
}


const JPG_TYPE = "image/jpeg"
const PNG_TYPE = "image/png"


const MAX_SIZE_PROFILE_PICTURE = 2073600 		//	1920x1080
const MAX_SIZE_QUESTION_PICTURE = 160000 // 400x400