package repousers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/users"
	"github.com/google/uuid"
)


var (
	errorCreatingUser = errors.New("there was an error creating the user")
	errorDeletingUser = errors.New("there was an error deleting the user")
	errorNoCodeAssociated = errors.New("the code is not associated with any temp couple")
	errorCheckingTempCouple = errors.New("there was an error checking the temp couple")
	errorGeneratingCode = errors.New("there was an error generating the couple's code")
)

type UsersPostgresRepo struct {
	db *sql.DB
}

func NewUsersPostgresRepo(db *sql.DB) users.UsersRepo{
	return &UsersPostgresRepo{
		db: db,
	}
}

func (r *UsersPostgresRepo) CreateUser(ctx context.Context, user *users.UserModel) error{
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO users(id, first_name, last_name, gender, birth_date, created_at, active, country_code, language_code, nickname)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		user.Id, user.FirstName, user.LastName, user.Gender, user.BirthDate, user.CreatedAt, user.Active, user.CountryCode, user.LanguageCode, user.NickName,
	)
	if err != nil {
		return errorCreatingUser
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return errorCreatingUser
	}
	return nil
}

func (r *UsersPostgresRepo) DeleteUser(ctx context.Context, userId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx, 
		`DELETE FROM users WHERE id = $1`,
		userId,
	)
	if err != nil{
		return errorDeletingUser
	}

	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return errorDeletingUser
	}
	return nil
}

func (r *UsersPostgresRepo) GetTempCoupleUserIdByCode(ctx context.Context, code int) (*uuid.UUID, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT user_id FROM temp_couples WHERE code = $1`,
		code,
	)
	userId := new(uuid.UUID) 
	if err := row.Scan(&userId); err != nil{
		return nil, errorNoCodeAssociated
	}
	return userId, nil
}

func (r *UsersPostgresRepo) CheckTempCoupleById(ctx context.Context, userId uuid.UUID) (bool, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id FROM temp_couples WHERE user_id = $1`,
		userId,
	)
	var result int 
	if err := row.Scan(&result); errors.Is(err, sql.ErrNoRows){
		return false, nil
	}else if err != nil{
		return false, errorCheckingTempCouple
	}
	return true, nil
}

func (r *UsersPostgresRepo) UpdateTempCouple(ctx context.Context, id uuid.UUID, code int) error{
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE temp_couples SET code = $1, updated_at = $2 WHERE user_id = $3`,
		code, time.Now(), id,
	)
	if err != nil{
		return errorGeneratingCode
	}
	if num, err := result.RowsAffected(); err != nil || num == 0{
		return errorGeneratingCode
	}
	return nil
}

func (r *UsersPostgresRepo) CreateTempCouple(ctx context.Context, id uuid.UUID, code int) error{
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO temp_couples (user_id, code, updated_at, created_at) VALUES($1, $2, $3, $4)`,
		id, code, time.Now(), time.Now(),
	)
	if err != nil{
		return errorGeneratingCode
	}
	if num, err := result.RowsAffected(); err != nil || num == 0{
		return errorGeneratingCode
	}
	return nil
}