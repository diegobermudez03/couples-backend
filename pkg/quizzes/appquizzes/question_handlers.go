package appquizzes

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/google/uuid"
)

type inputOption struct{
	Text 		string	`json:"text" validate:"required"`
	ImageName	string	`json:"imageName"`
}

type questionOption struct{
	OptId		int 	`json:"optId" validate:"required"`
	Text 		string	`json:"txt" validate:"required"`
	ImageId		*uuid.UUID	`json:"imId"`
	ImageUrl	*string		`json:"imUrl"`
}

// ORDERING QUESTIONS MODELS
type orderingInput struct{
	SortingType		string 			`json:"sortingType" validate:"required"`
	Options 		[]inputOption	`json:"options" validate:"required"`
}

type orderingOptionsFormat struct{
	SortingType		string 				`json:"sortTp" validate:"required"`
	Options 		[]questionOption	`json:"opts" validate:"required"`
}

// OPEN QUESTION MODELS 

type openInput struct{
	NumAnswers	int 	`json:"numAnswers" validate:"required"`
}

type openOptionsFormat struct{
	NumAnswers int		`json:"nAnsw" validate:"required"`
}


// MULTIPLE CHOICE OPTIONS


type multipleInput struct {
	MultipleAnswer 		*bool 	`json:"multipleAnswer"`
	Options 			[]inputOption	`json:"options" validate:"required"`
}

type multipleOptionsFormat struct{
	MultipleAnswer 		bool 	`json:"multAns" validate:"required"`
	Options 			[]questionOption	`json:"opts" validate:"required"`
}

//MATCHING OPTIONS

type matchingInput struct{
	Options1 		[]inputOption	`json:"options1" validate:"required"`
	Options2 		[]inputOption	`json:"options2" validate:"required"`
}

type matchingOptionsFormat struct{
	Options1		[]questionOption	`json:"opts1" validate:"required"`
	Options2		[]questionOption	`json:"opts2" validate:"required"`
}

// DRAG AND DROP OPTIONS

type dragAndDropInput struct{
	Boxes 		[]inputOption	`json:"boxes" validate:"required"`
	Options		[]inputOption	`json:"options" validate:"required"`
}


type dragAndDropOptionsFormat struct{
	Boxes 	[]questionOption	`json:"boxes" validate:"required"`
	Options []questionOption	`json:"options" validate:"required"`

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///// 			CREATORS 			//////

func (s *UserService) trueFalseCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	return "{}", nil
}

func (s *UserService) sliderCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	return "{}", nil
}

func (s *UserService) orderingCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	var input orderingInput
	if err := s.readJson(optionsJSON, &input); err != nil{
		return "", quizzes.ErrInvalidQuestionOptions
	}

	if input.SortingType != quizzes.LEAST_TO_MOST && input.SortingType != quizzes.MOST_TO_LEAST{
		return "", quizzes.ErrInvalidQuestionOptions
	}

	//create base for or
	var output orderingOptionsFormat
	output.SortingType = input.SortingType

	output.Options = make([]questionOption, 0, len(input.Options))
	output.Options = s.readOptions(ctx, input.Options, output.Options, images, quiz.Id,questionId)

	jsonBytes, err := json.Marshal(output)
	if err != nil{
		return "", quizzes.ErrCreatingQuestion
	}
	return string(jsonBytes), nil
}



func (s *UserService) openCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	var input openInput
	if err := s.readJson(optionsJSON, &input); err != nil{
		return "", quizzes.ErrInvalidQuestionOptions
	}

	var output openOptionsFormat
	output.NumAnswers = input.NumAnswers

	jsonBytes, err := json.Marshal(output)
	if err != nil{
		return "", quizzes.ErrCreatingQuestion
	}
	return string(jsonBytes), nil
}


func (s *UserService) multipleChCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	var input multipleInput
	if err := s.readJson(optionsJSON, &input); err != nil{
		return "", quizzes.ErrInvalidQuestionOptions
	}

	output := multipleOptionsFormat{
		MultipleAnswer: false,
	}
	if input.MultipleAnswer != nil{
		output.MultipleAnswer = *input.MultipleAnswer
	}

	output.Options = make([]questionOption, 0, len(input.Options))
	output.Options = s.readOptions(ctx, input.Options, output.Options, images, quiz.Id, questionId)

	jsonBytes, err := json.Marshal(output)
	if err != nil{
		return "", quizzes.ErrCreatingQuestion
	}
	return string(jsonBytes), nil
}


func (s *UserService) matchingCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	var input matchingInput
	if err := s.readJson(optionsJSON, &input); err != nil{
		return "", quizzes.ErrInvalidQuestionOptions
	}
	numberOptions := len(input.Options1)
	if numberOptions != len(input.Options2){
		return "", quizzes.ErrInvalidQuestionOptions
	}

	output := matchingOptionsFormat{}

	output.Options1 = make([]questionOption, 0, numberOptions)
	output.Options1 = make([]questionOption, 0,  numberOptions)
	output.Options1 = s.readOptions(ctx, input.Options1, output.Options1, images, quiz.Id, questionId)
	output.Options2 = s.readOptions(ctx, input.Options2, output.Options1, images, quiz.Id, questionId)
	
	jsonBytes, err := json.Marshal(output)
	if err != nil{
		return "", quizzes.ErrCreatingQuestion
	}
	return string(jsonBytes), nil
}


func (s *UserService) dragAndDropCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	var input dragAndDropInput
	if err := s.readJson(optionsJSON, &input); err != nil{
		return "", quizzes.ErrInvalidQuestionOptions
	}

	output := dragAndDropOptionsFormat{}

	output.Boxes = make([]questionOption,  0, len(input.Boxes))
	output.Options = make([]questionOption, 0, len(input.Options))
	output.Boxes = s.readOptions(ctx, input.Boxes, output.Boxes, images, quiz.Id, questionId)
	output.Options = s.readOptions(ctx, input.Options, output.Options, images, quiz.Id, questionId)

	jsonBytes, err := json.Marshal(output)
	if err != nil{
		return "", quizzes.ErrCreatingQuestion
	}
	return string(jsonBytes), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///// 			DELETORS 			//////

func (s *UserService)  deleteTrueFalse(ctx context.Context, question *quizzes.QuestionPlainModel) error{
	return nil
}

func (s *UserService)  deleteSlider(ctx context.Context, question *quizzes.QuestionPlainModel) error{
	return nil
}

func (s *UserService)  deleteOpen(ctx context.Context, question *quizzes.QuestionPlainModel) error{
	return nil
}

func (s *UserService)  deleteOrdering(ctx context.Context, question *quizzes.QuestionPlainModel) error{
	var options orderingOptionsFormat
	if err := json.Unmarshal([]byte(question.OptionsJson), &options); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	if err := s.deleteOptionsImages(ctx,options.Options); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	return nil
}

func (s *UserService)  deleteMultipleCh(ctx context.Context, question *quizzes.QuestionPlainModel) error{
	var options multipleOptionsFormat
	if err := json.Unmarshal([]byte(question.OptionsJson), &options); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	if err := s.deleteOptionsImages(ctx,options.Options); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	return nil
}

func (s *UserService)  deleteMatching(ctx context.Context, question *quizzes.QuestionPlainModel) error{
	var options matchingOptionsFormat
	if err := json.Unmarshal([]byte(question.OptionsJson), &options); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	if err := s.deleteOptionsImages(ctx,options.Options1); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	if err := s.deleteOptionsImages(ctx,options.Options2); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	return nil
}

func (s *UserService)  deleteDragAndDrop(ctx context.Context, question *quizzes.QuestionPlainModel) error{
	var options matchingOptionsFormat
	if err := json.Unmarshal([]byte(question.OptionsJson), &options); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	if err := s.deleteOptionsImages(ctx,options.Options1); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	if err := s.deleteOptionsImages(ctx,options.Options2); err != nil{
		return quizzes.ErrDeletingQuestion
	}
	return nil
}


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///// 			HELPERS 			//////


func (s *UserService) readOptions(ctx context.Context, inputOptions []inputOption, outputOptions []questionOption, images map[string]io.Reader, quizId uuid.UUID, questionId uuid.UUID) []questionOption{
	lock := sync.Mutex{}
	waitGroup := sync.WaitGroup{}
	for ind, inputOption := range inputOptions{
		option := questionOption{
			OptId: ind,
		}
		option.Text = inputOption.Text

		if inputOption.ImageName != ""{
			waitGroup.Add(1)
			go func(){
				file, ok := images[inputOption.ImageName]
				if ok{
					fileId, url, err := s.fileService.UploadImage(ctx, file, files.MAX_SIZE_QUESTION_PICTURE, true,  s.getOptionImagePath(quizId,questionId, inputOption.ImageName)...)
					if err == nil{
						option.ImageId = fileId
						option.ImageUrl = url
					}
				}
				lock.Lock()
				outputOptions = append(outputOptions, option)
				lock.Unlock()
				waitGroup.Done()
			}()
		}else{
			lock.Lock()
			outputOptions = append(outputOptions, option)
			lock.Unlock()
		}
	}
	waitGroup.Wait()
	return outputOptions
}


func (s *UserService) deleteOptionsImages(ctx context.Context, options []questionOption) error{
	for _, opt := range options{
		if opt.ImageId != nil{
			if err := s.fileService.DeleteImage(ctx, *opt.ImageId); err != nil{
				return quizzes.ErrDeletingQuestion
			}
		}
	}
	return nil
}