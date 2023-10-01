package transport

import (
	"context"
	"github.com/Kreg101/metrics/internal/entity"
	"github.com/Kreg101/metrics/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Repository interface {
	Add(ctx context.Context, m entity.Metric)
	Get(ctx context.Context, name string) (entity.Metric, bool)
	GetAll(ctx context.Context) entity.Metrics
	Ping(ctx context.Context) error
}

// Mux - структура для соединение transport запроса и хранилища
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
	router.Get("/", logging(compression(mux.mainPage)))
	router.Get("/ping", mux.ping)
	router.Get("/value/{type}/{name}", logging(compression(mux.metricPage)))
	router.Post("/update/{type}/{name}/{value}", logging(compression(mux.updateMetric)))
	router.Post("/update/", mux.check(logging(compression(mux.updateMetricWithBody))))
	router.Post("/value/", mux.check(logging(compression(mux.getMetric))))
	router.Post("/updates/", mux.check(logging(compression(mux.updates))))
	return router
}
