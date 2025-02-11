package repofiles

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diegobermudez03/couples-backend/pkg/files"
	"github.com/diegobermudez03/couples-backend/pkg/infraestructure"
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
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx, 
			`INSERT INTO files(id, bucket, grouping, object_key, url, public, type)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			file.Id, file.Bucket, file.Group, file.ObjectKey, file.Url, file.Public, file.Type,
		)
	})
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

func (r *FilesPostgresRepo) DeleteFileById(ctx context.Context, id uuid.UUID) (int, error){
	return infraestructure.ExecSQL(ctx, r.db, func(ex infraestructure.Executor) (sql.Result, error) {
		return ex.ExecContext(
			ctx,
			`DELETE FROM files WHERE id = $1`,
			id,  
		)
	})
}