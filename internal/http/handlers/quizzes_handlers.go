package handlers

import (
	"encoding/json"
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
}


///////////////////////////////// HANDLERS 

func (h *QuizzesHandler) PostAdminQuizCategory(w http.ResponseWriter, r *http.Request){
	const maxUploadSize = 50 << 20 //6MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}
	cat := r.FormValue("category")
	if cat == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("MISSING_FIELDS"))
		return 
	}
	var payload PostCategoryDTO
	if err := json.Unmarshal([]byte(cat), &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, errors.New("MISSING_FIELDS"))
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, errors.New("MISSING_FIELDS"))
		return
	}
	defer file.Close()

	//perhaps here I'll add type of image verification
	err = h.service.CreateQuizCategory(r.Context(), payload.Name, payload.Description, file)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *QuizzesHandler) PostAdminQuiz(w http.ResponseWriter, r *http.Request){}