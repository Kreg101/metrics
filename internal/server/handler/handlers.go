package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Kreg101/metrics/internal/entity"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// mainPage в теле ответа запишет все существующие в хранилище метрики
func (mux *Mux) mainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(mux.storage.GetAll(r.Context()).String()))
}

// metricPage проверит, что существует запрашиваемая метрика и в теле ответа запишет ее
func (mux *Mux) metricPage(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")

	// Проверяю, есть ли метрика с данным именем и типом в хранилище
	if m, ok := mux.storage.Get(r.Context(), name); ok && mType == m.MType {
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte(m.String()))
		return
	}

	mux.log.Infof("no entity %s and type %s in inmemstore", name, mType)
	w.WriteHeader(http.StatusNotFound)
}

// updateMetric проверит, что запрашиваемая метрика существует и переданные тип и значение соответствуют нормам
func (mux *Mux) updateMetric(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "name") == "" {
		w.WriteHeader(http.StatusNotFound)
		mux.log.Errorf("empty entity name")
		return
	}

	switch chi.URLParam(r, "type") {
	case "gauge":
		res, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			mux.log.Errorf("can't parse %s in float", chi.URLParam(r, "value"))
			return
		}
		w.WriteHeader(http.StatusOK)
		mux.storage.Add(r.Context(), entity.Metric{
			ID:    chi.URLParam(r, "name"),
			MType: "gauge",
			Value: &res,
			Delta: nil,
		})

	case "counter":
		res, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			mux.log.Errorf("can't parse %s in double", chi.URLParam(r, "value"))
			return
		}
		w.WriteHeader(http.StatusOK)
		mux.storage.Add(r.Context(), entity.Metric{
			ID:    chi.URLParam(r, "name"),
			MType: "counter",
			Value: nil,
			Delta: &res,
		})

	default:
		w.WriteHeader(http.StatusBadRequest)
		mux.log.Infof("wrong entity type: %s", chi.URLParam(r, "type"))
	}
}

// updateMetricsWithBody в формате json передается метрика, проверяется ее корректность и возвращается она же обновленная
func (mux *Mux) updateMetricWithBody(w http.ResponseWriter, r *http.Request) {
	var m entity.Metric

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		mux.log.Errorf("can't unmarshal json %s to %v", r.Body, m)
		return
	}

	// Проверяем что метрика с таким типом и значением корректна
	if (m.MType == "counter" && m.Delta == nil) || (m.MType == "gauge" && m.Value == nil) {
		w.WriteHeader(http.StatusBadRequest)
		mux.log.Errorf("empty delta or value in request")
		return
	}

	if m.MType == "counter" {
		fmt.Println("update", m.ID, m.MType, *m.Delta, m.Value)
	}

	// За счет того, что поля Delta и Value - указатели, когда мы положим их в хранилище, они обновятся
	// значит не нужно их снова доставать и вручную менять
	mux.storage.Add(r.Context(), m)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mux.log.Errorf("can't marshal %v", m)
		return
	}
}

// getMetric вернет метрику и StatusOk, если метрика с указанным именем и типом существует в хранилище
// иначе вернет StatusNotFound
func (mux *Mux) getMetric(w http.ResponseWriter, r *http.Request) {
	var m entity.Metric

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mux.log.Errorf("can't unmarshal json %s to %v", r.Body, m)
		return
	}

	if v, ok := mux.storage.Get(r.Context(), m.ID); ok && v.MType == m.MType {
		m = v
	} else {
		w.WriteHeader(http.StatusNotFound)
		mux.log.Infof("no %s entity in storage", m.ID)
		return
	}

	if m.MType == "counter" {
		fmt.Println("get", m.ID, m.MType, *m.Delta, m.Value)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mux.log.Errorf("can't marshal %v", m)
		return
	}
}

// ping проверяет соединение с хранилищем
func (mux *Mux) ping(w http.ResponseWriter, r *http.Request) {
	err := mux.storage.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mux.log.Errorf("can't connect to server: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// updates позволяет обновлять хранилище не 1 метрикой, а целым батчем
func (mux *Mux) updates(w http.ResponseWriter, r *http.Request) {
	var metrics []entity.Metric

	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mux.log.Errorf("can't unmarshal json %s to metrics slice", r.Body)
		return
	}

	for _, m := range metrics {

		if (m.MType == "counter" && m.Delta == nil) || (m.MType == "gauge" && m.Value == nil) {
			mux.log.Errorf("empty delta or value in request")
			continue
		}

		if m.MType == "counter" {
			fmt.Println("updates", m.ID, m.MType, *m.Delta, m.Value)
		}

		fmt.Println(m.ID, m.MType, m.Delta, m.Value)

		mux.storage.Add(r.Context(), m)
	}
}
