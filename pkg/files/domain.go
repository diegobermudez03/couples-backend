package files

import (
	"context"
	"io"
	"os"

	"github.com/google/uuid"
)

type Service interface {
	UploadImage(ctx context.Context, image io.Reader, maxSize int64, public bool, path ...string) (imId *uuid.UUID, url *string, err error)
	UpdateImage(ctx context.Context, image io.Reader, maxSize int64, id uuid.UUID) (error)
	GetImage(ctx context.Context, path string) (*os.File, string, error)
	DeleteImage(ctx context.Context, imageId uuid.UUID) error
	GetBatchUrls(ctx context.Context, imagesIds []uuid.UUID) (map[uuid.UUID]string, error)
}

type Repository interface{
	CreateFile(ctx context.Context, file *FileModel) (int,error)
	GetFileById(ctx context.Context, id uuid.UUID) (*FileModel, error)
	DeleteFileById(ctx context.Context, id uuid.UUID) (int, error)
	GetBatchUrls(ctx context.Context, imageIds []uuid.UUID)(map[uuid.UUID]string, error)
}

type FileRepository interface{
	StoreFile(ctx  context.Context, bucket, group, objectKey string, image io.Reader) error
	GetFile(ctx context.Context, path string) (*os.File, error)
	DeleteFile(ctx context.Context, bucket, group, objectKey string) error
}


const JPG_TYPE = "image/jpeg"
const PNG_TYPE = "image/png"


const MAX_SIZE_PROFILE_PICTURE = 2073600 		//	1920x1080
const MAX_SIZE_QUESTION_PICTURE = 160000 // 400x400