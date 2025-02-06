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


type orderingInput struct{
	SortingType		string 	`json:"sortingType" validate:"required"`
	Options 		[]struct{
		Text 		string	`json:"text" validate:"required"`
		ImageName	string	`json:"imageName"`
	}
}

type orderingOptionsFormat struct{
	SortingType		string 	`json:"sortTp" validate:"required"`
	Options 		[]orderingOption
}

type orderingOption struct{
	OptId		int 	`json:"optId" validate:"required"`
	Text 		string	`json:"txt" validate:"required"`
	ImageId		*uuid.UUID	`json:"imId"`
	ImageUrl	*string		`json:"imUrl"`
}

func (s *UserService) orderingCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	var input orderingInput
	if err := json.Unmarshal([]byte(optionsJSON), &input); err != nil{
		return "", quizzes.ErrInvalidQuestionOptions
	}

	if input.SortingType != quizzes.LEAST_TO_MOST && input.SortingType != quizzes.MOST_TO_LEAST{
		return "", quizzes.ErrInvalidQuestionOptions
	}

	//create base for or
	var output orderingOptionsFormat
	output.SortingType = input.SortingType
	output.Options = make([]orderingOption, len(input.Options))

	waitGroup := sync.WaitGroup{}

	for _, inputOption := range input.Options{
		option := orderingOption{}
		option.Text = inputOption.Text

		if inputOption.ImageName != ""{
			waitGroup.Add(1)
			go func(){
				file, ok := images[inputOption.ImageName]
				if ok{
					fileId, err := s.fileService.UploadImage(ctx, file, files.MAX_SIZE_QUESTION_PICTURE, true, quizzes.DOMAIN_NAME, quizzes.QUIZZES, quiz.Id.String(), questionId.String(), inputOption.ImageName)
					if err == nil{
						option.ImageId = fileId
						//add url
					}
				}
			}()
		}else{
			output.Options = append(output.Options, option)
		}
	}
	waitGroup.Wait()
	return "", nil
}
