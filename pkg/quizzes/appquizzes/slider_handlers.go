package appquizzes

import (
	"context"
	"io"

	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
)

func (s *UserService) sliderCreator(ctx context.Context, quiz *quizzes.QuizPlainModel, optionsJSON string, images map[string]io.Reader) (string, error) {
	return "{}", nil
}
