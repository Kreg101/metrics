package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Kreg101/metrics/internal/communication"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/Kreg101/metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Repository interface {
	Add(key string, value interface{})
	GetAll() storage.Metrics
	Get(name string) (interface{}, bool)
	CheckType(name string) string
}

type Mux struct {
	storage Repository
}

func NewMux(storage Repository) *Mux {
	mux := &Mux{}
	mux.storage = storage
	return mux
}

func withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// due to the use of aka singleton pattern this log will be the same
		// as log in main
		log := logger.Default()

		start := time.Now()

		responseData := &responseData{}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start).Milliseconds()

		log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
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
			w.Header().Set("content-type", "application/json")
			w.Write([]byte(singleMetricToString(v)))
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (mux *Mux) updateMetric(w http.ResponseWriter, r *http.Request) {
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

func (mux *Mux) updateMetricWithBody(w http.ResponseWriter, r *http.Request) {
	var m communication.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case "counter":
		if m.Delta == nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mux.storage.Add(m.ID, storage.Counter(*m.Delta))

		w.Header().Set("Content-Type", "application/json")

		e := json.NewEncoder(w).Encode(m)
		if e != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "gauge":
		if m.Value == nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mux.storage.Add(m.ID, storage.Gauge(*m.Value))

		if v, ok := mux.storage.Get(m.ID); ok {
			*m.Value = float64(v.(storage.Gauge))
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		e := json.NewEncoder(w).Encode(m)
		if e != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (mux *Mux) getMetric(w http.ResponseWriter, r *http.Request) {
	var m communication.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if v, ok := mux.storage.Get(m.ID); ok {
		if mux.storage.CheckType(m.ID) == m.MType {
			switch m.MType {
			case "counter":
				tmp := int64(v.(storage.Counter))
				m.Delta = &tmp
			case "gauge":
				tmp := float64(v.(storage.Gauge))
				m.Value = &tmp
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (mux *Mux) Router() chi.Router {
	router := chi.NewRouter()
	//router.Get("/", withLogging(mux.mainPage))
	//router.Get("/value/{type}/{name}", withLogging(mux.metricPage))
	//router.Post("/update/{type}/{name}/{value}", withLogging(mux.updateMetric))
	router.Post("/update/", withLogging(mux.updateMetricWithBody))
	router.Post("/value/", withLogging(mux.getMetric))
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
