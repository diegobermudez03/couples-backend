package appquizzes

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/diegobermudez03/couples-backend/pkg/localization"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type QuestionOptionsCreator func(ctx context.Context, quiz *quizzes.QuizPlainModel, inputOptions string, images map[string]io.Reader, questionId uuid.UUID) (string, error)

type UserService struct{
	fileService		files.Service
	loacalizationService localization.LocalizationService
	repo 			quizzes.QuizzesRepository
	creators 		map[string]QuestionOptionsCreator
	jsonValidator 	*validator.Validate
}

func NewUserService(
	fileService	files.Service, 
	loacalizationService localization.LocalizationService, 
	repo quizzes.QuizzesRepository) quizzes.UserService{
	service := &UserService{
		fileService: fileService,
		loacalizationService: loacalizationService,
		repo: repo,
		jsonValidator: validator.New(),
	}
	service.creators = map[string]QuestionOptionsCreator{
		quizzes.TRUE_FALSE_TYPE : service.trueFalseCreator,
		quizzes.SLIDER_TYPE : service.sliderCreator,
		quizzes.ORDERING_TYPE : service.orderingCreator,
		quizzes.OPEN_TYPE : service.openCreator,
		quizzes.MULTIPLE_CH_TYPE : service.multipleChCreator,
		quizzes.MATCHING_TYPE : service.matchingCreator,
		quizzes.DRAG_AND_DROP_TYPE : service.dragAndDropCreator,
	}
	return service
}

func (s *UserService) GetQuizById(ctx context.Context, quizId uuid.UUID)(*quizzes.QuizPlainModel, error){
	quiz, err := s.repo.GetQuizById(ctx, quizId)
	if err != nil{
		return nil, quizzes.ErrRetrievingQuiz
	}else if quiz == nil{
		return nil, quizzes.ErrQuizNotFound
	}
	return quiz, nil
}

func (s *UserService) CreateQuiz(ctx context.Context, name, description, languageCode string, categoryId, userId *uuid.UUID, image io.Reader) error{
	if name == ""{
		return quizzes.ErrEmptyQuizName
	}
	if err := s.loacalizationService.ValidateLanguage(languageCode); err != nil{
		return quizzes.ErrInvalidLanguage
	}
	if categoryId != nil{
		cat, err := s.repo.GetCategoryById(ctx, *categoryId)
		if err != nil{
			return quizzes.ErrCreatingQuiz
		} else if cat == nil{
			categoryId = nil
		}
	}

	quizId := uuid.New()
	var imageId *uuid.UUID
	if image != nil{
		imageId, _, _ = s.fileService.UploadImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, true, quizzes.DOMAIN_NAME, quizzes.QUIZZES, quizId.String(), quizzes.PROFILE)
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
		CreatorId: userId,
	}

	if num, err := s.repo.CreateQuiz(ctx, &model); err != nil || num == 0{
		log.Print(err.Error())
		return quizzes.ErrCreatingQuiz
	}
	return nil
}


func (s *UserService)  UpdateQuiz(ctx context.Context, quizId uuid.UUID, name, description, languageCode string, categoryId *uuid.UUID, image io.Reader) error{
	quiz, err := s.repo.GetQuizById(ctx, quizId)
	if err != nil{
		return quizzes.ErrUpdatingQuiz
	}else if quiz == nil{
		return quizzes.ErrQuizNotFound
	}

	//if we are updating the image
	if image != nil{
		//if there was already an image
		if quiz.ImageId != nil{
			s.fileService.UpdateImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, *quiz.ImageId)
		} else{
			//if its a new image
			id, _, err := s.fileService.UploadImage(ctx, image, files.MAX_SIZE_PROFILE_PICTURE, true, quizzes.DOMAIN_NAME, quizzes.QUIZZES, quiz.Id.String(), quizzes.PROFILE)
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
	if languageCode != ""{
		quiz.LanguageCode = languageCode
	}
	if categoryId != nil{
		if cat, err := s.repo.GetCategoryById(ctx, *categoryId); err == nil && cat != nil{
			quiz.CategoryId = categoryId
		} 
	}

	if num, err := s.repo.UpdateQuiz(ctx, quiz); num == 0 || err != nil{
		return quizzes.ErrUpdatingQuiz
	}
	return nil 
}

func (s *UserService) CreateQuestion(ctx context.Context, quizId uuid.UUID, parameters quizzes.CreateQuestionRequest, images map[string]io.Reader) error{
	quiz, err := s.repo.GetQuizById(ctx, quizId)
	if err != nil || quiz  == nil{
		return quizzes.ErrQuizNotFound
	}

	//call specific question type creator for options JSON
	creator, ok := s.creators[parameters.QType]
	if !ok{
		return quizzes.ErrInvalidQuestionType
	}
	questionId := uuid.New()
	inputOptionsJson, err := json.Marshal(parameters.OptionsJson)
	if err != nil{
		return quizzes.ErrInvalidQuestionOptions
	}
	optionsJson, err := creator(ctx, quiz, string(inputOptionsJson), images, questionId)
	if err != nil{
		return err
	}

	// create model to store
	questionModel := quizzes.QuestionPlainModel{
		Id:questionId,
		OptionsJson: optionsJson,
		Question: parameters.Question,
		QuestionType: parameters.QType,
		QuizId: quizId,
		Active: true,
	}

	//create strategic question if needed
	if parameters.StrategicAnswerId != nil{
		st, err := s.repo.GetStrategicTypeAnswerById(ctx, *parameters.StrategicAnswerId)
		if st != nil && err == nil{
			questionModel.StrategicAnswerId = &st.Id
		}
	}else if parameters.StrategicName != nil && parameters.StrategicDescription != nil{
		stId := uuid.New()
		stModel := quizzes.StrategicAnswerModel{
			Id: stId,
			Name: *parameters.StrategicName,
			Description: *parameters.StrategicDescription,
		}
		if num, err := s.repo.CreateStrategicTypeAnswer(ctx, &stModel); err == nil && num != 0{
			questionModel.StrategicAnswerId = &stId
		}
	}
	//calculate ordering 
	maxOrder, err := s.repo.GetMaxOrderQuestionFromQuiz(ctx, quizId)
	if err != nil{
		log.Print(err.Error())
		return quizzes.ErrCreatingQuestion
	}
	questionModel.Ordering = maxOrder + 1 

	//write question
	if num, err := s.repo.CreateQuestion(ctx, &questionModel); err != nil || num == 0{
		log.Print(err.Error())
		return quizzes.ErrCreatingQuestion
	}
	return nil
}

func (s *UserService) DeleteQuestion(ctx context.Context, questionId uuid.UUID) error{
	count, err := s.repo.GetUsersAnswersCountByQuestionId(ctx, questionId)
	if err != nil{
		return quizzes.ErrDeletingQuestion
	}

	var num int 
	if count == 0{
		num, err = s.repo.DeleteQuestions(ctx, quizzes.QuestionFilter{Id: &questionId})
	}else{
		num, err = s.repo.SoftDeleteQuestions(ctx, quizzes.QuestionFilter{Id: &questionId})
	}

	if err != nil || num == 0{ 
		return quizzes.ErrDeletingQuestion
	}
	return nil
}


func (s *UserService) DeleteQuiz(ctx context.Context, quizId uuid.UUID) error{
	count, err := s.repo.GetQuizzesPlayedCount(ctx, quizzes.QuizPlayedFilter{QuizId: &quizId})
	if err != nil{
		return quizzes.ErrDeletingQuiz
	}
	var num int 
	if count == 0{
		questions, auxErr := s.repo.GetQuestionsByQuizId(ctx, quizId)
		if auxErr != nil{
			return quizzes.ErrDeletingQuiz
		}
		for _, q := range questions{
			_, auxErr := s.repo.DeleteUsersAnswers(ctx, quizzes.UserAnswerFilter{QuestionId: &q.Id})
			if auxErr != nil{
				return quizzes.ErrDeletingQuiz
			}
		}
		num, err = s.repo.DeleteQuestions(ctx, quizzes.QuestionFilter{QuizId: &quizId})
	}else{
		num, err = s.repo.SoftDeleteQuizById(ctx, quizId)
	}
	if num == 0 || err != nil{
		return quizzes.ErrDeletingQuiz
	}
	return nil 
}


//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
///				PRIVATE METHODS				/////

func (s *UserService) getOptionImagePath(quizId uuid.UUID, questionId uuid.UUID, imageName string) []string{
	return []string{quizzes.DOMAIN_NAME, quizzes.QUIZZES, quizId.String(), questionId.String(), imageName}
}


func (s *UserService) readJson(jsonText string, payload any) error{
	if err := json.Unmarshal([]byte(jsonText), &payload); err != nil{
		return err
	}
	if err := s.jsonValidator.Struct(payload); err != nil{
		return err 
	}
	return nil
}


