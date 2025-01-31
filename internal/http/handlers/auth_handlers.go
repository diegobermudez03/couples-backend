package handlers

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService 	auth.AuthService
	adminService 	auth.AuthAdminService
	middlewares 	*middlewares.Middlewares
}

func NewAuthHandler(authService auth.AuthService, adminService 	auth.AuthAdminService, middlewares *middlewares.Middlewares) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		middlewares: middlewares,
		adminService :adminService,
	}
}

func (h *AuthHandler) RegisterRoutes(r *chi.Mux){
	router := chi.NewMux()
	r.Mount("/auth", router)

	router.Post("/register", h.registerEndpoint)
	router.Post("/login", h.LoginEndpoint)
	router.Post("/users", h.createUserEndpoint)
	router.Get("/users/status", h.checkExistanceEndpoint)
	router.Delete("/users/logout", h.userLogoutEndpoint)
	router.Post("/couples/temporal", h.postTempCoupleCodeEndpoint)
	router.Get("/couples/temporal", h.getTempCoupleCodeEndpoint)
	router.Get("/couples/temporal/notification", h.suscribeTempCoupleNotifications)
	router.Post("/couples/connect", h.connectWithCoupleEndpoint)
	router.Post("/accessToken", h.postAccessTokenEndpoint)
	router.With(h.middlewares.CheckAccessToken).Delete("/logout", h.logoutEndpoint)


	router.Post("/admin/accessToken", h.postAdminAccessTokenEndpoint)
}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//////////////////			ENDPOINTS						///////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

///// DTOS
type registerDTO struct{
	Email 		string	`json:"email" validate:"email"`
	Password 	string 	`json:"password"`
	Device 		string 	`json:"device" validate:"required"`
	Os			string	`json:"os" validate:"required"`
}

type loginDTO struct{
	Email 		string	`json:"email" validate:"email"`
	Password 	string 	`json:"password"`
	Device 		string 	`json:"device" validate:"required"`
	Os			string	`json:"os" validate:"required"`
}

type createUserDTO struct{
	FirstName 		string 	`json:"firstName" validate:"required"`
	LastName 		string 	`json:"lastName" validate:"required"`
	Gender 			string	`json:"gender" validate:"required"`
	BirthDate 		int 	`json:"birthDate" validate:"required"`
	CountryCode 	string 	`json:"countryCode" validate:"required"`
	LanguageCode 	string 	`json:"languageCode" validate:"required"`
}

type tempCoupleDTO struct{
	Code 		int 	`json:"code"`
	StartDate 	int64 	`json:"startDate"`
}

/////////////////////////////////// ERRORS CODES

var authErrorCodes = map[error] int{
	auth.ErrorCreatingSession : http.StatusInternalServerError,
	auth.ErrorCreatingAccessToken : http.StatusInternalServerError,
	auth.ErrorCheckingStatus : http.StatusInternalServerError,
	auth.ErrorcreatingTempCouple : http.StatusInternalServerError,
	auth.ErrorNoUserFoundEmail : http.StatusNotFound,
	auth.ErrorNoActiveUser : http.StatusNotFound,
	auth.ErrorNoActiveCoupleFromUser : http.StatusNotFound,
	auth.ErrorExpiredRefreshToken : http.StatusBadRequest,
	auth.ErrorExpiredAccessToken : http.StatusBadRequest,
	auth.ErrorMalformedAccessToken : http.StatusBadRequest,
	auth.ErrorCantLogoutAnonymousAcc : http.StatusBadRequest,
	auth.ErrorIncorrectPassword : http.StatusBadRequest,
	auth.ErrorInsecurePassword : http.StatusBadRequest,
	auth.ErrorEmailAlreadyUsed : http.StatusBadRequest,
	auth.ErrorCreatingAccount : http.StatusInternalServerError,
	auth.ErrorVinculatingAccount : http.StatusInternalServerError,
	auth.ErrorCreatingUser : http.StatusInternalServerError,
	auth.ErrorWithLogin : http.StatusInternalServerError,
	auth.ErrorInvalidRefreshToken : http.StatusBadRequest,
	auth.ErrorWithLogout : http.StatusInternalServerError,
	auth.ErrorNonExistingSession : http.StatusNotFound,
	auth.ErrorUserForAccountAlreadyExists : http.StatusConflict,
	auth.ErrorUnableToConnectCouple : http.StatusInternalServerError,
	auth.ErrorGettingTempCouple : http.StatusInternalServerError,
	auth.ErrTempCoupleNotFound : http.StatusNotFound,
	auth.ErrCantCreateNewCouple : http.StatusInternalServerError,
	auth.ErrUnableToSuscribe : http.StatusInternalServerError,
	auth.ErrNoCodeToSuscribe : http.StatusNotFound,
	auth.ErrNonExistingCode : http.StatusBadRequest,
	auth.ErrCantConnectWithYourself : http.StatusBadRequest,
}


///////////////////////////////// HANDLERS 

func (h *AuthHandler) registerEndpoint(w http.ResponseWriter, r *http.Request){
	// extract payload
	payload := registerDTO{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	token := r.Header.Get("token")
	// call service
	refreshToken, err := h.authService.RegisterUserAuth(
		r.Context(),
		payload.Email,
		payload.Password,
		payload.Device,
		payload.Os,
		token,
	)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return  
	}

	// Respond
	utils.WriteJSON(
		w, 
		http.StatusCreated, 
		map[string]any{
			"refreshToken" : refreshToken,
		},
	)
}


func (h *AuthHandler) LoginEndpoint(w http.ResponseWriter, r *http.Request){
	dto := loginDTO{}
	if err := utils.ReadJSON(r, &dto); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	refreshToken, err := h.authService.LoginUserAuth(r.Context(), dto.Email, dto.Password, dto.Device, dto.Os)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(
		w, 
		http.StatusOK, 
		map[string]any{
			"refreshToken" : refreshToken,
		},
	)
}


func (h *AuthHandler) createUserEndpoint(w http.ResponseWriter, r *http.Request){
	payload := createUserDTO{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	token := r.Header.Get("token")
	token, err := h.authService.CreateUser(
		r.Context(), 
		token,
		payload.FirstName,
		payload.LastName,
		payload.Gender,
		payload.CountryCode,
		payload.LanguageCode,
		payload.BirthDate,
	)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"refreshToken" : token,
	})
}

func (h *AuthHandler) checkExistanceEndpoint(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}

	status, err := h.authService.CheckUserAuthStatus(r.Context(), token)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status" : status,
	})
}

func (h *AuthHandler) userLogoutEndpoint(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}

	if err := h.authService.CloseUsersSession(r.Context(), token); err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) getTempCoupleCodeEndpoint(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}
	tempCouple, err := h.authService.GetTempCoupleOfUser(r.Context(), token)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	dto := tempCoupleDTO{
		Code: tempCouple.Code,
		StartDate: tempCouple.StartDate.Unix(),
	}
	utils.WriteJSON(
		w, http.StatusOK, dto,
	);
}


func (h *AuthHandler) postTempCoupleCodeEndpoint(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}
	payload := struct{
		StartDate	int `json:"startDate" validate:"required"`
	}{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	code, err := h.authService.CreateTempCouple(r.Context(), token, payload.StartDate)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(
		w, http.StatusCreated,
		map[string]int{
			"code" : code,
		},
	)
}


func (h *AuthHandler) connectWithCoupleEndpoint(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}

	payload := struct{
		Code 	int 	`json:"code" validate:"required"`
	}{}
	if err := utils.ReadJSON(r,&payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	accessToken, err := h.authService.ConnectCouple(r.Context(), token, payload.Code)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"accessToken" : accessToken})
}

func (h *AuthHandler) postAccessTokenEndpoint(w http.ResponseWriter, r *http.Request){
	payload := struct{
		RefreshToken 	string 	`json:"refreshToken" validate:"required"`
	}{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	accessToken, newRToken, err := h.authService.CreateAccessToken(r.Context(), payload.RefreshToken)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(
		w, 
		http.StatusCreated,
		 map[string]any{
			"accessToken" : accessToken,
			"refreshToken" : newRToken,
		},
	)
}


func (h *AuthHandler) postAdminAccessTokenEndpoint(w http.ResponseWriter, r *http.Request){
	payload := struct{
		RefreshToken 	string 	`json:"refreshToken" validate:"required"`
	}{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	accessToken, err := h.adminService.CreateAccessToken(r.Context(), payload.RefreshToken)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(
		w, 
		http.StatusCreated,
		 map[string]any{
			"accessToken" : accessToken,
		},
	)
}


func (h *AuthHandler) logoutEndpoint(w http.ResponseWriter, r *http.Request){
	sessionId := r.Context().Value(middlewares.SessionIdKey).(uuid.UUID)
	if err := h.authService.LogoutSession(r.Context(), sessionId); err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) suscribeTempCoupleNotifications(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}
	channel, userId, err := h.authService.SuscribeTempCoupleNot(r.Context(), token)
	if err != nil{
		errorCode, ok := authErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return
	}
	// SETTING SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	randId := rand.IntN(10000)+1000
	log.Printf("CONNECTED CLIENT %d", randId)
	w.Write([]byte("data: CONNECTED\n\n"))
	flusher.Flush()
	
	clientGone := r.Context().Done()
	select{
	case <- clientGone:
		log.Printf("CLIENT CLOSED SSE CONNECTION ID : %d", randId)
		h.authService.RemoveCodeSuscriber(*userId)
	case received, ok :=<- channel:
		if received == auth.StatusVinculated{
			w.Write([]byte(fmt.Sprintf("data: %s\n\n", received)))
		}else{
			w.Write([]byte("data: ERROR\n\n"))
		}
		flusher.Flush()
		w.Write([]byte("event:close\ndata: Connection closing\n\n"))
		flusher.Flush()
		if ok{
			//we only close remove the suscriber from here if the channel is still opened, if its not it means that someone else closed it
			h.authService.RemoveCodeSuscriber(*userId)
		}
		log.Printf("CLOSING SSE CONNECTION BY SERVER ID: %d", randId)
		//here it closes the connection but is not because of the message but because the handler function ends
	}
}
////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
/////////////////////// PRIVATE FUNCTIONS
