package handler

import (
	"context"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Repository interface {
	Add(ctx context.Context, m metric.Metric)
	Get(ctx context.Context, name string) (metric.Metric, bool)
	GetAll(ctx context.Context) metric.Metrics
	Ping(ctx context.Context) error
}

// Mux - структура для соединение http запроса и хранилища
type Mux struct {
	storage Repository
	log     *zap.SugaredLogger
	key     string
}

func NewMux(storage Repository, log *zap.SugaredLogger, key string) *Mux {
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
	router.Get("/", logging(compression(mux.mainPage)))
	router.Get("/ping", mux.ping)
	router.Get("/value/{type}/{name}", logging(compression(mux.metricPage)))
	router.Post("/update/{type}/{name}/{value}", logging(compression(mux.updateMetric)))
	router.Post("/update/", mux.check(logging(compression(mux.updateMetricWithBody))))
	router.Post("/value/", mux.check(logging(compression(mux.getMetric))))
	router.Post("/updates/", mux.check(logging(compression(mux.updates))))
	return router
}
