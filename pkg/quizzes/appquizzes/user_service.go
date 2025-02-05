package appquizzes

import (
	"context"
	"io"

	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/google/uuid"
)

type UserService struct{
	repo 	quizzes.QuizzesRepository
}

func NewUserService(repo quizzes.QuizzesRepository) quizzes.UserService{
	return &UserService{
		repo: repo,
	}
}


func (s *UserService) CreateQuestion(ctx context.Context, quizId uuid.UUID, parameters quizzes.CreateQuestionRequest, images map[string]io.Reader) error{
	quiz, err := s.repo.GetQuizById(ctx, quizId)
	if err != nil || quiz  == nil{
		return quizzes.ErrQuizNotFound
	}

	//call specific question type creator for options JSON
	var optionsJson string
	switch parameters.QType{
	case quizzes.TRUE_FALSE_TYPE: optionsJson, err = s.trueFalseCreator(ctx, quiz, parameters.OptionsJson)
	case quizzes.SLIDER_TYPE: optionsJson, err = s.sliderCreator(ctx, quiz, parameters.OptionsJson)
	case quizzes.ORDERING_TYPE: optionsJson, err = s.orderingCreator(ctx, quiz, parameters.OptionsJson)
	case quizzes.OPEN_TYPE: optionsJson, err = s.openCreator(ctx, quiz, parameters.OptionsJson)
	case quizzes.MULTIPLE_CH_TYPE: optionsJson, err = s.multipleChCreator(ctx, quiz, parameters.OptionsJson)
	case quizzes.MATCHING_TYPE: optionsJson, err = s.matchingCreator(ctx, quiz, parameters.OptionsJson)
	case quizzes.DRAG_AND_DROP_TYPE: optionsJson, err = s.dragAndDropCreator(ctx, quiz, parameters.OptionsJson)
	default : err = quizzes.ErrInvalidQuestionType
	}

	if err != nil{
		return err
	}

	//create strategic question if needed

	//calculate ordering and score value

	//write question


	return nil
}




//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
///				PRIVATE METHODS				/////


func (s *UserService) trueFalseCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string) (string, error){
	return "", nil
}

func (s *UserService) sliderCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string) (string, error){
	return "", nil
}

func (s *UserService) orderingCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string) (string, error){
	return "", nil
}

func (s *UserService) openCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string) (string, error){
	return "", nil
}

func (s *UserService) multipleChCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string) (string, error){
	return "", nil
}

func (s *UserService) matchingCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string) (string, error){
	return "", nil
}

func (s *UserService) dragAndDropCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string) (string, error){
	return "", nil
}