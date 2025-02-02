package repofiles

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/diegobermudez03/couples-backend/pkg/files"
)

const ROOT_FOLDER = "files"


type LocalStorage struct{}

func NewLocalStorage() files.FileRepository{
	return &LocalStorage{}
}

func (r *LocalStorage) StoreFile(ctx  context.Context, bucket, group, objectKey string, image io.Reader) error{
	path := r.getPath(bucket, group)
	if err := os.MkdirAll(path, os.ModePerm); err != nil{
		return err
	}
	file, err := os.Create(filepath.Join(path, objectKey))
	if err != nil{
		return err 
	}
	defer file.Close()

	if _, err := io.Copy(file, image); err != nil{
		return err
	}
	return nil
}


////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////
///////////////////			PRIVATE METHODS 						////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

func (s *LocalStorage) getPath(bucket, object string) string{
	return filepath.Join("..", ROOT_FOLDER, bucket, object)
}

