package transport

import (
	"context"
	"github.com/Kreg101/metrics/internal/server/entity"
	"github.com/Kreg101/metrics/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Logic interface {
	AddMetric(ctx context.Context, metric entity.Metric) error
	GetMetric(ctx context.Context, metric entity.Metric) (entity.Metric, error)
	GetAllMetrics(ctx context.Context, metrics entity.Metrics) (entity.Metrics, error)
	Ping(ctx context.Context) error
}

type Repository interface {
	Add(ctx context.Context, m entity.Metric)
	Get(ctx context.Context, name string) (entity.Metric, bool)
	GetAll(ctx context.Context) entity.Metrics
	Ping(ctx context.Context) error
}

// Mux - структура для соединения запроса и хранилища
type Mux struct {
	storage Repository
	log     *zap.SugaredLogger
	key     string
}

func New(storage Repository, log *zap.SugaredLogger, key string) *Mux {
	mux := &Mux{storage: storage}

	if log == nil {
		mux.log = logger.Default()
	} else {
		mux.log = log
	}

	mux.key = key

	return mux
}

// Router настроит роутер хэндлерами из handlers
func (mux *Mux) Router() chi.Router {
	router := chi.NewRouter()

	router.Get("/", logging(compression(mux.getMetrics)))
	router.Get("/ping", mux.ping)
	router.Get("/value/{type}/{name}", logging(compression(mux.getMetric1)))
	router.Post("/update/{type}/{name}/{value}", logging(compression(mux.updateMetric1)))
	router.Post("/update/", mux.check(logging(compression(mux.updateMetric2))))
	router.Post("/value/", mux.check(logging(compression(mux.getMetric2))))
	router.Post("/updates/", mux.check(logging(compression(mux.updateMetrics))))

	return router
}
