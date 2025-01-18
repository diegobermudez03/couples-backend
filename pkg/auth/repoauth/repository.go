package repoauth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth"
	"github.com/google/uuid"
)

var (
	errorCreatingUser = errors.New("there was an error creating the account for the user")
	errorCreatingSession = errors.New("there was an error creating the session for the user")
	errorRetrievingUserAuth = errors.New("there was an error retrieving the user auth")
	errorNoSessionFound = errors.New("no session found")
	errorRetrievingSession = errors.New("there was an error retrieving the session")
	errorNoUserFoundId = errors.New("no user found with the Id")
	errorVinculatingAccount = errors.New("there was an error vinculating the account")
	errorDeletingSession = errors.New("there was an error deleting the session")
	errorDeletingUserAuth = errors.New("there was an error deleting the user auth")
	errorsNoSessionToDelete = errors.New("there's no session to delete")
	errorNoUserAuthToDelete = errors.New("there's no user auth to delete")
	errorUpdatingUserAuth = errors.New("there was an error updating the user account")
)

type AuthPostgresRepo struct {
	db 		*sql.DB
}

func NewAuthPostgresRepo(db *sql.DB) auth.AuthRepository{
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
		return errorCreatingUser
	}
	if num, _ := result.RowsAffected();num == 0{
		return errorCreatingUser
	}
	return nil
}

func (r *AuthPostgresRepo) CreateSession(ctx context.Context,  id uuid.UUID, userId uuid.UUID, token string, device *string, os *string, expiresAt time.Time) error{
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO sessions(id, token, device, os, expires_at, created_at, last_used, user_auth_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		id, token, device, os, expiresAt, time.Now(), time.Now(), userId,
	)
	if err != nil{
		return errorCreatingSession
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorCreatingSession
	}
	return nil
}

func (r *AuthPostgresRepo) GetUserByEmail(ctx context.Context, email string) (*auth.UserAuthModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, hash, oauth_provider, oauth_id, created_at, user_id
		FROM users_auth WHERE email = $1`,
		email,
	)
	return r.readUser(row, auth.ErrorNoUserFoundEmail)
}

func (r *AuthPostgresRepo) GetSessionByToken(ctx context.Context, token string) (*auth.SessionModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, token, device, os, expires_at, created_at, last_used, user_auth_id
		FROM sessions WHERE token = $1`,
		token,
	)
	model := new(auth.SessionModel)

	err := row.Scan(&model.Id, &model.Token, &model.Device, &model.Os, &model.ExpiresAt, &model.CreatedAt, &model.LastUsed, &model.UserAuthId)
	if errors.Is(err, sql.ErrNoRows){
		return nil, errorNoSessionFound
	}else if err != nil{
		return nil, errorRetrievingSession
	}

	return model, nil
}


func (r *AuthPostgresRepo) GetUserById(ctx context.Context, id uuid.UUID) (*auth.UserAuthModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, hash, oauth_provider, oauth_id, created_at, user_id
		FROM users_auth WHERE id = $1`,
		id,
	)
	return r.readUser(row, errorNoUserFoundId)
}

func (r *AuthPostgresRepo) CreateEmptyUser(ctx context.Context, id uuid.UUID, userId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO users_auth(id, user_id, created_at) VALUES($1, $2, $3)`,
		id, userId, time.Now(),
	)
	if err != nil{
		return errorCreatingUser
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorCreatingUser
	}
	return nil
}

func (r *AuthPostgresRepo) UpdateAuthUserId(ctx context.Context, authId uuid.UUID, userId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx, 
		`UPDATE users_auth SET user_id = $1 WHERE id = $2`,
		userId, authId,
	)
	if err != nil{
		log.Print(err)
		return errorVinculatingAccount
	}
	if num, _ := result.RowsAffected(); num != 1{
		return errorVinculatingAccount
	}
	return nil
}

func (r *AuthPostgresRepo) DeleteSessionById(ctx context.Context, sessionId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx, 
		`DELETE FROM sessions WHERE id = $1`, 
		sessionId,
	)
	if err != nil{
		return errorDeletingSession
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorsNoSessionToDelete
	}
	return nil
}

func (r *AuthPostgresRepo) DeleteUserAuthById(ctx context.Context, authId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx, 
		`DELETE FROM users_auth WHERE id = $1`, 
		authId,
	)
	if err != nil{
		return errorDeletingUserAuth
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorNoUserAuthToDelete
	}
	return nil
}

func (r *AuthPostgresRepo) UpdateAuthUserById(ctx context.Context, authId uuid.UUID, authModel *auth.UserAuthModel) error{
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users_auth SET email = $1, hash = $2, oauth_provider = $3, oauth_id = $4
		WHERE id = $5`,
		authModel.Email, authModel.Hash, authModel.OauthProvider, authModel.OauthId, authId,
	)
	if err != nil{
		return errorUpdatingUserAuth
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorUpdatingUserAuth
	}
	return nil
}


///////////////////// HELPERS

func (r *AuthPostgresRepo) readUser(row *sql.Row, caseError error) (*auth.UserAuthModel, error){
	model := new(auth.UserAuthModel)
	err := row.Scan(&model.Id, &model.Email, &model.Hash, &model.OauthProvider, &model.OauthId, &model.CreatedAt, &model.UserId)
	if err == sql.ErrNoRows{
		return nil, caseError
	}
	if err != nil{
		log.Print(err)
		return nil, errorRetrievingUserAuth
	}
	return model, nil
}
	
	