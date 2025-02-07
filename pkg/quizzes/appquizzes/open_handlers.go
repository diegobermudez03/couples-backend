package appquizzes

import (
	"context"
	"encoding/json"
	"io"

	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/google/uuid"
)

type openInput struct{
	NumAnswers	int 	`json:"numAnswers" validate:"required"`
}

type openOptionsFormat struct{
	NumAnswers int		`json:"nAnsw" validate:"required"`
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