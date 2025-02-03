package repofiles

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/google/uuid"
)

type FilesPostgresRepo struct {
	db *sql.DB
}

func NewFilesPostgresRepo(db *sql.DB) files.Repository{
	return &FilesPostgresRepo{
		db: db,
	}
}


func (r *FilesPostgresRepo) CreateFile(ctx context.Context, file *files.FileModel) (int,error){
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO files(id, bucket, grouping, object_key, url, public, type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		file.Id, file.Bucket, file.Group, file.ObjectKey, file.Url, file.Public, file.Type,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil 
}

func (r *FilesPostgresRepo) GetFileById(ctx context.Context, id uuid.UUID) (*files.FileModel, error){
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT id, bucket, grouping, object_key, url, public, type
		FROM files WHERE id = $1`,
		id,
	)
	model := new(files.FileModel)
	err := row.Scan(&model.Id, &model.Bucket, &model.Group, &model.ObjectKey, &model.Url, &model.Public, &model.Type)
	if errors.Is(err, sql.ErrNoRows){
		return nil, nil
	}
	if err != nil{
		return nil, err 
	}
	return model, nil
}