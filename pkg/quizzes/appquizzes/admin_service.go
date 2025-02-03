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

	imageId, err := s.filesService.UploadImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, true, quizzes.DOMAIN_NAME, quizzes.CATEGORIES, categoryId.String(), quizzes.PROFILE)
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


func (s *AdminServiceImpl) CreateQuiz(ctx context.Context, name, description, languageCode string, categoryId uuid.UUID, image io.Reader) error{
	if name == ""{
		return quizzes.ErrEmptyQuizName
	}
	if err := s.loacalizationService.ValidateLanguage(languageCode); err != nil{
		return quizzes.ErrInvalidLanguage
	}

	cat, err := s.quizzesRepo.GetCategoryById(ctx, categoryId)
	if err != nil{
		log.Print(err.Error())
		return quizzes.ErrCreatingQuiz
	} else if cat == nil{
		return quizzes.ErrCategoryDontExists
	}

	quizId := uuid.New()
	var imageId *uuid.UUID
	if image != nil{
		imageId, _ = s.filesService.UploadImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, true, quizzes.DOMAIN_NAME, quizzes.QUIZZES, quizId.String(), quizzes.PROFILE)
	}

	model := quizzes.QuizPlainModel{
		Id: quizId,
		Name: name,
		Description: description,
		LanguageCode: languageCode,
		ImageId: imageId,
		Published: false,
		Active: true,
		CreatedAt: time.Now(),
		CategoryId: categoryId,
		CreatorId: nil,
	}

	if num, err := s.quizzesRepo.CreateQuiz(ctx, &model); err != nil || num == 0{
		log.Print(err.Error())
		return quizzes.ErrCreatingQuiz
	}
	return nil
}


func (s *AdminServiceImpl)  UpdateQuiz(ctx context.Context, quizId uuid.UUID, name, description string, categoryId *uuid.UUID, image io.Reader) error{
	quiz, err := s.quizzesRepo.GetQuizById(ctx, quizId)
	if err != nil{
		return quizzes.ErrUpdatingQuiz
	}else if quiz == nil{
		return quizzes.ErrQuizNotFound
	}

	//if we are updating the image
	if image != nil{
		//if there was already an image
		if quiz.ImageId != nil{
			s.filesService.UpdateImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, *quiz.ImageId)
		} else{
			//if its a new image
			id, err := s.filesService.UploadImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, true, quizzes.DOMAIN_NAME, quizzes.QUIZZES, quiz.Id.String(), quizzes.PROFILE)
			if err == nil{
				quiz.ImageId = id
			}
		}
	}

	if name != ""{
		quiz.Name = name
	}
	if description != ""{
		quiz.Description = description
	}
	if categoryId != nil{
		if cat, err := s.quizzesRepo.GetCategoryById(ctx, *categoryId); err == nil && cat != nil{
			quiz.CategoryId = *categoryId
		} 
	}

	if num, err := s.quizzesRepo.UpdateQuiz(ctx, quiz); num == 0 || err != nil{
		return quizzes.ErrUpdatingQuiz
	}
	return nil 
}