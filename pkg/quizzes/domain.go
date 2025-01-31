package quizzes

import "context"

type AdminService interface {
	CreateQuizCategory(ctx context.Context, name, description, imageType string, imageBytes []byte) error
}

type Service interface{}