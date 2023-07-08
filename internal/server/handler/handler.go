package handler

import (
	"fmt"
	"github.com/Kreg101/metrics/internal/server/constants"
	"github.com/Kreg101/metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Mux struct {
	storage *storage.Storage
	router  chi.Router
}

func NewMux() *Mux {
	mux := &Mux{}
	mux.storage = storage.NewStorage()
	mux.router = chi.NewRouter()
	return mux
}

func mainPage(mux *Mux) {
	mux.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "text/html")
		w.Write([]byte(metricsToString(mux.storage.GetAll())))
	})
}

func metricPage(mux *Mux) {
	mux.router.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		if v, ok := mux.storage.Get(name); ok {
			if mux.storage.CheckType(name) == chi.URLParam(r, "type") {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("content-type", "text/html")
				w.Write([]byte(singleMetricToString(v)))
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})
}

func updateMetric(mux *Mux) {
	mux.router.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		if chi.URLParam(r, "name") == constants.EmptyString {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch chi.URLParam(r, "type") {
		case constants.GaugeType:
			fmt.Println(chi.URLParam(r, "type"), " ", chi.URLParam(r, "name"), " ", chi.URLParam(r, "value"))
			res, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
			mux.storage.Add(chi.URLParam(r, "name"), storage.Gauge(res))
		case constants.CounterType:
			res, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
			mux.storage.Add(chi.URLParam(r, "name"), storage.Counter(res))
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

func (mux *Mux) Apply() chi.Router {
	mainPage(mux)
	metricPage(mux)
	updateMetric(mux)
	return mux.router
}

func metricsToString(m storage.Metrics) string {
	list := make([]string, 0)
	for k, v := range m {
		var keyValue = k + ":"
		switch res := v.(type) {
		case storage.Gauge:
			keyValue += fmt.Sprintf("%.3f", res)
		case storage.Counter:
			keyValue += fmt.Sprintf("%d", res)
		}
		list = append(list, keyValue)
	}
	return strings.Join(list, ", ")
}

func singleMetricToString(v interface{}) string {
	switch x := v.(type) {
	case storage.Gauge:
		return float2String(float64(x))
	case storage.Counter:
		return fmt.Sprintf("%d", x)
	}
	return ""
}

func float2String(x float64) string {
	if math.Trunc(x) == x {
		return fmt.Sprintf("%.0f", x)
	}
	return strings.TrimRight(fmt.Sprintf("%.3f", x), "0")
}
