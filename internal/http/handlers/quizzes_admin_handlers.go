package handlers

import (
	"fmt"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)


type QuizzesAdminHandler struct {
	service quizzes.AdminService
	middlewares *middlewares.Middlewares
}

func NewQuizzesAdminHandler(service quizzes.AdminService, middlewares *middlewares.Middlewares) *QuizzesAdminHandler {
	return &QuizzesAdminHandler{
		service: service,
		middlewares: middlewares,
	}
}

func (h *QuizzesAdminHandler) RegisterRoutes(r *chi.Mux){
	router := chi.NewMux()
	//router.Use(h.middlewares.CheckAdminAccessToken)

	r.Mount("/quizzes", router)
	r.Mount("/admin/quizzes", router)

	router.Post("/categories", h.postAdminQuizCategory)
	router.Patch(fmt.Sprintf("/categories/{%s}", CAT_ID_URL_PARAM), h.patchAdminQuizCategory)
	router.Post(fmt.Sprintf("/categories/{%s}/quizzes", CAT_ID_URL_PARAM), h.postAdminQuiz)
	router.Patch(fmt.Sprintf("/{%s}", QUIZ_ID_URL_PARAM), h.patchAdminQuiz)
	//router.Delete(fmt.Sprintf("/categories/{%s}", CAT_ID_URL_PARAM), h.deleteAdminQuizCategory)
}


///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//////////////////			ENDPOINTS						///////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

///// DTOS
type postCategoryAdminDTO struct{
	Name 		string 	`json:"name" validate:"required"`
	Description string 	`json:"description" validate:"required"`
}

type patchCategoryAdminDTO struct{
	Name 		string 	`json:"name"`
	Description string 	`json:"description"`
}

type postAdminQuizAdminDTO struct{
	Name 		string 	`json:"name" validate:"required"`
	Description string 	`json:"description" validate:"required"`
	LanguageCode string 	`json:"languageCode" validate:"required"`
}

type patchQuizAdminDTO struct{
	Name 		string 	`json:"name"`
	Description string 	`json:"description"`
	LanguageCode string 	`json:"languageCode"`
	CategoryId 	*uuid.UUID	`json:"categoryId"`
}


///////////////////////////////// ADMIN HANDLERS	/////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

func (h *QuizzesAdminHandler) postAdminQuizCategory(w http.ResponseWriter, r *http.Request){
	const maxUploadSize = 5 << 20 //5MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest,utils.ErrFileTooBig)
		return 
	}

	// reading json
	var payload postCategoryAdminDTO
	if err :=utils.ReadFormJson(r, "category", &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// reading image
	file, _, err := r.FormFile("image")
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrMissingFields)
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


func (h *QuizzesAdminHandler) patchAdminQuizCategory(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r, CAT_ID_URL_PARAM)
	parsedId, err := uuid.Parse(id)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrInvalidId)
		return 
	}

	const maxUploadSize = 5 << 20 //5MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrFileTooBig)
		return 
	}

	// reading json
	var payload patchCategoryAdminDTO
	utils.ReadFormJson(r, "quiz", &payload)

	// reading image
	file, _, err := r.FormFile("image")
	if err != nil{
		file = nil
	}else{
		defer file.Close()
	}

	if err := h.service.UpdateQuizCategory(r.Context(), parsedId, payload.Name, payload.Description, file ); err != nil{
		code, ok := quizzessErrorCodes[err]
		if !ok{
			code = 500
		}
		utils.WriteError(w, code, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *QuizzesAdminHandler) postAdminQuiz(w http.ResponseWriter, r *http.Request){
	categoryId := chi.URLParam(r, CAT_ID_URL_PARAM)
	catParsed, err := uuid.Parse(categoryId)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptyCategoryId)
		return 
	}
	const maxUploadSize = 5 << 20 //5MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrFileTooBig)
		return 
	}

	var payload postAdminQuizAdminDTO
	if err := utils.ReadFormJson(r, "quiz", &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	// reading image
	file, _, err := r.FormFile("image")
	if err != nil{
		file = nil
	}else{
		defer file.Close()
	}

	err = h.service.CreateQuiz(r.Context(), payload.Name, payload.Description, payload.LanguageCode,catParsed, file)
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



func (h *QuizzesAdminHandler) patchAdminQuiz(w http.ResponseWriter, r *http.Request){
	quizId := chi.URLParam(r, QUIZ_ID_URL_PARAM)
	quizParsed, err := uuid.Parse(quizId)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptuQuizId)
		return 
	}
	const maxUploadSize = 5 << 20 //5MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrFileTooBig)
		return 
	}

	var payload patchQuizAdminDTO
	if err := utils.ReadFormJson(r, "quiz", &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	// reading image
	file, _, err := r.FormFile("image")
	if err != nil{
		file = nil
	}else{
		defer file.Close()
	}

	err = h.service.UpdateQuiz(r.Context(), quizParsed, payload.Name, payload.Description,payload.LanguageCode, payload.CategoryId, file)
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

