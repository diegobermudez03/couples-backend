package repofiles

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"github.com/diegobermudez03/couples-backend/pkg/files"
)

const ROOT_FOLDER = "files"


type LocalStorage struct{}

func NewLocalStorage() files.FileRepository{
	return &LocalStorage{}
}

func (r *LocalStorage) StoreFile(ctx  context.Context, bucket, group, objectKey string, image image.Image) error{
	path := r.getPath(bucket, group)
	log.Print(path)
	if err := os.MkdirAll(path, os.ModePerm); err != nil{
		return err
	}
	outWriter, err := os.Create(fmt.Sprint(filepath.Join(path, objectKey), ".jpg"))
	if err != nil{
		return err 
	}
	defer outWriter.Close()
	if err := jpeg.Encode(outWriter, image, &jpeg.Options{Quality: 80}); err != nil{
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

