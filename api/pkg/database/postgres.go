package database

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"api/pkg/config"
	"api/pkg/logging"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresOptions struct {
	Name               string          `json:"name"`
	User               string          `json:"user"`
	Host               string          `json:"host"`
	Port               string          `json:"port"`
	SSLMode            string          `json:"ssl_mode"`
	ConnectionTimeout  int             `json:"connection_timeout,omitempty"`
	Password           string          `json:"-"`
	SSLCertPath        string          `json:"ssl_cert_path,omitempty"`
	SSLKeyPath         string          `json:"ssl_key_path,omitempty"`
	SSLRootCertPath    string          `json:"ssl_root_cert_path,omitempty"`
	PoolMinConnections string          `json:"pool_min_connections,omitempty"`
	PoolMaxConnections string          `json:"pool_max_connections,omitempty"`
	PoolMaxConnLife    config.Duration `json:"pool_max_conn_life,omitempty"`
	PoolMaxConnIdle    config.Duration `json:"pool_max_conn_idle,omitempty"`
	PoolHealthCheck    config.Duration `json:"pool_health_check,omitempty"`
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

func NewPostgresDB(ctx context.Context, opts *PostgresOptions) (*PostgresDB, error) {
	pgxConfig, err := pgxpool.ParseConfig(pgDSN(opts))
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

func pgDSN(opts *PostgresOptions) string {
	vals := pgValues(opts)
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

func setIfPositiveDuration(m map[string]string, key string, d config.Duration) {
	if d.Duration > 0 {
		m[key] = d.Duration.String()
	}
}

func pgValues(opts *PostgresOptions) map[string]string {
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", opts.Name)
	setIfNotEmpty(p, "user", opts.User)
	setIfNotEmpty(p, "host", opts.Host)
	setIfNotEmpty(p, "port", opts.Port)
	setIfNotEmpty(p, "sslmode", opts.SSLMode)
	setIfPositive(p, "connect_timeout", opts.ConnectionTimeout)
	setIfNotEmpty(p, "password", opts.Password)
	setIfNotEmpty(p, "sslcert", opts.SSLCertPath)
	setIfNotEmpty(p, "sslkey", opts.SSLKeyPath)
	setIfNotEmpty(p, "sslrootcert", opts.SSLRootCertPath)
	setIfNotEmpty(p, "pool_min_conns", opts.PoolMinConnections)
	setIfNotEmpty(p, "pool_max_conns", opts.PoolMaxConnections)
	setIfPositiveDuration(p, "pool_max_conn_lifetime", opts.PoolMaxConnLife)
	setIfPositiveDuration(p, "pool_max_conn_idle_time", opts.PoolMaxConnIdle)
	setIfPositiveDuration(p, "pool_health_check_period", opts.PoolHealthCheck)
	return p
}

func (o *PostgresOptions) ConnectionURL() string {
	if o == nil {
		return ""
	}

	host := o.Host
	if v := o.Port; v != "" {
		host = host + ":" + v
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   o.Name,
	}

	if o.User != "" || o.Password != "" {
		u.User = url.UserPassword(o.User, o.Password)
	}

	q := u.Query()
	if v := o.ConnectionTimeout; v > 0 {
		q.Add("connect_timeout", strconv.Itoa(v))
	}
	if v := o.SSLMode; v != "" {
		q.Add("sslmode", v)
	}
	if v := o.SSLCertPath; v != "" {
		q.Add("sslcert", v)
	}
	if v := o.SSLKeyPath; v != "" {
		q.Add("sslkey", v)
	}
	if v := o.SSLRootCertPath; v != "" {
		q.Add("sslrootcert", v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
