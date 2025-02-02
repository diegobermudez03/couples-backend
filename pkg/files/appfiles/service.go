package appfiles

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"path/filepath"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/google/uuid"
	"golang.org/x/image/draw"
)

type FilesServiceImpl struct{
	filesRepo 		files.FileRepository
}

func NewFilesServiceImpl(filesRepo files.FileRepository) files.Service{
	return &FilesServiceImpl{
		filesRepo: filesRepo,
	}
}


func (s *FilesServiceImpl) UploadImage(ctx context.Context, imageReader io.Reader,  maxSize int64, path ...string) (*uuid.UUID,error){
	if len(path) < 3{
		return nil, files.ErrPathNotLongEnough
	}

	//detect image type
	imageBytes := make([]byte, 512)
	_, err := imageReader.Read(imageBytes)
	if err != nil{
		return nil, files.ErrDetectingImageType
	}
	imageType := http.DetectContentType(imageBytes)
	//go back to the beginning
	_, err = imageReader.(io.Seeker).Seek(0, io.SeekStart)
	if err != nil{
		return nil, files.ErrDetectingImageType
	}

	var imageSrc image.Image
	switch imageType{
	case files.JPG_TYPE:
		imageSrc, err = jpeg.Decode(imageReader)
	case files.PNG_TYPE:
		imageSrc, err = png.Decode(imageReader)
	default:
		return nil, files.ErrInvalidImageType
	}

	if err != nil{
		return nil, files.ErrDetectingImageType 
	}

	width := imageSrc.Bounds().Dx()
	height := imageSrc.Bounds().Dy()
	nPixels := width*height

	//if image is bigger than maximum
	var finalImage image.Image
	if nPixels > int(maxSize){
		//the factor is the value we need to multiply in order to meet the maximum size
		factor := math.Sqrt(float64(maxSize)/float64(nPixels))
		width = int(float64(width)*factor)
		height = int(float64(height)*factor)
		imageDst := image.NewRGBA(image.Rect(0,0, width, height))
		draw.CatmullRom.Scale(imageDst, imageDst.Rect, imageSrc, imageSrc.Bounds(), draw.Over, nil)
		finalImage = imageDst
	}else{
		finalImage = imageSrc
	}

	// store image 
	bucket := path[0]
	object := path[len(path)-1]
	path = path[1:len(path)-1]
	group := filepath.Join(path...)
	if err := s.filesRepo.StoreFile(ctx, bucket, group, object, finalImage); err != nil{
		return nil, err
	}

	// add to database
	
	return nil, nil
}
