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
			insert into devices (name, api_key_hash, location, active)
			values ($1, $2, $3, $4)
			returning id, config, created_at, updated_at
		`
		var dbID int
		err := tx.QueryRow(ctx, query,
			device.Name,
			device.APIKeyHash,
			device.Location,
			device.Active,
		).Scan(&dbID, &device.Config, &device.CreatedAt, &device.UpdatedAt)

		device.ID = fmt.Sprintf("%d", dbID)
		return err
	})
}

func (r *PostgreSQLRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	var device models.Device
	var dbID int

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			select id, name, api_key_hash, location, active, config,
			       last_seen_at, last_status, created_at, updated_at
			from devices
			where id = $1
		`
		return tx.QueryRow(ctx, query, id).Scan(
			&dbID,
			&device.Name,
			&device.APIKeyHash,
			&device.Location,
			&device.Active,
			&device.Config,
			&device.LastSeenAt,
			&device.LastStatus,
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
			select id, name, api_key_hash, location, active, config,
			       last_seen_at, last_status, created_at, updated_at
			from devices
			where api_key_hash = $1
		`
		return tx.QueryRow(ctx, query, hash).Scan(
			&dbID,
			&device.Name,
			&device.APIKeyHash,
			&device.Location,
			&device.Active,
			&device.Config,
			&device.LastSeenAt,
			&device.LastStatus,
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
			select id, name, api_key_hash, location, active, config,
			       last_seen_at, last_status, created_at, updated_at
			from devices
			order by created_at desc
		`
		rows, err := tx.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var d models.Device
			var dbID int
			if err := rows.Scan(
				&dbID, &d.Name, &d.APIKeyHash, &d.Location, &d.Active, &d.Config,
				&d.LastSeenAt, &d.LastStatus, &d.CreatedAt, &d.UpdatedAt,
			); err != nil {
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
			update devices
			set name = $1, api_key_hash = $2, location = $3, active = $4, config = $5,
			    updated_at = current_timestamp
			where id = $6
			returning updated_at
		`
		return tx.QueryRow(ctx, query,
			device.Name,
			device.APIKeyHash,
			device.Location,
			device.Active,
			device.Config,
			device.ID,
		).Scan(&device.UpdatedAt)
	})
}

func (r *PostgreSQLRepository) TouchLastSeen(ctx context.Context, id string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `update devices set last_seen_at = current_timestamp where id = $1`
		_, err := tx.Exec(ctx, query, id)
		return err
	})
}

func (r *PostgreSQLRepository) UpdateStatus(ctx context.Context, id string, status []byte) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			update devices
			set last_status = $1, last_seen_at = current_timestamp
			where id = $2
		`
		_, err := tx.Exec(ctx, query, status, id)
		return err
	})
}

func (r *PostgreSQLRepository) ListStatus(ctx context.Context) ([]models.DeviceStatus, error) {
	var statuses []models.DeviceStatus

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			select id, name, location, active, config, last_seen_at, last_status
			from devices
			order by name asc
		`
		rows, err := tx.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var s models.DeviceStatus
			var dbID int
			if err := rows.Scan(
				&dbID, &s.Name, &s.Location, &s.Active, &s.Config,
				&s.LastSeenAt, &s.LastStatus,
			); err != nil {
				return err
			}
			s.ID = fmt.Sprintf("%d", dbID)
			statuses = append(statuses, s)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (r *PostgreSQLRepository) Delete(ctx context.Context, id string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `delete from devices where id = $1`
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
