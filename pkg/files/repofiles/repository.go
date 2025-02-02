package repofiles

import (
	"context"
	"database/sql"

	"github.com/diegobermudez03/couples-backend/pkg/files"
)

type FilesPostgresRepo struct {
	db *sql.DB
}

func NewFilesPostgresRepo(db *sql.DB) files.Repository{
	return &FilesPostgresRepo{
		db: db,
	}
}


func (r *FilesPostgresRepo)CreateFile(ctx context.Context, file *files.FileModel) (int,error){
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO files(id, bucket, grouping, object_key, created_at, type)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		file.Id, file.Bucket, file.Group, file.ObjectKey, file.CreatedAt, file.Type,
	)
	if err != nil{
		return 0, err 
	}
	num, _ := result.RowsAffected()
	return int(num), nil 
}