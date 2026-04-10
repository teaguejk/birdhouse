package device

import (
	"api/internal/api/interfaces"
	"api/internal/shared/models"
	"api/pkg/database"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type PostgreSQLRepository struct {
	db *database.PostgresDB
}

func NewPostgreSQLRepository(db *database.PostgresDB) interfaces.DeviceRepository {
	return &PostgreSQLRepository{db: db}
}

func (r *PostgreSQLRepository) Create(ctx context.Context, device *models.Device) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			INSERT INTO devices (name, api_key_hash, location, active)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at
		`
		var dbID int
		err := tx.QueryRow(ctx, query,
			device.Name,
			device.APIKeyHash,
			device.Location,
			device.Active,
		).Scan(&dbID, &device.CreatedAt, &device.UpdatedAt)

		device.ID = fmt.Sprintf("%d", dbID)
		return err
	})
}

func (r *PostgreSQLRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	var device models.Device
	var dbID int

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			SELECT id, name, api_key_hash, location, active, created_at, updated_at
			FROM devices
			WHERE id = $1
		`
		return tx.QueryRow(ctx, query, id).Scan(
			&dbID,
			&device.Name,
			&device.APIKeyHash,
			&device.Location,
			&device.Active,
			&device.CreatedAt,
			&device.UpdatedAt,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("device not found")
		}
		return nil, err
	}

	device.ID = fmt.Sprintf("%d", dbID)
	return &device, nil
}

func (r *PostgreSQLRepository) GetByAPIKeyHash(ctx context.Context, hash string) (*models.Device, error) {
	var device models.Device
	var dbID int

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			SELECT id, name, api_key_hash, location, active, created_at, updated_at
			FROM devices
			WHERE api_key_hash = $1
		`
		return tx.QueryRow(ctx, query, hash).Scan(
			&dbID,
			&device.Name,
			&device.APIKeyHash,
			&device.Location,
			&device.Active,
			&device.CreatedAt,
			&device.UpdatedAt,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	device.ID = fmt.Sprintf("%d", dbID)
	return &device, nil
}

func (r *PostgreSQLRepository) List(ctx context.Context) ([]models.Device, error) {
	var devices []models.Device

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			SELECT id, name, api_key_hash, location, active, created_at, updated_at
			FROM devices
			ORDER BY created_at DESC
		`
		rows, err := tx.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var d models.Device
			var dbID int
			if err := rows.Scan(&dbID, &d.Name, &d.APIKeyHash, &d.Location, &d.Active, &d.CreatedAt, &d.UpdatedAt); err != nil {
				return err
			}
			d.ID = fmt.Sprintf("%d", dbID)
			devices = append(devices, d)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *PostgreSQLRepository) Update(ctx context.Context, device *models.Device) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			UPDATE devices
			SET name = $1, api_key_hash = $2, location = $3, active = $4, updated_at = CURRENT_TIMESTAMP
			WHERE id = $5
			RETURNING updated_at
		`
		return tx.QueryRow(ctx, query,
			device.Name,
			device.APIKeyHash,
			device.Location,
			device.Active,
			device.ID,
		).Scan(&device.UpdatedAt)
	})
}

func (r *PostgreSQLRepository) Delete(ctx context.Context, id string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `DELETE FROM devices WHERE id = $1`
		result, err := tx.Exec(ctx, query, id)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("device not found")
		}
		return nil
	})
}
