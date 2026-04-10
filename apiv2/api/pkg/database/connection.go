package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"api/pkg/logging"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context)
	GetDB() interface{}
}

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func (db *PostgresDB) Close(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Info("closing postgres connection pool.")
	db.Pool.Close()
}

func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func (db *PostgresDB) GetDB() interface{} {
	return db.Pool
}

func NewPostgresDB(ctx context.Context, cfg *Config) (*PostgresDB, error) {
	pgxConfig, err := pgxpool.ParseConfig(dbDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pgxConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

func dbDSN(cfg *Config) string {
	vals := dbValues(cfg)
	p := make([]string, 0, len(vals))
	for k, v := range vals {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(p, " ")
}

func setIfNotEmpty(m map[string]string, key, val string) {
	if val != "" {
		m[key] = val
	}
}

func setIfPositive(m map[string]string, key string, val int) {
	if val > 0 {
		m[key] = fmt.Sprintf("%d", val)
	}
}

func setIfPositiveDuration(m map[string]string, key string, d time.Duration) {
	if d > 0 {
		m[key] = d.String()
	}
}

func dbValues(cfg *Config) map[string]string {
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", cfg.Name)
	setIfNotEmpty(p, "user", cfg.User)
	setIfNotEmpty(p, "host", cfg.Host)
	setIfNotEmpty(p, "port", cfg.Port)
	setIfNotEmpty(p, "sslmode", cfg.SSLMode)
	setIfPositive(p, "connect_timeout", cfg.ConnectionTimeout)
	setIfNotEmpty(p, "password", cfg.Password)
	setIfNotEmpty(p, "sslcert", cfg.SSLCertPath)
	setIfNotEmpty(p, "sslkey", cfg.SSLKeyPath)
	setIfNotEmpty(p, "sslrootcert", cfg.SSLRootCertPath)
	setIfNotEmpty(p, "pool_min_conns", cfg.PoolMinConnections)
	setIfNotEmpty(p, "pool_max_conns", cfg.PoolMaxConnections)
	setIfPositiveDuration(p, "pool_max_conn_lifetime", cfg.PoolMaxConnLife)
	setIfPositiveDuration(p, "pool_max_conn_idle_time", cfg.PoolMaxConnIdle)
	setIfPositiveDuration(p, "pool_health_check_period", cfg.PoolHealthCheck)
	return p
}

func NewPostgresDBWithName(ctx context.Context, cfg *Config, dbName string) (*PostgresDB, error) {
	cfgCopy := *cfg
	cfgCopy.Name = dbName

	return NewPostgresDB(ctx, &cfgCopy)
}
