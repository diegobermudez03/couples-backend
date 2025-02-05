package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const CAT_ID_URL_PARAM = "catId"
const QUIZ_ID_URL_PARAM = "quizId"

type QuizzesHandler struct {
	service     quizzes.UserService
	middlewares *middlewares.Middlewares
}

func NewQuizzesHandler(service quizzes.UserService, middlewares *middlewares.Middlewares) *QuizzesHandler {
	return &QuizzesHandler{
		service:     service,
		middlewares: middlewares,
	}
}

func (h *QuizzesHandler) RegisterRoutes(r *chi.Mux) {
	router := chi.NewMux()
	router.Use(h.middlewares.CheckAccessToken)

	r.Mount("/quizzes", router)

	router.Post(fmt.Sprintf("/{%s}/questions", QUIZ_ID_URL_PARAM), h.postQuestionHandler)
}


///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//////////////////			ENDPOINTS						///////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

///// DTOS

type postQuestionDTO struct{
	Question 			string 		`json:"question" validate:"required"`
	QuestionType		string 		`json:"questionType" validate:"required"`
	ScoreValue			*string 	`json:"scoreValue"`
	OptionsJson			string 		`json:"optionsJson" validate:"required"`
	TimerSeconds 		*int 		`json:"timerSeconds"`
	StrategicAnswerId 	*uuid.UUID	`json:"strategicAnswerId"`
	StrategicName 		*string 	`json:"strategicName"`
	StrategicDescription *string	`json:"strategicDescription"`
}

/////////////////////////////////// ERRORS CODES

var quizzessErrorCodes = map[error] int{
	quizzes.ErrCategoryAlreadyExists : http.StatusConflict,
	quizzes.ErrMissingCategoryAttributes : http.StatusBadRequest,
	quizzes.ErrCreatingCategory : http.StatusInternalServerError,
}



///////////////////////////////// ADMIN HANDLERS	/////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////


func (h *QuizzesHandler) postQuestionHandler(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r, QUIZ_ID_URL_PARAM)
	parsedId, err := uuid.Parse(id)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptuQuizId)
		return
	}

	const maxUploadSize = 20 << 20	//15MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrFileTooBig)
	}
	
	var payload postQuestionDTO
	if err := utils.ReadFormJson(r, "question", &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	//reading all images passed
	images := map[string]io.Reader{}

	files := r.MultipartForm.File["images"]
	for _, header := range files{
		file, err := header.Open()
		if err != nil{
			continue 
		}
		defer file.Close()
		images[header.Filename] = file
	}
	
	err = h.service.CreateQuestion(
		r.Context(), 
		parsedId,
		quizzes.CreateQuestionRequest{
			Question: payload.Question,
			QType: payload.QuestionType,
			ScoreValue: payload.ScoreValue,
			Timer: payload.TimerSeconds,
			StrategicAnswerId: payload.StrategicAnswerId,
			StrategicName: payload.StrategicName,
			StrategicDescription: payload.StrategicDescription,
		},
		images,
	)
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