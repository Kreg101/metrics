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

		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(ow, r)
	}
}

func (mux *Mux) mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(metrics2String(mux.storage.GetAll())))
}

func (mux *Mux) metricPage(w http.ResponseWriter, r *http.Request) {
	log := logger.Default()
	name := chi.URLParam(r, "name")
	if v, ok := mux.storage.Get(name); ok {
		if mux.storage.CheckType(name) == chi.URLParam(r, "type") {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(singleMetric2String(v)))
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
		mux.storage.Add(chi.URLParam(r, "name"), storage.Gauge(res))

	case "counter":
		res, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("can't parse %s in double", chi.URLParam(r, "value"))
			return
		}
		w.WriteHeader(http.StatusOK)
		mux.storage.Add(chi.URLParam(r, "name"), storage.Counter(res))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Infof("wrong metric type: %s", chi.URLParam(r, "type"))
	}
}

func (mux *Mux) updateMetricWithBody(w http.ResponseWriter, r *http.Request) {
	log := logger.Default()

	var m communication.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("can't unmarshal json %s to %v", r.Body, m)
		return
	}

	switch m.MType {
	case "counter":
		if m.Delta == nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("empty counter value in request")
			return
		}
		mux.storage.Add(m.ID, storage.Counter(*m.Delta))

		if v, ok := mux.storage.Get(m.ID); ok {
			help := int64(v.(storage.Counter))
			m.Delta = &help
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("internal error; can't find the metric, which should be in storage")
			return
		}

	case "gauge":
		if m.Value == nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("empty gauge value in request")
			return
		}
		mux.storage.Add(m.ID, storage.Gauge(*m.Value))
	}

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

	var m communication.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("can't unmarshal json %s to %v", r.Body, m)
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
	router.Get("/", usingLogger(mux.mainPage))
	router.Get("/value/{type}/{name}", usingLogger(mux.metricPage))
	router.Post("/update/{type}/{name}/{value}", usingLogger(mux.updateMetric))
	router.Post("/update/", usingLogger(mux.updateMetricWithBody))
	router.Post("/value/", usingLogger(mux.getMetric))
	return router
}

func metrics2String(m storage.Metrics) string {
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

func singleMetric2String(v interface{}) string {
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
