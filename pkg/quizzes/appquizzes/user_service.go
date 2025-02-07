package appquizzes

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type QuestionOptionsCreator func(ctx context.Context, quiz *quizzes.QuizPlainModel, inputOptions string, images map[string]io.Reader, questionId uuid.UUID) (string, error)

type UserService struct{
	fileService		files.Service
	repo 			quizzes.QuizzesRepository
	creators 		map[string]QuestionOptionsCreator
	jsonValidator 	*validator.Validate
}

func NewUserService(fileService	files.Service, repo quizzes.QuizzesRepository) quizzes.UserService{
	service := &UserService{
		fileService: fileService,
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


func (s *UserService) CreateQuestion(ctx context.Context, userId *uuid.UUID, quizId uuid.UUID, parameters quizzes.CreateQuestionRequest, images map[string]io.Reader) error{
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