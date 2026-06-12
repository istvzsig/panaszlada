package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/istvzsig/panaszlada/models"
)

type MemoryStore struct {
	mu      sync.RWMutex
	reports map[string]models.Report
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		reports: make(map[string]models.Report),
	}
}

func (m *MemoryStore) Create(ctx context.Context, r models.Report) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reports[r.ID] = r
	return nil
}

func (m *MemoryStore) GetByCode(ctx context.Context, code string) (models.Report, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, r := range m.reports {
		if r.TrackingCode == code {
			return r, nil
		}
	}

	return models.Report{}, errors.New("not found")
}

func (m *MemoryStore) List(ctx context.Context) ([]models.Report, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]models.Report, 0, len(m.reports))
	for _, r := range m.reports {
		out = append(out, r)
	}
	return out, nil
}

func (m *MemoryStore) UpdateStatus(ctx context.Context, id, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.reports[id]
	if !ok {
		return errors.New("not found")
	}

	r.Status = models.ReportStatus(status)
	m.reports[id] = r
	return nil
}

func (m *MemoryStore) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.reports, id)
	return nil
}
