package handlers

import (
	"errors"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/users"
	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	service users.UsersService
}


func NewUsersHandler(service users.UsersService) *UsersHandler{
	return &UsersHandler{
		service: service,
	}
}

func (h *UsersHandler) RegisterRoutes(r *chi.Mux){
	router := chi.NewMux()
	r.Mount("/users", router)

	router.Post("/", h.createUserEndpoint)
	router.Get("/exists", h.checkExistanceEndpoint)
	router.Delete("/logout", h.logoutEndpoint)
	router.Post("/couples/temporal", h.createTempCoupleEndpoint)
}


///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//////////////////			ENDPOINTS						///////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

///// DTOS
type createUserDTO struct{
	FirstName 		string 	`json:"firstName" validate:"required"`
	LastName 		string 	`json:"lastName" validate:"required"`
	Gender 			string	`json:"gender" validate:"required"`
	BirthDate 		int 	`json:"birthDate" validate:"required"`
	CountryCode 	string 	`json:"countryCode" validate:"required"`
	LanguageCode 	string 	`json:"languageCode" validate:"required"`
}




////// HANDLERS

func (h *UsersHandler) createUserEndpoint(w http.ResponseWriter, r *http.Request){
	payload := createUserDTO{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	token := r.Header.Get("token")
	token, err := h.service.CreateUser(
		r.Context(), 
		payload.FirstName,
		payload.LastName,
		payload.Gender,
		payload.CountryCode,
		payload.LanguageCode,
		payload.BirthDate,
		token,
	)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"refreshToken" : token,
	})
}

func (h *UsersHandler) checkExistanceEndpoint(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}

	if err := h.service.CheckUserExistance(r.Context(), token); err != nil{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no user associated"))
		return 
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *UsersHandler) logoutEndpoint(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("token")
	if token == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("no token provided"))
		return 
	}

	if err := h.service.CloseOngoingSession(r.Context(), token); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *UsersHandler) createTempCoupleEndpoint(w http.ResponseWriter, r *http.Request){
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

	code, err := h.service.CreateTempCouple(r.Context(), token, payload.StartDate)
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