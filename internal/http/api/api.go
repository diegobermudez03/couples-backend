package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/diegobermudez03/couples-backend/internal/config"
	"github.com/diegobermudez03/couples-backend/internal/http/handlers"
	"github.com/diegobermudez03/couples-backend/internal/http/middlewares"
	"github.com/diegobermudez03/couples-backend/pkg/auth/appauth"
	"github.com/diegobermudez03/couples-backend/pkg/auth/repoauth"
	"github.com/diegobermudez03/couples-backend/pkg/localization/applocalization"
	"github.com/diegobermudez03/couples-backend/pkg/users/appusers"
	"github.com/diegobermudez03/couples-backend/pkg/users/repousers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type APIServer struct {
	server 	http.Server
	config  *config.Config
	db 		*sql.DB
}

func NewAPIServer(config *config.Config, db *sql.DB) *APIServer {
	return &APIServer{
		config: config,
		db: db,
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
	authRepository := repoauth.NewAuthPostgresRepo(s.db)
	usersRepository := repousers.NewUsersPostgresRepo(s.db)

	//create services
	localizationService := applocalization.NewLocalizationServiceImpl()
	usersService := appusers.NewUsersServiceImpl(localizationService, usersRepository)
	authService := appauth.NewAuthService(authRepository, usersService, s.config.AuthConfig.AccessTokenLife, s.config.AuthConfig.RefreshTokenLife)

	//middlewares
	middlewares := middlewares.NewMiddlewares(authService)
	//create handlers
	authHandler := handlers.NewAuthHandler(authService, middlewares)
	usersHandler := handlers.NewUsersHandler(usersService, middlewares)

	//registering routes
	authHandler.RegisterRoutes(router)
	usersHandler.RegisterRoutes(router)
}