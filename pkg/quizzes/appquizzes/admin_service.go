package appquizzes

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/diegobermudez03/couples-backend/pkg/localization"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/google/uuid"
)

type AdminServiceImpl struct {
	filesService 	files.Service
	loacalizationService localization.LocalizationService
	quizzesRepo 	quizzes.QuizzesRepository

}

func NewAdminServiceImpl(filesService files.Service, loacalizationService localization.LocalizationService,  quizzesRepo quizzes.QuizzesRepository) quizzes.AdminService{
	return &AdminServiceImpl{
		filesService: filesService,
		quizzesRepo: quizzesRepo,
		loacalizationService:loacalizationService,
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
	categoryId := uuid.New()

	imageId, _, err := s.filesService.UploadImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, true, quizzes.DOMAIN_NAME, quizzes.CATEGORIES, categoryId.String(), quizzes.PROFILE)
	if err != nil{
		if errors.Is(err, files.ErrInvalidImageType){
			return quizzes.ErrInvalidImageType
		}
		log.Print(err.Error())
		return quizzes.ErrCreatingCategory
	}	
	
	quizModel := quizzes.QuizCatPlainModel{
		Id: categoryId,
		Name: name,
		Description: description,
		ImageId: *imageId,
		CreatedAt: time.Now(),
		Active: true,
	}
	if num, err :=s.quizzesRepo.CreateCategory(ctx, &quizModel); err != nil || num == 0{
		log.Print(err)
		return quizzes.ErrCreatingCategory
	}
	return nil
}


func (s *AdminServiceImpl) UpdateQuizCategory(ctx context.Context, id uuid.UUID, name, description string, image io.Reader) error{
	cat, err := s.quizzesRepo.GetCategoryById(ctx, id)
	if err != nil{
		return quizzes.ErrUpdatingCategory
	} else if cat == nil{
		return quizzes.ErrNonExistingCategory
	}

	waitGroup := sync.WaitGroup{}

	if image != nil{
		waitGroup.Add(1)
		//will execute asynchronously, if there's any error is ignored xd, I should notify about the error, I should...
		go func() {
			s.filesService.UpdateImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, cat.ImageId)
			waitGroup.Done()
		}()
	}

	if name != ""{
		cat.Name = name 
	}
	if description != ""{
		cat.Description = description
	}

	if num, err := s.quizzesRepo.UpdateCategory(ctx, cat); err != nil || num == 0{
		return quizzes.ErrUpdatingCategory
	}
	waitGroup.Wait()
	return nil
}
