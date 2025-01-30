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

type AuthPostgresRepo struct {
	db 		*sql.DB
}

func NewAuthPostgresRepo(db *sql.DB) auth.AuthRepository{
	return &AuthPostgresRepo{
		db: db,
	}
}

func (r *AuthPostgresRepo) CreateUserAuth(ctx context.Context, id uuid.UUID, email string, hash string) (int, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users_auth(id, email, hash, created_at) 
		VALUES ($1, $2, $3, $4)`,
		id, email, hash, time.Now(),
	)
	if err != nil{
		log.Print("error creating user: ", err.Error())
		return 0, err
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *AuthPostgresRepo) CreateSession(ctx context.Context,  sessionModel *auth.SessionModel) (int, error){
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO sessions(id, token, device, os, expires_at, created_at, last_used, user_auth_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		sessionModel.Id, sessionModel.Token, sessionModel.Device, sessionModel.Os, sessionModel.ExpiresAt, 
		sessionModel.CreatedAt, sessionModel.LastUsed, sessionModel.UserAuthId,
	)
	if err != nil{
		log.Print("error creating session: ", err.Error())
		return 0, err
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *AuthPostgresRepo) GetUserByEmail(ctx context.Context, email string) (*auth.UserAuthModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, hash, oauth_provider, oauth_id, created_at, user_id
		FROM users_auth WHERE email = $1`,
		email,
	)
	return r.readUser(row)
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
		return nil, nil
	}else if err != nil{
		log.Print("error getting session : ", err.Error())
		return nil, err
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
	return r.readUser(row)
}

func (r *AuthPostgresRepo) CreateEmptyUser(ctx context.Context, id uuid.UUID, userId uuid.UUID) (int, error) {
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO users_auth(id, user_id, created_at) VALUES($1, $2, $3)`,
		id, userId, time.Now(),
	)
	if err != nil{
		log.Print("error creating empty user: ", err.Error())
		return 0, err
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *AuthPostgresRepo) UpdateAuthUserId(ctx context.Context, authId uuid.UUID, userId uuid.UUID) (int, error){
	result, err := r.db.ExecContext(
		ctx, 
		`UPDATE users_auth SET user_id = $1 WHERE id = $2`,
		userId, authId,
	)
	if err != nil{
		log.Print("error vinculating account: ", err.Error())
		return 0, err
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *AuthPostgresRepo) DeleteSessionById(ctx context.Context, sessionId uuid.UUID) (int, error){
	result, err := r.db.ExecContext(
		ctx, 
		`DELETE FROM sessions WHERE id = $1`, 
		sessionId,
	)
	if err != nil{
		log.Print("error deleting session: ", err.Error())
		return 0, err
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *AuthPostgresRepo) DeleteUserAuthById(ctx context.Context, authId uuid.UUID) (int, error){
	result, err := r.db.ExecContext(
		ctx, 
		`DELETE FROM users_auth WHERE id = $1`, 
		authId,
	)
	if err != nil{
		log.Print("error deleting account: ", err.Error())
		return 0, err
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *AuthPostgresRepo) UpdateAuthUserById(ctx context.Context, authId uuid.UUID, authModel *auth.UserAuthModel) (int, error){
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users_auth SET email = $1, hash = $2, oauth_provider = $3, oauth_id = $4
		WHERE id = $5`,
		authModel.Email, authModel.Hash, authModel.OauthProvider, authModel.OauthId, authId,
	)
	if err != nil{
		log.Print("error updating account: ", err.Error())
		return 0, err
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func (r *AuthPostgresRepo) GetSessionById(ctx context.Context, id uuid.UUID) (*auth.SessionModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, token, device, os, expires_at, created_at, last_used, user_auth_id
		FROM sessions WHERE id = $1`,
		id,
	)
	model := new(auth.SessionModel)

	err := row.Scan(&model.Id, &model.Token, &model.Device, &model.Os, &model.ExpiresAt, &model.CreatedAt, &model.LastUsed, &model.UserAuthId)
	if errors.Is(err, sql.ErrNoRows){
		return nil, nil
	}else if err != nil{
		log.Print("error getting session: ", err.Error())
		return nil, err
	}
	return model, nil
}

func (r *AuthPostgresRepo)  UpdateSessionLastUsed(ctx context.Context, sessionId uuid.UUID, lastTime time.Time) (int, error){
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE sessions SET last_used = $1 WHERE id = $2`,
		lastTime, sessionId,
	)
	if err != nil{
		return 0,err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

///////////////////// HELPERS

func (r *AuthPostgresRepo) readUser(row *sql.Row) (*auth.UserAuthModel, error){
	model := new(auth.UserAuthModel)
	err := row.Scan(&model.Id, &model.Email, &model.Hash, &model.OauthProvider, &model.OauthId, &model.CreatedAt, &model.UserId)
	if err == sql.ErrNoRows{
		return nil, nil
	}
	if err != nil{
		log.Print("error reading user: ", err.Error())
		return nil, err
	}
	return model, nil
}
	
	