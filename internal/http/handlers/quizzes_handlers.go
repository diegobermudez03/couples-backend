package handlers

import (
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/go-chi/chi/v5"
)

type QuizzesHandler struct {
	service quizzes.Service
	middlewares middlewares.Middlewares
}

func NewQuizzesHandler(service quizzes.Service, middlewares middlewares.Middlewares) *QuizzesHandler {
	return &QuizzesHandler{
		service: service,
		middlewares: middlewares,
	}
}

func (h *QuizzesHandler) RegisterRoutes(r *chi.Mux){
	routerUsers := chi.NewMux()
	routerAdmin := chi.NewMux()
	routerUsers.With(h.middlewares.CheckAccessToken)

	r.Mount("/quizzes", routerUsers)
	r.Mount("/admin/quizzes", routerAdmin)

	routerAdmin.Post("/categories", h.PostAdminQuizCategory)
	routerAdmin.Post("/", h.PostAdminQuiz)
}


///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//////////////////			ENDPOINTS						///////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

///// DTOS


/////////////////////////////////// ERRORS CODES

var quizzessErrorCodes = map[error] int{
}


///////////////////////////////// HANDLERS 

func (h *QuizzesHandler) PostAdminQuizCategory(w http.ResponseWriter, r *http.Request){}

func (h *QuizzesHandler) PostAdminQuiz(w http.ResponseWriter, r *http.Request){}