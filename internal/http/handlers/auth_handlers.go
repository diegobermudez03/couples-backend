package handlers

import (
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


///////////////////////////////// HANDLERS 

func (h *AuthHandler) registerEndpoint(w http.ResponseWriter, r *http.Request){
	// extract payload
	payload := registerDTO{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	// call service
	refreshToken, err := h.authService.RegisterUser(
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

	refreshToken, err := h.authService.LoginUser(r.Context(), dto.Email, dto.Password, dto.Device, dto.Os)
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

