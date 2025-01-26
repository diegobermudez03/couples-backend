package handlers

import (
	"errors"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService 	auth.AuthService
	middlewares 	*middlewares.Middlewares
}

func NewAuthHandler(authService auth.AuthService, middlewares *middlewares.Middlewares) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		middlewares: middlewares,
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
	router.Post("/couples/connect", h.connectWithCoupleEndpoint)
	router.Get("/accessToken", h.getAccessTokenEndpoint)
	router.With(h.middlewares.CheckAccessToken).Delete("/logout", h.logoutEndpoint)

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
		utils.WriteError(w, http.StatusInternalServerError, err)
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
		utils.WriteError(w, http.StatusInternalServerError, err)
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
		utils.WriteError(w, http.StatusInternalServerError, err)
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
		utils.WriteError(w, http.StatusInternalServerError, err)
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
		utils.WriteError(w, http.StatusBadRequest, err)
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
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(
		w,
		http.StatusCreated,
		tempCouple,
	)
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
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(
		w,
		http.StatusCreated,
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
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"accessToken" : accessToken})
}

func (h *AuthHandler) getAccessTokenEndpoint(w http.ResponseWriter, r *http.Request){
	payload := struct{
		RefreshToken 	string 	`json:"refreshToken" validate:"required"`
	}{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	accessToken, err := h.authService.CreateAccessToken(r.Context(), payload.RefreshToken)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"accessToken" : accessToken})
}


func (h *AuthHandler) logoutEndpoint(w http.ResponseWriter, r *http.Request){
	sessionId := r.Context().Value(middlewares.SessionIdKey).(uuid.UUID)
	if err := h.authService.LogoutSession(r.Context(), sessionId); err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}