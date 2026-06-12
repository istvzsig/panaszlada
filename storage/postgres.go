package storage

import (
	"context"
	"database/sql"

	"github.com/istvzsig/panaszlada/models"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// =====================
// CREATE (INSERT)
// =====================
func (s *PostgresStore) Create(ctx context.Context, r models.Report) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO reports (
			id, tracking_code, title, description,
			category, latitude, longitude, status, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		r.ID, r.TrackingCode, r.Title, r.Description,
		r.Category, r.Latitude, r.Longitude, r.Status, r.CreatedAt,
	)
	return err
}

// =====================
// GET BY CODE
// =====================
func (s *PostgresStore) GetByCode(ctx context.Context, code string) (models.Report, error) {
	var r models.Report

	err := s.db.QueryRowContext(ctx, `
		SELECT id, tracking_code, title, description,
		       category, latitude, longitude, status, created_at
		FROM reports WHERE tracking_code=$1
	`, code).Scan(
		&r.ID, &r.TrackingCode, &r.Title, &r.Description,
		&r.Category, &r.Latitude, &r.Longitude,
		&r.Status, &r.CreatedAt,
	)

	return r, err
}

// =====================
// LIST
// =====================
func (s *PostgresStore) List(ctx context.Context) ([]models.Report, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tracking_code, title, description,
		       category, latitude, longitude, status, created_at
		FROM reports
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Report

	for rows.Next() {
		var r models.Report

		if err := rows.Scan(
			&r.ID, &r.TrackingCode, &r.Title, &r.Description,
			&r.Category, &r.Latitude, &r.Longitude,
			&r.Status, &r.CreatedAt,
		); err != nil {
			return nil, err
		}

		out = append(out, r)
	}

	return out, rows.Err()
}

// =====================
// UPDATE STATUS
// =====================
func (s *PostgresStore) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE reports SET status=$1 WHERE id=$2`,
		status,
		id,
	)
	return err
}

// =====================
// DELETE
// =====================
func (s *PostgresStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM reports WHERE id=$1`,
		id,
	)
	return err
}
