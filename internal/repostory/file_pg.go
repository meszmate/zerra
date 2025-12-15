package repostory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
)

type FileRepostory interface {
	GetFilesByParent(ctx context.Context, parentType models.FileParentType, parentID string) ([]models.File, *errx.Error)
	PutFile(ctx context.Context, tx pgx.Tx, fileKey string, fileType string, parentType models.FileParentType, parentID, name string) (*models.File, *errx.Error)
	DeleteFile(ctx context.Context, tx pgx.Tx, ID string) (*models.File, *errx.Error)
}

type fileRepostory struct {
	DB *db.DB
}

func NewFileRepostory(db *db.DB) FileRepostory {
	return &fileRepostory{
		DB: db,
	}
}

func (r *fileRepostory) GetFilesByParent(ctx context.Context, parentType models.FileParentType, parentID string) ([]models.File, *errx.Error) {
	query := `
		SELECT id, parent_type, parent_id,
		 name, file_type, created_at
		FROM files
		WHERE owner_type = $1
		 AND owner_id = $2
	`

	params := []any{
		parentType,
		parentID,
	}

	var files []models.File = make([]models.File, 0)

	rows, err := r.DB.Query(
		ctx,
		query,
		params...,
	)
	if err != nil {
		db.CaptureError(err, query, params, "query")
		return nil, errx.InternalError()
	}

	for rows.Next() {
		var f models.File
		if err := rows.Scan(
			&f.ID, &f.ParentType, &f.ParentID,
			&f.Name, &f.FileType, &f.CreatedAt,
		); err != nil {
			db.CaptureError(err, "", nil, "")
			return nil, errx.InternalError()
		}
		files = append(files, f)
	}

	return files, nil
}

func (r *fileRepostory) PutFile(ctx context.Context, tx pgx.Tx, fileKey, fileType string, parentType models.FileParentType, parentID, name string) (*models.File, *errx.Error) {
	file := &models.File{
		ID:         uuid.NewString(),
		ParentType: parentType,
		ParentID:   parentID,
		Name:       name,
		FileType:   fileType,
		CreatedAt:  time.Now(),
	}

	query := `
		INSERT INTO files (
		 id, parent_type, parent_id,
		 name, file_type, created_at
		)
	`

	params := []any{
		file.ID, string(parentType), parentID,
		name, fileType, file.CreatedAt,
	}

	if _, err := tx.Exec(ctx, query, params...); err != nil {
		db.CaptureError(err, query, params, "exec")
		return nil, errx.InternalError()
	}

	return file, nil
}

func (r *fileRepostory) DeleteFile(ctx context.Context, tx pgx.Tx, ID string) (*models.File, *errx.Error) {
	var file models.File

	query := `
		DELETE FROM files
		WHERE id = $1
		RETURNING id, parent_type, parent_id,
		 name, file_type, created_at
	`

	params := []any{
		ID,
	}

	if err := tx.QueryRow(ctx, query, params...).Scan(
		&file.ID, &file.ParentType, &file.ParentID,
		&file.Name, &file.FileType, &file.CreatedAt,
	); err != nil {
		db.CaptureError(err, query, params, "exec")
		return nil, errx.InternalError()
	}

	return &file, nil
}
