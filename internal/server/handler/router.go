package handler

import (
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Repository interface {
	Add(metric.Metric)
	GetAll() metric.Metrics
	Get(name string) (metric.Metric, bool)
}

type Mux struct {
	storage Repository
	log     *zap.SugaredLogger
}

func NewMux(storage Repository, log *zap.SugaredLogger) *Mux {
	mux := &Mux{}
	mux.storage = storage

	if log == nil {
		mux.log = logger.Default()
	} else {
		mux.log = log
	}

	return mux
}

// Router настроит роутер хэндлерами из handlers
func (mux *Mux) Router() chi.Router {
	router := chi.NewRouter()
	router.Get("/", logging(compression(mux.mainPage)))
	router.Get("/value/{type}/{name}", logging(compression(mux.metricPage)))
	router.Post("/update/{type}/{name}/{value}", logging(compression(mux.updateMetric)))
	router.Post("/update/", logging(compression(mux.updateMetricWithBody)))
	router.Post("/value/", logging(compression(mux.getMetric)))
	return router
}
