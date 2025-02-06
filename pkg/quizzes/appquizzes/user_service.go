package appquizzes

import (
	"context"
	"io"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/google/uuid"
)

type QuestionOptionsCreator func(ctx context.Context, quiz *quizzes.QuizPlainModel, inputOptions string, images map[string]io.Reader) (string, error)

type UserService struct{
	fileService	files.Service
	repo 	quizzes.QuizzesRepository
	creators map[string]QuestionOptionsCreator
}

func NewUserService(fileService	files.Service, repo quizzes.QuizzesRepository) quizzes.UserService{
	service := &UserService{
		fileService: fileService,
		repo: repo,
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
	optionsJson, err := creator(ctx, quiz, parameters.OptionsJson, images)
	if err != nil{
		return err
	}

	// create model to store
	questionModel := quizzes.QuestionPlainModel{
		Id: uuid.New(),
		OptionsJson: optionsJson,
		QuestionType: parameters.QType,
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
		return quizzes.ErrCreatingQuestion
	}
	questionModel.Ordering = maxOrder + 1 

	//write question
	if err := s.repo.CreateQuestion(ctx, &questionModel); err != nil{
		return quizzes.ErrCreatingQuestion
	}
	return nil
}




//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
///				PRIVATE METHODS				/////


func (s *UserService) openCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader) (string, error){
	return "", nil
}

func (s *UserService) multipleChCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader) (string, error){
	return "", nil
}

func (s *UserService) matchingCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader) (string, error){
	return "", nil
}

func (s *UserService) dragAndDropCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader) (string, error){
	return "", nil
}