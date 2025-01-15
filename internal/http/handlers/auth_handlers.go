package handlers

import (
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/auth/domainauth"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authService 	domainauth.AuthService
}

func NewAuthHandler(authService domainauth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(r *chi.Mux){
	router := chi.NewMux()
	r.Mount("/auth", router)

	router.Post("/register", h.registerEndpoint)
}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//////////////////			ENDPOINTS						///////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

type registerDTO struct{
	Email 		string	`json:"email" validate:"email"`
	Password 	string 	`json:"password"`
	Device 		string 	`json:"device" validate:"required"`
	Os			string	`json:"os" validate:"required"`
}


func (h *AuthHandler) registerEndpoint(w http.ResponseWriter, r *http.Request){
	// extract payload
	payload := registerDTO{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	// call service
	refreshToken, err := h.authService.RegisterUser(
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