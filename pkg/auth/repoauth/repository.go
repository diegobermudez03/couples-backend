package repoauth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth/domainauth"
	"github.com/google/uuid"
)

var (
	errorCreatingUser = errors.New("there was an error creating the account for the user")
	errorCreatingSession = errors.New("there was an error creating the session for the user")
	errorRetrievingUserAuth = errors.New("there was an error retrieving the user auth")
)

type AuthPostgresRepo struct {
	db 		*sql.DB
}

func NewAuthPostgresRepo(db *sql.DB) domainauth.AuthRepository{
	return &AuthPostgresRepo{
		db: db,
	}
}

func (r *AuthPostgresRepo) CreateUserAuth(ctx context.Context, id uuid.UUID, email string, hash string) error{
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users_auth(id, email, hash, created_at) 
		VALUES ($1, $2, $3, $4)`,
		id, email, hash, time.Now(),
	)
	if err != nil{
		log.Print(err)
		return errorCreatingUser
	}
	if num, err := result.RowsAffected(); err != nil || num == 0{
		log.Print(err)
		return errorCreatingUser
	}
	return nil
}

func (r *AuthPostgresRepo) CreateSession(ctx context.Context,  id uuid.UUID, userId uuid.UUID, token string, device string, os string, expiresAt time.Time) error{
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO sessions(id, token, device, os, expires_at, created_at, last_used, user_auth_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		id, token, device, os, expiresAt, time.Now(), time.Now(), userId,
	)
	if err != nil{
		log.Print(err)
		return errorCreatingSession
	}
	if num, err := result.RowsAffected(); err != nil || num == 0{
		log.Print(err)
		return errorCreatingSession
	}
	return nil
}

func (r *AuthPostgresRepo) GetUserByEmail(ctx context.Context, email string) (*domainauth.UserAuthModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, hash, oauth_provider, oauth_id, created_at, user_id
		FROM users_auth WHERE email = $1`,
		email,
	)

	model := new(domainauth.UserAuthModel)
	*model = domainauth.UserAuthModel{}

	err := row.Scan(&model.Id, &model.Email, &model.Hash, &model.OauthProvider, &model.OauthId, &model.CreatedAt, &model.UserId)
	if err == sql.ErrNoRows{
		return nil, domainauth.ErrorNoUserFoundEmail
	}
	if err != nil{
		log.Print(err)
		return nil, errorRetrievingUserAuth
	}
	return model, nil
}