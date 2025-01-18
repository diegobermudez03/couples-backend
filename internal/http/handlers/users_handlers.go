package handlers

import (
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsersHandler struct {
	service users.UsersService
	middlewares *middlewares.Middlewares
}


func NewUsersHandler(service users.UsersService, middlewares *middlewares.Middlewares) *UsersHandler{
	return &UsersHandler{
		service: service,
		middlewares: middlewares,
	}
}

func (h *UsersHandler) RegisterRoutes(r *chi.Mux){
	router := chi.NewMux()
	router.Use(h.middlewares.CheckAccessToken)

	r.Mount("/users", router)
	
	router.Patch("/partners/nickname", h.PatchPartnersNickNameEndpoint)
}


///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//////////////////			ENDPOINTS						///////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

///// DTOS
type patchNicknameDTO struct{
	Nickname 	string 	`json:"nickname" validate:"required"`
}

////// HANDLERS
func (h *UsersHandler) PatchPartnersNickNameEndpoint(w http.ResponseWriter, r *http.Request){
	userId := r.Context().Value(middlewares.UserIdKey).(uuid.UUID)
	coupleId := r.Context().Value(middlewares.CoupleIdKey).(uuid.UUID)

	payload := patchNicknameDTO{}
	if err := utils.ReadJSON(r, &payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	if err := h.service.EditPartnersNickname(r.Context(), userId, coupleId, payload.Nickname); err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
