package handlers

import (
	"errors"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/go-chi/chi/v5"
)

type QuizzesHandler struct {
	service quizzes.AdminService
	middlewares *middlewares.Middlewares
}

func NewQuizzesHandler(service quizzes.AdminService, middlewares *middlewares.Middlewares) *QuizzesHandler {
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
type PostCategoryDTO struct{
	Name 		string 	`json:"name" validate:"required"`
	Description string 	`json:"description" validate:"required"`
}


/////////////////////////////////// ERRORS CODES

var quizzessErrorCodes = map[error] int{
	quizzes.ErrCategoryAlreadyExists : http.StatusConflict,
	quizzes.ErrMissingCategoryAttributes : http.StatusBadRequest,
	quizzes.ErrCreatingCategory : http.StatusInternalServerError,
}


///////////////////////////////// HANDLERS 

func (h *QuizzesHandler) PostAdminQuizCategory(w http.ResponseWriter, r *http.Request){
	const maxUploadSize = 5 << 20 //5MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	// reading json
	var payload PostCategoryDTO
	if err :=utils.ReadFormJson(r, "category", &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// reading image
	file, _, err := r.FormFile("image")
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, errors.New("MISSING_FIELDS"))
		return
	}
	defer file.Close()

	// i could check image type here, howeever I'll leave it to the files service
	err = h.service.CreateQuizCategory(r.Context(), payload.Name, payload.Description, file)
	if err != nil{
		code, ok := quizzessErrorCodes[err]
		if !ok{
			code = 500
		}
		utils.WriteError(w, code, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *QuizzesHandler) PostAdminQuiz(w http.ResponseWriter, r *http.Request){}