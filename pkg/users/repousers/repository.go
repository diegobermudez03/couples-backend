package repousers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diegobermudez03/couples-backend/pkg/users"
)


var (
	errorCreatingUser = errors.New("there was an error creating the user")
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
		`INSERT INTO users(id, first_name, last_name, gender, birth_date, created_at, active, country_code, language_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.Id, user.FirstName, user.LastName, user.Gender, user.BirthDate, user.CreatedAt, user.Active, user.CountryCode, user.LanguageCode,
	)
	if err != nil {
		return errorCreatingUser
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return errorCreatingUser
	}
	return nil
}