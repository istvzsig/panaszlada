package testutil

import (
	"context"

	"github.com/istvzsig/panaszlada/models"
)

type MockStore struct {
	CreateFn       func(ctx context.Context, r models.Report) error
	GetByCodeFn    func(ctx context.Context, code string) (models.Report, error)
	ListFn         func(ctx context.Context) ([]models.Report, error)
	UpdateStatusFn func(ctx context.Context, id string, status string) error
	DeleteFn       func(ctx context.Context, id string) error
}

func NewMockStore() *MockStore {
	return &MockStore{
		CreateFn: func(ctx context.Context, r models.Report) error { return nil },
		GetByCodeFn: func(ctx context.Context, code string) (models.Report, error) {
			return models.Report{}, nil
		},
		ListFn: func(ctx context.Context) ([]models.Report, error) {
			return []models.Report{}, nil
		},
		UpdateStatusFn: func(ctx context.Context, id string, status string) error {
			return nil
		},
		DeleteFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
}

func (m *MockStore) Create(ctx context.Context, r models.Report) error {
	return m.CreateFn(ctx, r)
}

func (m *MockStore) GetByCode(ctx context.Context, code string) (models.Report, error) {
	return m.GetByCodeFn(ctx, code)
}

func (m *MockStore) List(ctx context.Context) ([]models.Report, error) {
	return m.ListFn(ctx)
}

func (m *MockStore) UpdateStatus(ctx context.Context, id string, status string) error {
	return m.UpdateStatusFn(ctx, id, status)
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
