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

type multipleInput struct {
	MultipleAnswer 		*bool 	`json:"multipleAnswer"`
	Options 			[]struct{
		Text 		string 	`json:"text" validate:"required"`
		ImageName 	string	`json:"imageName"`
	}	`json:"options" validate:"required"`
}

type multipleOptionsFormat struct{
	MultipleAnswer 		bool 	`json:"multAns" validate:"required"`
	Options 			[]multipleOption	`json:"opts" validate:"required"`
}

type multipleOption struct{
	OptId 				int 	`json:"optId" validate:"required"`
	Text 				string 	`json:"txt" validate:"required"`
	ImageId 			*uuid.UUID	`json:"imId"`
	ImageUrl 			*string 	`json:"imUrl"`
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

	output.Options = make([]multipleOption, len(input.Options))
	lock := sync.Mutex{}
	waitGroup := sync.WaitGroup{}
	

	for ind, inputOption := range input.Options{
		option := multipleOption{
			OptId: ind,
		}
		option.Text = inputOption.Text

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
				output.Options = append(output.Options, option)
				lock.Unlock()
				waitGroup.Done()
			}()
		}else{
			lock.Lock()
			output.Options = append(output.Options, option)
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