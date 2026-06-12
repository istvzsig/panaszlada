package storage

import (
	"context"

	"github.com/istvzsig/panaszlada/models"
)

type ReportStore interface {
	Create(ctx context.Context, r models.Report) error
	GetByCode(ctx context.Context, code string) (models.Report, error)
	List(ctx context.Context) ([]models.Report, error)
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
}
