package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const CAT_ID_URL_PARAM = "catId"
const QUIZ_ID_URL_PARAM = "quizId"
const QUESTION_ID_URL_PARAM = "questionId"

type QuizzesHandler struct {
	service     quizzes.UserService
	middlewares *middlewares.Middlewares
	adminService quizzes.AdminService
}

func NewQuizzesHandler(adminService quizzes.AdminService, service quizzes.UserService, middlewares *middlewares.Middlewares) *QuizzesHandler {
	return &QuizzesHandler{
		adminService: adminService,
		service:     service,
		middlewares: middlewares,
	}
}

func (h *QuizzesHandler) RegisterRoutes(r *chi.Mux) {
	routerUsers := chi.NewMux()
	routerUsers.Use(h.middlewares.CheckAccessToken)
	routerAdmin := chi.NewMux()
	//routerAdmin.Use(h.middlewares.CheckAdminAccessToken)

	r.Mount("/quizzes", routerUsers)

	//	quiz handlers
	routerUsers.Post("/quizes", h.postQuiz)
	routerUsers.With(h.middlewares.CheckUserQuizPermissions).Patch(fmt.Sprintf("/quizes/{%s}", QUIZ_ID_URL_PARAM), h.patchQuizHandler)
	routerUsers.With(h.middlewares.CheckUserQuizPermissions).Delete(fmt.Sprintf("/quizes/{%s}", QUIZ_ID_URL_PARAM), h.deleteQuiz)
	// 	question handlers
	routerUsers.With(h.middlewares.CheckUserQuizPermissions).Post(fmt.Sprintf("/quizes/{%s}/questions", QUIZ_ID_URL_PARAM), h.postQuestionHandler)
	routerUsers.With(h.middlewares.CheckUserQuizPermissions).Patch(fmt.Sprintf("/questions/{%s}", QUESTION_ID_URL_PARAM), h.patchQuestion)
	routerUsers.With(h.middlewares.CheckUserQuizPermissions).Delete(fmt.Sprintf("/questions/{%s}", QUESTION_ID_URL_PARAM), h.deleteQuestion)

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	r.Mount("/admin/quizzes", routerAdmin)

	//	categories handlers
	routerAdmin.Post("/categories", h.postAdminQuizCategory)
	routerAdmin.Patch(fmt.Sprintf("/categories/{%s}", CAT_ID_URL_PARAM), h.patchAdminQuizCategory)
	routerAdmin.Delete(fmt.Sprintf("/categories/{%s}", CAT_ID_URL_PARAM), h.deleteCategory)
	//	quiz handlers
	routerAdmin.Post(fmt.Sprintf("/categories/{%s}/quizes", CAT_ID_URL_PARAM), h.postQuiz)
	routerAdmin.Patch(fmt.Sprintf("/quizes/{%s}", QUIZ_ID_URL_PARAM), h.patchQuizHandler)
	routerAdmin.Delete(fmt.Sprintf("/quizes/{%s}", QUIZ_ID_URL_PARAM), h.deleteQuiz)
	//question handlers
	routerAdmin.Post(fmt.Sprintf("/quizes/{%s}/questions", QUIZ_ID_URL_PARAM), h.postQuestionHandler)
	routerUsers.Patch(fmt.Sprintf("/questions/{%s}", QUESTION_ID_URL_PARAM), h.patchQuestion)
	routerAdmin.Delete(fmt.Sprintf("/questions/{%s}", QUESTION_ID_URL_PARAM), h.deleteQuestion)
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
	OptionsJson			map[string]any 	`json:"optionsJson"`
	StrategicAnswerId 	*uuid.UUID	`json:"strategicAnswerId"`
	StrategicName 		*string 	`json:"strategicName"`
	StrategicDescription *string	`json:"strategicDescription"`
}

type patchQuestionDTO struct{
	Question 			*string	`json:"question"`
	OptionsJson 		map[string]any `json:"optionsJson"`
	StrategicAnswerId 	*uuid.UUID	`json:"strategicAnswerId"`
	StrategicName 		*string 	`json:"strategicName"`
	StrategicDescription *string	`json:"strategicDescription"`
}

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

type patchQuizDTO struct{
	Name 		string 	`json:"name"`
	Description string 	`json:"description"`
	LanguageCode string 	`json:"languageCode"`
	CategoryId 	*uuid.UUID	`json:"categoryId"`
}


/////////////////////////////////// ERRORS CODES

var quizzessErrorCodes = map[error] int{
	quizzes.ErrCategoryAlreadyExists : http.StatusConflict,
	quizzes.ErrMissingCategoryAttributes : http.StatusBadRequest,
	quizzes.ErrCreatingCategory : http.StatusInternalServerError,
}



func (h *QuizzesHandler) postAdminQuizCategory(w http.ResponseWriter, r *http.Request){
	const maxUploadSize = 5 << 20 //5MB
	var payload postCategoryAdminDTO
	if err := utils.ParseAndReadMultiPartForm(w, r,maxUploadSize,&payload, "category"); err != nil{
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
	err = h.adminService.CreateQuizCategory(r.Context(), payload.Name, payload.Description, file)
	if err != nil{
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}


func (h *QuizzesHandler) patchAdminQuizCategory(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r, CAT_ID_URL_PARAM)
	parsedId, err := uuid.Parse(id)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrInvalidId)
		return 
	}

	const maxUploadSize = 5 << 20 //5MB
	var payload patchCategoryAdminDTO
	if err := utils.ParseAndReadMultiPartForm(w, r, maxUploadSize,&payload, "quiz"); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// reading image
	file, _, err := r.FormFile("image")
	if err != nil{
		file = nil
	}else{
		defer file.Close()
	}

	if err := h.adminService.UpdateQuizCategory(r.Context(), parsedId, payload.Name, payload.Description, file ); err != nil{
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *QuizzesHandler) postQuiz(w http.ResponseWriter, r *http.Request){
	userId := h.getUserId(r)

	categoryId := chi.URLParam(r, CAT_ID_URL_PARAM)
	catParsed, err := uuid.Parse(categoryId)
	var catPtr *uuid.UUID  = nil
	if err == nil{
		catPtr = &catParsed
	}
	const maxUploadSize = 5 << 20 //5MB
	var payload postAdminQuizAdminDTO
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := utils.ParseAndReadMultiPartForm(w, r, maxUploadSize, &payload, "quiz"); err != nil{
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

	err = h.service.CreateQuiz(r.Context(), payload.Name, payload.Description, payload.LanguageCode, catPtr, userId,file)
	if err != nil{
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}



func (h *QuizzesHandler) patchQuizHandler(w http.ResponseWriter, r *http.Request){
	log.Println("patch quiz")
	quizId := chi.URLParam(r, QUIZ_ID_URL_PARAM)
	quizParsed, err := uuid.Parse(quizId)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptyQuizId)
		return 
	}
	const maxUploadSize = 5 << 20 //5MB
	var payload patchQuizDTO
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := utils.ParseAndReadMultiPartForm(w, r, maxUploadSize, &payload, "quiz"); err != nil{
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
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}


func (h *QuizzesHandler) postQuestionHandler(w http.ResponseWriter, r *http.Request){
	log.Print("post question handler")
	id := chi.URLParam(r, QUIZ_ID_URL_PARAM)
	parsedId, err := uuid.Parse(id)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptyQuizId)
		return
	}

	const maxUploadSize = 20 << 20	//15MB
	var payload postQuestionDTO
	if err := utils.ParseAndReadMultiPartForm(w, r, maxUploadSize, &payload, "question"); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//reading all images passed
	images := map[string]io.Reader{}
	files := r.MultipartForm.File["images"]
	for _, header := range files{
		file, err := header.Open()
		if err == nil{
			defer file.Close()
			images[header.Filename] = file
		}
	}
	
	err = h.service.CreateQuestion(r.Context(), parsedId,
		quizzes.CreateQuestionRequest{
			Question: payload.Question,
			QType: payload.QuestionType,
			OptionsJson: payload.OptionsJson,
			StrategicAnswerId: payload.StrategicAnswerId,
			StrategicName: payload.StrategicName,
			StrategicDescription: payload.StrategicDescription,
		},
		images,
	)
	if err != nil{
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}


func (h *QuizzesHandler) deleteCategory(w http.ResponseWriter, r *http.Request){
	auxId := chi.URLParam(r, CAT_ID_URL_PARAM)
	catId, err := uuid.Parse(auxId)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptyQuestionId)
		return
	}
	if err := h.adminService.DeleteQuizCategory(r.Context(), catId); err != nil{
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return 
	}
	utils.WriteJSON(w, http.StatusOK, nil )
}


func (h *QuizzesHandler) deleteQuiz(w http.ResponseWriter, r *http.Request){
	auxId := chi.URLParam(r, QUIZ_ID_URL_PARAM)
	quizId, err := uuid.Parse(auxId)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptyQuestionId)
		return
	}
	if err := h.service.DeleteQuiz(r.Context(), quizId); err != nil{
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return 
	}
	utils.WriteJSON(w, http.StatusOK, nil )
}


func (h *QuizzesHandler) deleteQuestion(w http.ResponseWriter, r *http.Request){
	auxId := chi.URLParam(r, QUESTION_ID_URL_PARAM)
	questionId, err := uuid.Parse(auxId)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptyQuestionId)
		return
	}
	if err := h.service.DeleteQuestion(r.Context(), questionId); err != nil{
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return 
	}
	utils.WriteJSON(w, http.StatusOK, nil )
}

func (h *QuizzesHandler) patchQuestion(w http.ResponseWriter, r *http.Request){
	auxId := chi.URLParam(r, QUESTION_ID_URL_PARAM)
	questionId, err := uuid.Parse(auxId)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptyQuestionId)
		return
	}
	const maxSize = 15 << 20 //15MB
	var payload patchQuestionDTO
	if err := utils.ParseAndReadMultiPartForm(w, r, maxSize,&payload, "question"); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	//reading all images passed
	images := map[string]io.Reader{}
	files := r.MultipartForm.File["images"]
	for _, header := range files{
		file, err := header.Open()
		if err == nil{
			defer file.Close()
			images[header.Filename] = file
		}
	}

	err = h.service.UpdateQuestion(r.Context(), questionId, quizzes.UpdateQuestionRequest{
		Question: payload.Question,
		OptionsJson: payload.OptionsJson,
		StrategicAnswerId: payload.StrategicAnswerId,
		StrategicName: payload.StrategicName,
		StrategicDescription: payload.StrategicDescription,
	}, images)
	if err != nil{
		code := utils.GetErrorCode(err, quizzessErrorCodes, 500)
		utils.WriteError(w, code, err)
		return 
	}
}

//////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////
func (h *QuizzesHandler) getUserId(r *http.Request) *uuid.UUID{
	userId, ok := r.Context().Value(middlewares.UserIdKey{}).(uuid.UUID)
	var ptUserId *uuid.UUID = nil
	if ok{
		ptUserId = &userId
	}
	return ptUserId
}