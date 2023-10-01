package usecases

import (
	"context"
	"github.com/Kreg101/metrics/internal/server/entity"
)

type Repository interface {
	Add(ctx context.Context, m entity.Metric)
	Get(ctx context.Context, name string) (entity.Metric, bool)
	GetAll(ctx context.Context) entity.Metrics
	Ping(ctx context.Context) error
}

type UseCases struct {
	repo Repository
}

func (uc *UseCases) AddMetric(ctx context.Context, metric entity.Metric) error {
	return nil
}

func (uc *UseCases) GetMetric(ctx context.Context, metric entity.Metric) (entity.Metric, error) {
	return entity.Metric{}, nil
}

// GetAllMetrics TODO: change Metrics type from map to slice
func (uc *UseCases) GetAllMetrics(ctx context.Context, metrics entity.Metrics) (entity.Metrics, error) {
	return nil, nil
}
func (uc *UseCases) Ping(ctx context.Context) error {
	return nil
}
