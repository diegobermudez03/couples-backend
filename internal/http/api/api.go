package api

import (
	"context"
	"net/http"
	"time"

	"github.com/diegobermudez03/couples-backend/internal/config"
	"github.com/diegobermudez03/couples-backend/internal/http/handlers"
	"github.com/diegobermudez03/couples-backend/pkg/auth/appauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type APIServer struct {
	server 	http.Server
	config  *config.Config
}

func NewAPIServer(config *config.Config) *APIServer {
	return &APIServer{
		config: config,
	}
}


func (s *APIServer) Run() error{
	r := chi.NewMux()
	router := chi.NewMux()

	//ROUTER CONFIG
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/v1", router)
	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hellow world"))
	})


	//depencency injections
	s.injectDependencies(router)

	s.server = http.Server{
		Addr: "localhost:" + s.config.Port,
		Handler: r,
	}
	return s.server.ListenAndServe()
}


func (s *APIServer) Shutdown() error{
	return s.server.Shutdown(context.TODO())
}


func (s *APIServer) injectDependencies(router *chi.Mux){
	//create respositories

	//create services
	authService := appauth.NewAuthService()

	//create handlers
	authHandler := handlers.NewAuthHandler(authService)

	//registering routes
	authHandler.RegisterRoutes(router)
}