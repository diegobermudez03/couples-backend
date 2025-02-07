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

type matchingInputOption struct{
	Text 		string `json:"text" validate:"required"`
	ImageName 	string 	`json:"imageName"`
}
type matchingInput struct{
	Options1 		[]matchingInputOption	`json:"options1" validate:"required"`
	Options2 		[]matchingInputOption	`json:"options2" validate:"required"`
}

type matchingOptionsFormat struct{
	Options1		[]matchingOption	`json:"opts1" validate:"required"`
	Options2		[]matchingOption	`json:"opts2" validate:"required"`
}

type matchingOption struct{
	OptId 				int 		`json:"optId" validate:"required"`
	Text 				string 		`json:"text" validate:"required"`
	ImageId 			*uuid.UUID	`json:"imId"`
	ImageUrl 			*string 	`json:"imUrl"`
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

	output.Options1 = make([]matchingOption, numberOptions)
	output.Options1 = make([]matchingOption, numberOptions)
	lock := sync.Mutex{}
	waitGroup := sync.WaitGroup{}
	
	allOptions := input.Options1
	allOptions = append(allOptions, input.Options2...)
	for ind, inputOption := range allOptions{
		option := matchingOption{
			OptId: ind,
		}
		option.Text = inputOption.Text

		options := &output.Options1
		if ind >= numberOptions{
			options = &output.Options2
		}

		if inputOption.ImageName != ""{
			waitGroup.Add(1)
			go func(){
				file, ok := images[inputOption.ImageName]
				if ok{
					fileId, url, err := s.fileService.UploadImage(ctx, file, files.MAX_SIZE_QUESTION_PICTURE, true, s.getOptionImagePath(quiz.Id,questionId, inputOption.ImageName)...)
					if err == nil{
						option.ImageId = fileId
						option.ImageUrl = url
					}
				}
				lock.Lock()
				*options = append(*options, option)
				lock.Unlock()
				waitGroup.Done()
			}()
		}else{
			lock.Lock()
			*options = append(*options, option)
			lock.Unlock()
		}
	}
	waitGroup.Wait()

	jsonBytes, err := json.Marshal(output)
	if err != nil{
		return "", quizzes.ErrCreatingQuestion
	}
	return string(jsonBytes), nil
}
