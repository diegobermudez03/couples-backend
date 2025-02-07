package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/diegobermudez03/couples-backend/pkg/quizzes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	ErrNoAccessToken = errors.New("NO_ACCESS_TOKEN_PROVIDED")
	ErrUnathorized  = errors.New("USER_UNATHORIZED_TO_PERFORM_OPERATION")
)

const QUIZ_ID_URL_PARAM = "quizId"

type UserIdKey struct{}
type CoupleIdKey struct{}
type SessionIdKey struct{}



type Middlewares struct {
	authService auth.AuthService
	adminService auth.AuthAdminService
	quizzService quizzes.UserService
}

func NewMiddlewares(authService auth.AuthService, adminService auth.AuthAdminService, quizzService quizzes.UserService) *Middlewares{
	return &Middlewares{
		authService: authService,
		adminService: adminService,
		quizzService: quizzService,
	}
}

func (m *Middlewares) CheckAccessToken(handler http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request){
			//get header
			tokenString := r.Header.Get("Authorization")
			if tokenString == ""{
				utils.WriteError(w, http.StatusBadRequest, ErrNoAccessToken)
				return 
			}
			tokenString = strings.Split(tokenString, " ")[1]	

			//validate
			claims, err := m.authService.ValidateAccessToken(r.Context(), tokenString)
			if err != nil{
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}
			
			ctx := context.WithValue(r.Context(), UserIdKey{}, claims.UserId)
			ctx = context.WithValue(ctx, CoupleIdKey{}, claims.CoupleId)
			ctx = context.WithValue(ctx, SessionIdKey{}, claims.SessionId)
			r = r.WithContext(ctx)
			log.Printf("Succesfully validated %s", claims.UserId)
			handler.ServeHTTP(w, r)
		},
	)
}

func (m *Middlewares) CheckAdminAccessToken(handler http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request){
			//get header
			tokenString := r.Header.Get("Authorization")
			if tokenString == ""{
				utils.WriteError(w, http.StatusBadRequest, ErrNoAccessToken)
				return 
			}
			tokenString = strings.Split(tokenString, " ")[1]

			//validate
			claims, err := m.adminService.ValidateAccessToken(r.Context(), tokenString)
			if err != nil{
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}
			
			ctx := context.WithValue(r.Context(), SessionIdKey{}, claims.SessionId)
			r = r.WithContext(ctx)
			log.Printf("Succesfully validated %s", claims.SessionId)
			handler.ServeHTTP(w, r)
		},
	)
}


func (m *Middlewares) CheckUserQuizPermissions(handler http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			quizId := chi.URLParam(r, QUIZ_ID_URL_PARAM)
			quizParsed, err := uuid.Parse(quizId)
			if err != nil{
				utils.WriteError(w, http.StatusBadRequest, utils.ErrEmptyQuizId)
				return 
			}

			userId, _ := r.Context().Value(UserIdKey{}).(uuid.UUID)
			
			quiz, _ := m.quizzService.GetQuizById(r.Context(), quizParsed)
			if quiz == nil || quiz.CreatorId == nil || userId != *quiz.CreatorId{
				log.Print("user attempt to modify other quiz")
				utils.WriteError(w, http.StatusUnauthorized, ErrUnathorized)
				return 
			}
			handler.ServeHTTP(w, r)
		},
	)
}