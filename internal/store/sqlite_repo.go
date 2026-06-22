package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"github.com/google/uuid"
)

type sqliteRepo struct{ db *sql.DB }

func NewSQLiteRepo(db *sql.DB) (Repository, error) {
	if db == nil { return nil, fmt.Errorf("store: nil db") }
	return &sqliteRepo{db: db}, nil
}

func (r *sqliteRepo) CreateEndpoint(ctx context.Context, url string) (*Endpoint, error) {
	now := utcNow()
	e := &Endpoint{ID: url, URL: url, CreatedAt: now, UpdatedAt: now}
	_, err := r.db.ExecContext(ctx, `INSERT INTO endpoints (id, url, created_at, updated_at) VALUES (?, ?, ?, ?)`, e.ID, e.URL, e.CreatedAt, e.UpdatedAt)
	if err != nil { return nil, fmt.Errorf("store: create endpoint: %w", err) }
	return e, nil
}

func (r *sqliteRepo) GetEndpoint(ctx context.Context, id string) (*Endpoint, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, url, created_at, updated_at FROM endpoints WHERE id = ?`, id)
	e := &Endpoint{}
	if err := row.Scan(&e.ID, &e.URL, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return nil, fmt.Errorf("store: get endpoint: %w", err)
	}
	return e, nil
}

func (r *sqliteRepo) AppendRequest(ctx context.Context, endpointID string, req *Request) error {
	req.ID = uuid.NewString()
	req.EndpointID = endpointID
	req.CreatedAt = utcNow()
	_, err := r.db.ExecContext(ctx, `INSERT INTO requests (id, endpoint_id, method, path, headers, query, body, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.ID, req.EndpointID, req.Method, req.Path, req.Headers, req.Query, req.Body, req.CreatedAt)
	if err != nil { return fmt.Errorf("store: append request: %w", err) }
	return nil
}

func (r *sqliteRepo) ListRequests(ctx context.Context, endpointID string, limit int) ([]*Request, error) {
	if limit <= 0 { limit = 50 }
	rows, err := r.db.QueryContext(ctx, `SELECT id, endpoint_id, method, path, headers, query, body, created_at FROM requests WHERE endpoint_id = ? ORDER BY created_at DESC LIMIT ?`, endpointID, limit)
	if err != nil { return nil, fmt.Errorf("store: list requests: %w", err) }
	defer rows.Close()
	var out []*Request
	for rows.Next() {
		var req Request
		if err := rows.Scan(&req.ID, &req.EndpointID, &req.Method, &req.Path, &req.Headers, &req.Query, &req.Body, &req.CreatedAt); err != nil {
			return nil, fmt.Errorf("store: scan request: %w", err)
		}
		out = append(out, &req)
	}
	return out, rows.Err()
}

func (r *sqliteRepo) Close() error { return r.db.Close() }

func utcNow() string { return time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00") }
