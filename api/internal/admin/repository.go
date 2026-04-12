package admin

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

func NewPostgreSQLRepository(db *database.PostgresDB) interfaces.AdminRepository {
	return &PostgreSQLRepository{db: db}
}

func (r *PostgreSQLRepository) GetByEmail(ctx context.Context, email string) (*models.AdminUser, error) {
	var user models.AdminUser
	var dbID int

	err := r.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		query := `
			SELECT id, email, name, active, created_at, updated_at
			FROM admins
			WHERE email = $1 AND active = true
		`
		return tx.QueryRow(ctx, query, email).Scan(
			&dbID,
			&user.Email,
			&user.Name,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.ID = fmt.Sprintf("%d", dbID)
	return &user, nil
}
