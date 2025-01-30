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

/////////////////////////////////// ERRORS CODES

var usersErrorCodes = map[error] int{
	users.ErrorUnableCreateUser : http.StatusInternalServerError,
	users.ErrorInvalidCountryCode : http.StatusBadRequest,
	users.ErrorInvalidLanguageCode : http.StatusBadRequest,
	users.ErrorInvalidGender : http.StatusBadRequest,
	users.ErrorTooYoung : http.StatusBadGateway,
	users.ErrorNoCoupleFound : http.StatusNotFound,
	users.ErrorUserHasActiveCouple : http.StatusConflict,
	users.ErrorDeletingUser : http.StatusInternalServerError,
	users.ErrorCreatingTempCouple : http.StatusInternalServerError,
	users.ErrorInvalidCode : http.StatusBadRequest,
	users.ErrorCantConnectWithYourself : http.StatusBadRequest,
	users.ErrorConnectingCouple : http.StatusInternalServerError,
	users.ErrorCreatingPoints : http.StatusInternalServerError,
	users.ErrorUpdatingNickname : http.StatusInternalServerError,
	users.ErrorUnableToCheckPartnerNickname : http.StatusInternalServerError,
	users.ErrorUnableToGetTempCouple : http.StatusInternalServerError,
	users.ErrorNoTempCoupleFound : http.StatusInternalServerError,
	users.ErrorCantGetCouple: http.StatusInternalServerError,
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
		errorCode, ok := usersErrorCodes[err]
		if !ok{
			errorCode = 500
		}
		utils.WriteError(w, errorCode, err)
		return 
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
