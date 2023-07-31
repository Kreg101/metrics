package handler

import (
	"encoding/json"
	"github.com/Kreg101/metrics/internal/algorithms"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/Kreg101/metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Repository interface {
	Add(metric.Metric)
	GetAll() storage.Metrics
	Get(name string) (metric.Metric, bool)
}

type Mux struct {
	storage Repository
}

func NewMux(storage Repository) *Mux {
	mux := &Mux{}
	mux.storage = storage
	return mux
}

func usingLogger(h http.HandlerFunc) http.HandlerFunc {
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

func usingCompression(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			ow.Header().Set("Content-Encoding", "gzip")
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	}
}

func (mux *Mux) mainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(algorithms.Metrics2String(mux.storage.GetAll())))
}

func (mux *Mux) metricPage(w http.ResponseWriter, r *http.Request) {
	log := logger.Default()
	name := chi.URLParam(r, "name")
	if m, ok := mux.storage.Get(name); ok {
		if m.MType == chi.URLParam(r, "type") {
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(algorithms.SingleMetric2String(m)))
			return
		}
		log.Infof("mismatch metric type of %s in request and storage", name)
	}
	log.Infof("no metric %s in storage", name)
	w.WriteHeader(http.StatusNotFound)
}

func (mux *Mux) updateMetric(w http.ResponseWriter, r *http.Request) {
	log := logger.Default()
	if chi.URLParam(r, "name") == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Errorf("empty metric name")
		return
	}

	switch chi.URLParam(r, "type") {
	case "gauge":
		res, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("can't parse %s in float", chi.URLParam(r, "value"))
			return
		}
		w.WriteHeader(http.StatusOK)
		mux.storage.Add(metric.Metric{
			ID:    chi.URLParam(r, "name"),
			MType: "gauge",
			Value: &res,
			Delta: nil,
		})

	case "counter":
		res, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("can't parse %s in double", chi.URLParam(r, "value"))
			return
		}
		w.WriteHeader(http.StatusOK)
		mux.storage.Add(metric.Metric{
			ID:    chi.URLParam(r, "name"),
			MType: "counter",
			Value: nil,
			Delta: &res,
		})

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Infof("wrong metric type: %s", chi.URLParam(r, "type"))
	}
}

func (mux *Mux) updateMetricWithBody(w http.ResponseWriter, r *http.Request) {
	log := logger.Default()

	var m metric.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("can't unmarshal json %s to %v", r.Body, m)
		return
	}

	if (m.MType == "counter" && m.Delta == nil) || (m.MType == "gauge" && m.Value == nil) {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("empty delta or value in request")
		return
	}
	//switch m.MType {
	//case "counter":
	//	if m.Delta == nil {
	//
	//	}
	//case "gauge":
	//	if m.Value == nil {
	//		w.WriteHeader(http.StatusBadRequest)
	//		log.Errorf("empty gauge value in request")
	//		return
	//	}
	//}
	mux.storage.Add(m)

	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w).Encode(m)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("can't marshal %v", m)
		return
	}
}

func (mux *Mux) getMetric(w http.ResponseWriter, r *http.Request) {
	log := logger.Default()

	var m metric.Metric
	err := json.NewDecoder(r.Body).Decode(&m)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("can't unmarshal json %s to %v", r.Body, m)
		return
	}

	if v, ok := mux.storage.Get(m.ID); ok {
		if v.MType == m.MType {
			switch m.MType {
			case "counter":
				tmp := *v.Delta
				m.Delta = &tmp
			case "gauge":
				tmp := *v.Value
				m.Value = &tmp
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			log.Infof("wrong type %s of metric %s", m.MType, m.ID)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		log.Infof("no %s metric in storage", m.ID)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("can't marshal %v", m)
		return
	}
}

func (mux *Mux) Router() chi.Router {
	router := chi.NewRouter()
	router.Get("/", usingLogger(usingCompression(mux.mainPage)))
	router.Get("/value/{type}/{name}", usingLogger(usingCompression(mux.metricPage)))
	router.Post("/update/{type}/{name}/{value}", usingLogger(usingCompression(mux.updateMetric)))
	router.Post("/update/", usingLogger(usingCompression(mux.updateMetricWithBody)))
	router.Post("/value/", usingLogger(usingCompression(mux.getMetric)))
	return router
}
