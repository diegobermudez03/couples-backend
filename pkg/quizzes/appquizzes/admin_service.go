package appquizzes

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/google/uuid"
)

type AdminServiceImpl struct {
	filesService 	files.Service
	quizzesRepo 	quizzes.QuizzesRepository

}

func NewAdminServiceImpl(filesService 	files.Service, quizzesRepo quizzes.QuizzesRepository) quizzes.AdminService{
	return &AdminServiceImpl{
		filesService: filesService,
		quizzesRepo: quizzesRepo,
	}
}


func (s *AdminServiceImpl) CreateQuizCategory(ctx context.Context, name, description string, image io.Reader) error{
	if name == "" || description == "" || image == nil{
		return quizzes.ErrMissingCategoryAttributes
	}

	cat, err := s.quizzesRepo.GetCategoryByName(ctx, strings.ToLower(name))
	if err == nil && cat != nil{
		return quizzes.ErrCategoryAlreadyExists
	}

	imageId, err := s.filesService.UploadImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, quizzes.DOMAIN_NAME,quizzes.CATEGORIES, name)
	if err != nil{
		if errors.Is(err, files.ErrInvalidImageType){
			return quizzes.ErrInvalidImageType
		}
		log.Print(err.Error())
		return quizzes.ErrCreatingCategory
	}	
	
	quizModel := quizzes.QuizCatPlainModel{
		Id: uuid.New(),
		Name: name,
		Description: description,
		ImageId : imageId,
		CreatedAt: time.Now(),
	}
	if num, err :=s.quizzesRepo.CreateCategory(ctx, &quizModel); err != nil || num == 0{
		log.Print(err)
		return quizzes.ErrCreatingCategory
	}
	return nil
}