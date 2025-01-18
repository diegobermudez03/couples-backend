package handlers

import (
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/pkg/users"
	"github.com/go-chi/chi/v5"
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


////// HANDLERS
func (h *UsersHandler) PatchPartnersNickNameEndpoint(w http.ResponseWriter, r *http.Request){
	
}
