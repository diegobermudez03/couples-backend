package repousers

import (
	"context"
	"database/sql"
	"errors"
	"log"
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
	errorRettrievingCouple = errors.New("there was an error retrieving the couple")
	errorDeletingTempCouple = errors.New("there was an error deleting the temp couple")
	errorNoTempCoupleDeleted = errors.New("there was no temp couple to delete")
	errorNoUserToDelete = errors.New("there's no user to delete")
	errorCreatingCouple = errors.New("there was an error creating the couple")
	errorAddingPoints = errors.New("there was an error adding the points")
	errorNoUserFound = errors.New("no user found")
	errorUpdatingNickname = errors.New("there was an error updating the nickname")
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
	if num, _ := result.RowsAffected(); num == 0{
		return errorCreatingUser
	}
	return nil
}

func (r *UsersPostgresRepo) DeleteUserById(ctx context.Context, userId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx, 
		`DELETE FROM users WHERE id = $1`,
		userId,
	)
	if err != nil{
		return errorDeletingUser
	}

	if num, _ := result.RowsAffected(); num == 0{
		return errorNoUserToDelete
	}
	return nil
}

func (r *UsersPostgresRepo) GetTempCoupleByCode(ctx context.Context, code int) (*users.TempCoupleModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT user_id, code, start_date, created_at, updated_at
		FROM temp_couples WHERE code = $1`,
		code,
	)
	tempCouple := new(users.TempCoupleModel) 
	if err := row.Scan(&tempCouple.UserId, &tempCouple.Code, &tempCouple.StartDate, &tempCouple.CreatedAt, &tempCouple.UpdatedAt); err != nil{
		return nil, errorNoCodeAssociated
	}
	return tempCouple, nil
}

func (r *UsersPostgresRepo) CheckTempCoupleById(ctx context.Context, userId uuid.UUID) (bool, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT user_id FROM temp_couples WHERE user_id = $1`,
		userId,
	)
	var result uuid.UUID 
	if err := row.Scan(&result); errors.Is(err, sql.ErrNoRows){
		return false, nil
	}else if err != nil{
		return false, errorCheckingTempCouple
	}
	return true, nil
}

func (r *UsersPostgresRepo) UpdateTempCouple(ctx context.Context, tempCouple *users.TempCoupleModel) error{
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE temp_couples SET code = $1, start_date = $2, updated_at = $3 WHERE user_id = $4`,
		tempCouple.Code, tempCouple.StartDate ,time.Now(), tempCouple.UserId,
	)
	if err != nil{
		return errorGeneratingCode
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorGeneratingCode
	}
	return nil
}

func (r *UsersPostgresRepo) CreateTempCouple(ctx context.Context, tempCouple *users.TempCoupleModel) error{
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO temp_couples (user_id, code, start_date, updated_at, created_at) VALUES($1, $2, $3, $4, $5)`,
		tempCouple.UserId, tempCouple.Code, tempCouple.StartDate ,time.Now(), time.Now(),
	)
	if err != nil{
		return errorGeneratingCode
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorGeneratingCode
	}
	return nil
}

func (r *UsersPostgresRepo)  GetCoupleByUserId(ctx context.Context, userId uuid.UUID) (*users.CoupleModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, he_id, she_id, relation_start, end_date
		FROM couples WHERE he_id = $1 OR she_id = $1`,
		userId,
	)
	coupleModel := new(users.CoupleModel)

	err := row.Scan(&coupleModel.Id, &coupleModel.HeId, &coupleModel.SheId, &coupleModel.RelationStart, &coupleModel.EndDate)
	if errors.Is(err, sql.ErrNoRows){
		return nil, users.ErrorNoCoupleFound
	}else if err != nil{
		return nil, errorRettrievingCouple
	}
	return coupleModel, nil
}

func (r *UsersPostgresRepo)  DeleteTempCoupleById(ctx context.Context, id uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx, 
		`DELETE FROM temp_couples WHERE user_id = $1`,
		id,
	)
	if err != nil{
		return errorDeletingTempCouple
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorNoTempCoupleDeleted
	}
	return nil
}

func (r *UsersPostgresRepo) CreateCouple(ctx context.Context, couple *users.CoupleModel) error{
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO couples(id, he_id, she_id, relation_start)
		VALUES($1, $2, $3, $4)`,
		couple.Id, couple.HeId, couple.SheId, couple.RelationStart,
	)
	if err != nil{
		return errorCreatingCouple
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorCreatingCouple
	}
	return nil
}

func (r *UsersPostgresRepo) CreateCouplePoints(ctx context.Context, points *users.PointsModel) error{
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO points(id, points, day, couple_id)
		VALUES($1, $2, $3, $4)`,
		points.Id, points.Points, points.Day, *points.CoupleId,
	)
	if err != nil{
		log.Print(err)
		return errorAddingPoints
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorAddingPoints
	}
	return nil
}

func (r *UsersPostgresRepo) GetUserById(ctx context.Context, userId uuid.UUID) (*users.UserModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, first_name, last_name, gender, birth_date, created_at, 
		active, country_code, language_code, nickname
		FROM users WHERE id = $1`,
		userId,
	)
	user := new(users.UserModel)

	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Gender, &user.BirthDate, &user.CreatedAt, &user.Active, &user.CountryCode, &user.LanguageCode, &user.NickName)
	if err != nil{
		return nil, errorNoUserFound
	}
	return user, nil
}

func (r *UsersPostgresRepo) GetCoupleById(ctx context.Context, coupleId uuid.UUID) (*users.CoupleModel, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, he_id, she_id, relation_start, end_date
		FROM couples WHERE id = $1`,
		coupleId,
	)
	coupleModel := new(users.CoupleModel)

	err := row.Scan(&coupleModel.Id, &coupleModel.HeId, &coupleModel.SheId, &coupleModel.RelationStart, &coupleModel.EndDate)
	if errors.Is(err, sql.ErrNoRows){
		return nil, users.ErrorNoCoupleFound
	}else if err != nil{
		return nil, errorRettrievingCouple
	}
	return coupleModel, nil
}

func (r *UsersPostgresRepo) UpdateUserNicknameById(ctx context.Context, userId uuid.UUID, nickname string) error{
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET nickname = $1 WHERE id = $2`,
		nickname, userId,
	)
	if err != nil{
		return errorUpdatingNickname
	}
	if num, _ := result.RowsAffected(); num == 0{
		return errorUpdatingNickname
	}
	return nil
}