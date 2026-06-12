package storage

import (
	"context"

	"github.com/istvzsig/panaszlada/models"
	"github.com/istvzsig/retryx"
)

type RetryStore struct {
	base ReportStore
	w    retryx.Wrapper
}

func NewRetryStore(base ReportStore, w retryx.Wrapper) *RetryStore {
	return &RetryStore{
		base: base,
		w:    w,
	}
}

func (s *RetryStore) Create(ctx context.Context, r models.Report) error {
	return s.w.Do(ctx, func(ctx context.Context) error {
		return s.base.Create(ctx, r)
	})
}

func (s *RetryStore) GetByCode(ctx context.Context, code string) (models.Report, error) {
	var report models.Report

	err := s.w.Do(ctx, func(ctx context.Context) error {
		r, err := s.base.GetByCode(ctx, code)
		if err != nil {
			return err
		}
		report = r
		return nil
	})

	return report, err
}

func (s *RetryStore) List(ctx context.Context) ([]models.Report, error) {
	var reports []models.Report

	err := s.w.Do(ctx, func(ctx context.Context) error {
		r, err := s.base.List(ctx)
		if err != nil {
			return err
		}
		reports = r
		return nil
	})

	return reports, err
}

func (s *RetryStore) UpdateStatus(ctx context.Context, id, status string) error {
	return s.w.Do(ctx, func(ctx context.Context) error {
		return s.base.UpdateStatus(ctx, id, status)
	})
}

func (s *RetryStore) Delete(ctx context.Context, id string) error {
	return s.w.Do(ctx, func(ctx context.Context) error {
		return s.base.Delete(ctx, id)
	})
}
