package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MikeCHOKKI/nexusflow/service-user/internal/model"
)

// Postgres gère la connexion et les opérations PostgreSQL.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres initialise le pool de connexions PostgreSQL.
func NewPostgres(databaseURL string) (*Postgres, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &Postgres{pool: pool}, nil
}

// Close ferme le pool de connexions.
func (p *Postgres) Close() {
	p.pool.Close()
}

// CreateUser insère un nouvel utilisateur et retourne l'enregistrement complet.
func (p *Postgres) CreateUser(ctx context.Context, email, password, name string) (*model.User, error) {
	u := &model.User{}
	err := p.pool.QueryRow(ctx,
		`INSERT INTO users.users (email, password, name)
		 VALUES ($1, $2, $3)
		 RETURNING id::text, email, password, name, role, is_active,
		           to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		           to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`,
		email, password, name,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Role, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

// GetUser récupère un utilisateur par son ID.
func (p *Postgres) GetUser(ctx context.Context, id string) (*model.User, error) {
	u := &model.User{}
	err := p.pool.QueryRow(ctx,
		`SELECT id::text, email, password, name, role, is_active,
		        to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		        to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		 FROM users.users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Role, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

// GetUserByEmail récupère un utilisateur par son email.
func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	err := p.pool.QueryRow(ctx,
		`SELECT id::text, email, password, name, role, is_active,
		        to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		        to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		 FROM users.users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Role, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

// UpdateUser met à jour les champs fournis pour un utilisateur.
func (p *Postgres) UpdateUser(ctx context.Context, id string, fields map[string]any) (*model.User, error) {
	if len(fields) == 0 {
		return p.GetUser(ctx, id)
	}

	// Construire le SET dynamiquement.
	setClauses := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields)+1)
	idx := 1
	for col, val := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, idx))
		args = append(args, val)
		idx++
	}
	args = append(args, id)

	query := fmt.Sprintf(
		`UPDATE users.users SET %s, updated_at = NOW()
		 WHERE id = $%d
		 RETURNING id::text, email, password, name, role, is_active,
		           to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		           to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`,
		strings.Join(setClauses, ", "),
		idx,
	)

	u := &model.User{}
	err := p.pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.Password, &u.Name, &u.Role, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

// ListUsers retourne une liste paginée d'utilisateurs, filtrée optionnellement par rôle.
func (p *Postgres) ListUsers(ctx context.Context, page, limit int32, role string) ([]*model.User, int32, error) {
	// Compter le total.
	var total int32
	countQuery := "SELECT COUNT(*) FROM users.users"
	countArgs := make([]any, 0)

	if role != "" {
		countQuery += " WHERE role = $1"
		countArgs = append(countArgs, role)
	}
	if err := p.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Pagination.
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := `SELECT id::text, email, password, name, role, is_active,
	                 to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
	                 to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	          FROM users.users`
	queryArgs := make([]any, 0)
	argIdx := 1

	if role != "" {
		query += fmt.Sprintf(" WHERE role = $%d", argIdx)
		queryArgs = append(queryArgs, role)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := p.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Role, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, total, nil
}

// ─── Helpers ──────────────────────────────────────────────────

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique")
}

func isNoRows(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no rows in result set")
}
