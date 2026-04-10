package upload

import (
	"api/internal/api/interfaces"
	"api/internal/shared/models"
	"api/pkg/database"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

const (
	ImageStatusPending  = "pending"
	ImageStatusComplete = "complete"
	ImageStatusAssigned = "assigned"
	ImageStatusOrphaned = "orphaned"
	ImageStatusDeleted  = "deleted"
)

type PostgreSQLRepository struct {
	db *database.PostgresDB
}

func NewPostgreSQLRepository(db *database.PostgresDB) interfaces.UploadRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}

func (r *PostgreSQLRepository) Create(ctx context.Context, image *models.File) error {
	var dbID int

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			INSERT INTO uploads (user_id, resource_type, resource_id, status, filename, original_name, mime_type, size, url, expires_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`

		err := tx.QueryRow(ctx, query,
			image.UserID,
			image.ResourceType,
			image.ResourceID,
			image.Status,
			image.Filename,
			image.OriginalName,
			image.MimeType,
			image.Size,
			image.URL,
			image.ExpiresAt).Scan(&dbID)

		image.ID = fmt.Sprintf("%d", dbID)

		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgreSQLRepository) AssignToResource(ctx context.Context, imageIDs []string, resourceType, resourceID string) error {
	if len(imageIDs) == 0 {
		return nil
	}

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			UPDATE uploads
			SET resource_type = $1, resource_id = $2, status = $3, updated_at = NOW()
			WHERE id = ANY($4) and status = $5
		`

		_, err := tx.Exec(ctx, query, resourceType, resourceID, ImageStatusAssigned, imageIDs, ImageStatusPending)
		return err
	})

	return err
}

func (r *PostgreSQLRepository) HardDelete(ctx context.Context, id string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			DELETE FROM uploads
			WHERE id = $1
		`

		_, err := tx.Exec(ctx, query, id)
		return err
	})
}

func (r *PostgreSQLRepository) Complete(ctx context.Context, uploadKey string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			UPDATE uploads
			SET status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE filename = $2 AND status = $3
		`

		result, err := tx.Exec(ctx, query, ImageStatusComplete, uploadKey, ImageStatusPending)
		if err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			return fmt.Errorf("no pending upload found for key: %s", uploadKey)
		}

		return nil
	})
}

func (r *PostgreSQLRepository) Delete(ctx context.Context, id string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			UPDATE uploads
			SET status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`

		_, err := tx.Exec(ctx, query, ImageStatusDeleted, id)
		return err
	})
}

func (r *PostgreSQLRepository) HardDeleteByResource(ctx context.Context, resourceType, resourceID string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			DELETE FROM uploads
			WHERE resource_type = $1 AND resource_id = $2
		`

		_, err := tx.Exec(ctx, query, resourceType, resourceID)
		return err
	})
}

func (r *PostgreSQLRepository) DeleteByResource(ctx context.Context, resourceType, resourceID string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			UPDATE uploads
			SET status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE resource_type = $2 AND resource_id = $3
		`

		_, err := tx.Exec(ctx, query, ImageStatusDeleted, resourceType, resourceID)
		return err
	})
}

func (r *PostgreSQLRepository) DeleteByResourceAndFilenames(ctx context.Context, resourceType, resourceID string, filenames []string) error {
	if len(filenames) == 0 {
		return nil
	}

	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			UPDATE uploads
			SET status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE resource_type = $2 AND resource_id = $3 AND filename = ANY($4)
		`

		_, err := tx.Exec(ctx, query, ImageStatusDeleted, resourceType, resourceID, filenames)
		return err
	})
}

func (r *PostgreSQLRepository) GetByID(ctx context.Context, id string) (*models.File, error) {
	var image models.File
	var dbId int
	var dbUserID *int

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			SELECT id, user_id, resource_type, resource_id, status, filename, original_name, mime_type, size, url, sort_order, expires_at, created_at, updated_at
			FROM uploads
			WHERE id = $1
		`

		err := tx.QueryRow(ctx, query, id).Scan(
			&dbId,
			&dbUserID,
			&image.ResourceType,
			&image.ResourceID,
			&image.Status,
			&image.Filename,
			&image.OriginalName,
			&image.MimeType,
			&image.Size,
			&image.URL,
			&image.SortOrder,
			&image.ExpiresAt,
			&image.CreatedAt,
			&image.UpdatedAt)

		if err != nil {
			if err == pgx.ErrNoRows {
				return fmt.Errorf("image not found")
			}
			return err
		}

		image.ID = fmt.Sprintf("%d", dbId)
		if dbUserID == nil {
			image.UserID = ""
		} else {
			image.UserID = fmt.Sprintf("%d", *dbUserID)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (r *PostgreSQLRepository) GetByResource(ctx context.Context, resourceType, resourceID string, assignedOnly bool) ([]models.File, error) {
	var images []models.File

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			SELECT id, user_id, resource_type, resource_id, status, filename, original_name, mime_type, size, url, sort_order, expires_at, created_at, updated_at
			FROM uploads
			WHERE resource_type = $1 AND resource_id = $2
		`

		var rows pgx.Rows
		var err error

		if assignedOnly {
			query += " AND status = $3 ORDER BY sort_order ASC, created_at DESC"
			rows, err = tx.Query(ctx, query, resourceType, resourceID, ImageStatusAssigned)
		} else {
			query += " ORDER BY sort_order ASC, created_at DESC"
			rows, err = tx.Query(ctx, query, resourceType, resourceID)
		}

		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var image models.File
			var dbID int
			var dbUserID *int

			if err := rows.Scan(&dbID, &dbUserID, &image.ResourceType, &image.ResourceID,
				&image.Status, &image.Filename,
				&image.OriginalName, &image.MimeType, &image.Size,
				&image.URL, &image.SortOrder, &image.ExpiresAt, &image.CreatedAt,
				&image.UpdatedAt); err != nil {
				return err
			}

			image.ID = fmt.Sprintf("%d", dbID)
			if dbUserID == nil {
				image.UserID = ""
			} else {
				image.UserID = fmt.Sprintf("%d", *dbUserID)
			}
			images = append(images, image)

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return images, nil
}

func (r *PostgreSQLRepository) GetByFilename(ctx context.Context, filename string) (*models.File, error) {
	var image models.File
	var dbId int
	var dbUserID *int

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			SELECT id, user_id, resource_type, resource_id, status, filename, original_name, mime_type, size, url, sort_order, expires_at, created_at, updated_at
			FROM uploads
			WHERE filename = $1
		`

		err := tx.QueryRow(ctx, query, filename).Scan(
			&dbId,
			&dbUserID,
			&image.ResourceType,
			&image.ResourceID,
			&image.Status,
			&image.Filename,
			&image.OriginalName,
			&image.MimeType,
			&image.Size,
			&image.URL,
			&image.SortOrder,
			&image.ExpiresAt,
			&image.CreatedAt,
			&image.UpdatedAt)

		if err != nil {
			if err == pgx.ErrNoRows {
				return fmt.Errorf("image not found")
			}
			return err
		}

		image.ID = fmt.Sprintf("%d", dbId)
		if dbUserID == nil {
			image.UserID = ""
		} else {
			image.UserID = fmt.Sprintf("%d", *dbUserID)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (r *PostgreSQLRepository) GetExpiredPending(ctx context.Context) ([]models.File, error) {
	var images []models.File

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			SELECT id, user_id, resource_type, resource_id, status, filename, original_name, mime_type, size, url, sort_order, expires_at, created_at, updated_at
			FROM uploads
			WHERE status = $1 AND expires_at < NOW()
		`

		rows, err := tx.Query(ctx, query, ImageStatusPending)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var image models.File
			var dbId int
			var dbUserID *int

			if err := rows.Scan(&dbId, &dbUserID, &image.ResourceType, &image.ResourceID,
				&image.Status, &image.Filename,
				&image.OriginalName, &image.MimeType, &image.Size,
				&image.URL, &image.SortOrder, &image.ExpiresAt, &image.CreatedAt,
				&image.UpdatedAt); err != nil {
				return err
			}

			image.ID = fmt.Sprintf("%d", dbId)
			if dbUserID == nil {
				image.UserID = ""
			} else {
				image.UserID = fmt.Sprintf("%d", *dbUserID)
			}
			images = append(images, image)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return images, nil
}

func (r *PostgreSQLRepository) Assign(ctx context.Context, resourceType, resourceID string, filenames []string) error {
	if len(filenames) == 0 {
		return nil
	}

	var notFound []string

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		// Update each file individually to set the sort_order based on array position
		query := `
			UPDATE uploads
			SET resource_type = $1, resource_id = $2, status = $3, sort_order = $4, updated_at = NOW()
			WHERE filename = $5 AND status = $6
			RETURNING filename
		`

		found := make(map[string]bool)
		for i, filename := range filenames {
			var returnedFilename string
			err := tx.QueryRow(ctx, query, resourceType, resourceID, ImageStatusAssigned, i, filename, ImageStatusPending).Scan(&returnedFilename)
			if err == nil {
				found[returnedFilename] = true
			} else if err != pgx.ErrNoRows {
				return err
			}
		}

		for _, filename := range filenames {
			if !found[filename] {
				notFound = append(notFound, filename)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(notFound) > 0 {
		return fmt.Errorf("some images were not found while performing assignment: %v", notFound)
	}

	return nil
}

func (r *PostgreSQLRepository) UpdateSortOrder(ctx context.Context, resourceType, resourceID string, filenames []string) error {
	if len(filenames) == 0 {
		return nil
	}

	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		// Update sort_order for already-assigned images based on array position
		query := `
			UPDATE uploads
			SET sort_order = $1, updated_at = NOW()
			WHERE filename = $2 AND resource_type = $3 AND resource_id = $4 AND status = $5
		`

		for i, filename := range filenames {
			_, err := tx.Exec(ctx, query, i, filename, resourceType, resourceID, ImageStatusAssigned)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
