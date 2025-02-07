package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/diegobermudez03/couples-backend/internal/utils"
	"github.com/diegobermudez03/couples-backend/pkg/auth"
)

var (
	ErrNoAccessToken = errors.New("NO_ACCESS_TOKEN_PROVIDED")
)

type UserIdKey struct{}
type CoupleIdKey struct{}
type SessionIdKey struct{}



type Middlewares struct {
	authService auth.AuthService
	adminService auth.AuthAdminService
}

func NewMiddlewares(authService auth.AuthService, adminService auth.AuthAdminService) *Middlewares{
	return &Middlewares{
		authService: authService,
		adminService: adminService,
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