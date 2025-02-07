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

type dragAndDropInput struct{
	Boxes 		[]struct{
		Name 	string 	`json:"name" validate:"required"`
		ImageName 	string `json:"imageName" validate:"required"`
	}	`json:"boxes" validate:"required"`
	Options		[]struct{
		Text 	string 	`json:"text" validate:"required"`
		ImageName string 	`json:"imageName"`
	}	`json:"options" validate:"required"`
}


type dragAndDropOptionsFormat struct{
	Boxes 	[]dragAndDropBox	`json:"boxes" validate:"required"`
	Options []dragAndDropOption	`json:"options" validate:"required"`

}

type dragAndDropBox struct{
	BoxId 		int 		`json:"bxId" validate:"required"`
	Name 		string 		`json:"name" validate:"required"`
	ImageId 	*uuid.UUID	`json:"imId" validate:"required"`
	ImageUrl 	*string 	`json:"imUrl" validate:"required"`
}

type dragAndDropOption struct{
	OptId 		int			`json:"optId" validate:"required"`
	Text 		string 		`json:"txt" validate:"required"`
	ImageId 	*uuid.UUID	`json:"imId" validate:"required"`
	ImageUrl 	*string 	`json:"imUrl" validate:"required"`
}

func (s *UserService) dragAndDropCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader, questionId uuid.UUID) (string, error) {
	var input dragAndDropInput
	if err := s.readJson(optionsJSON, &input); err != nil{
		return "", quizzes.ErrInvalidQuestionOptions
	}

	output := dragAndDropOptionsFormat{}

	output.Boxes = make([]dragAndDropBox, len(input.Boxes))
	output.Options = make([]dragAndDropOption, len(input.Options))
	lock := sync.Mutex{}
	waitGroup := sync.WaitGroup{}
	

	for ind, inputBox := range input.Boxes{
		box := dragAndDropBox{
			BoxId: ind,
		}
		box.Name = inputBox.Name

		if inputBox.ImageName != ""{
			waitGroup.Add(1)
			go func(){
				file, ok := images[inputBox.ImageName]
				if ok{
					fileId, url, err := s.fileService.UploadImage(ctx, file, files.MAX_SIZE_QUESTION_PICTURE, true, s.getOptionImagePath(quiz.Id,questionId, inputBox.ImageName)...)
					if err == nil{
						box.ImageId = fileId
						box.ImageUrl = url
					}
				}
				lock.Lock()
				output.Boxes = append(output.Boxes, box)
				lock.Unlock()
				waitGroup.Done()
			}()
		}else{
			lock.Lock()
			output.Boxes = append(output.Boxes, box)
			lock.Unlock()
		}
	}


	for ind, inputOption := range input.Options{
		option := dragAndDropOption{
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