package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/go-chi/chi/v5"
)

type FilesHandler struct {
	service files.Service
}

func NewFilesHandler(service files.Service) *FilesHandler{
	return &FilesHandler{
		service: service,
	}
}


func (h *FilesHandler) RegisterRoutes(r *chi.Mux){
	router := chi.NewMux()

	r.Mount("/files", router)

	router.Get("/images/*", h.getPublicImage)
}


func (h *FilesHandler) getPublicImage(w http.ResponseWriter, r *http.Request){
	path := chi.URLParam(r, "*")
	if path == ""{
		utils.WriteError(w, http.StatusBadRequest, errors.New("NO_PATH_GIVEN"))
		return 
	}
	file, contentType, err := h.service.GetImage(r.Context(), path)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}
	defer file.Close()
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, file)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, errors.New("UNABLE_TO_LOAD"))
		return 
	}
}