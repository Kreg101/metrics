package handler

import (
	"fmt"
	"github.com/Kreg101/metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Repository interface {
	Add(key string, value interface{})
	GetAll() storage.Metrics
	Get(name string) (interface{}, bool)
	CheckType(name string) string
}

type Mux struct {
	storage Repository
	sync.Mutex
}

func NewMux(storage Repository) *Mux {
	mux := &Mux{}
	mux.storage = storage
	return mux
}

func (mux *Mux) mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(metricsToString(mux.storage.GetAll())))
}

func (mux *Mux) metricPage(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if v, ok := mux.storage.Get(name); ok {
		if mux.storage.CheckType(name) == chi.URLParam(r, "type") {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(singleMetricToString(v)))
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (mux *Mux) updateMetric(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()
	if chi.URLParam(r, "name") == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch chi.URLParam(r, "type") {
	case "gauge":
		res, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		mux.storage.Add(chi.URLParam(r, "name"), storage.Gauge(res))
	case "counter":
		res, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		mux.storage.Add(chi.URLParam(r, "name"), storage.Counter(res))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (mux *Mux) Router() chi.Router {
	router := chi.NewRouter()
	router.Get("/", mux.mainPage)
	router.Get("/value/{type}/{name}", mux.metricPage)
	router.Post("/update/{type}/{name}/{value}", mux.updateMetric)
	return router
}

func metricsToString(m storage.Metrics) string {
	list := make([]string, 0)
	for k, v := range m {
		var keyValue = k + ":"
		switch res := v.(type) {
		case storage.Gauge:
			keyValue += float2String(float64(res))
		case storage.Counter:
			keyValue += fmt.Sprintf("%d", res)
		}
		list = append(list, keyValue)
	}
	return strings.Join(list, ", ")
}

func singleMetricToString(v interface{}) string {
	switch res := v.(type) {
	case storage.Gauge:
		return float2String(float64(res))
	case storage.Counter:
		return fmt.Sprintf("%d", res)
	}
	return ""
}

func float2String(v float64) string {
	if math.Trunc(v) == v {
		return fmt.Sprintf("%.0f", v)
	}
	return strings.TrimRight(fmt.Sprintf("%.3f", v), "0")
}
