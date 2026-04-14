package command

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

func NewPostgreSQLRepository(db *database.PostgresDB) interfaces.CommandRepository {
	return &PostgreSQLRepository{db: db}
}

func (r *PostgreSQLRepository) Create(ctx context.Context, cmd *models.Command) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			insert into commands (device_id, action, payload)
			values ($1, $2, $3)
			returning id, status, created_at, updated_at
		`
		var dbID int
		err := tx.QueryRow(ctx, query,
			cmd.DeviceID,
			cmd.Action,
			cmd.Payload,
		).Scan(&dbID, &cmd.Status, &cmd.CreatedAt, &cmd.UpdatedAt)

		cmd.ID = fmt.Sprintf("%d", dbID)
		return err
	})
}

func (r *PostgreSQLRepository) ListPending(ctx context.Context, deviceID string) ([]models.Command, error) {
	var commands []models.Command

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			select id, device_id, action, payload, status, created_at, updated_at
			from commands
			where device_id = $1 and status = 'pending'
			order by created_at asc
		`
		rows, err := tx.Query(ctx, query, deviceID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var c models.Command
			var dbID, dbDeviceID int
			if err := rows.Scan(&dbID, &dbDeviceID, &c.Action, &c.Payload, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
				return err
			}
			c.ID = fmt.Sprintf("%d", dbID)
			c.DeviceID = fmt.Sprintf("%d", dbDeviceID)
			commands = append(commands, c)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (r *PostgreSQLRepository) Acknowledge(ctx context.Context, id string) error {
	return r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			update commands
			set status = 'completed', updated_at = current_timestamp
			where id = $1 and status = 'pending'
		`
		result, err := tx.Exec(ctx, query, id)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("command not found or already completed")
		}
		return nil
	})
}
