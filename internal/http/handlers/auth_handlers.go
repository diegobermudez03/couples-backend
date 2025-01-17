package handlers

import (
	"errors"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authService 	auth.AuthService
}

func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(r *chi.Mux){
	router := chi.NewMux()
	r.Mount("/auth", router)

	router.Post("/register", h.registerEndpoint)
	router.Get("/login", h.LoginEndpoint)
	router.Post("/users", h.createUserEndpoint)
	router.Get("/users/exists", h.checkExistanceEndpoint)
	router.Delete("/users/logout", h.logoutEndpoint)
	router.Get("/couples/temporal", h.getTempCoupleCodeEndpoint)
	router.Post("/couples/temporal", h.connectWithCoupleEndpoint)
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

	// call service
	refreshToken, err := h.authService.RegisterUserAuth(
		r.Context(),
		payload.Email,
		payload.Password,
		payload.Device,
		payload.Os,
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

func (h *AuthHandler) logoutEndpoint(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}

	if err := h.authService.CloseSession(r.Context(), token); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *AuthHandler) getTempCoupleCodeEndpoint(w http.ResponseWriter, r *http.Request){
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

	if err := h.authService.ConnectCouple(r.Context(), token, payload.Code); err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

