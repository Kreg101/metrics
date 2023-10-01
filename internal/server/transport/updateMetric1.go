package transport

import (
	"github.com/Kreg101/metrics/internal/server/entity"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// updateMetric1 проверит, что запрашиваемая метрика существует и переданные тип и значение соответствуют нормам
func (mux *Mux) updateMetric1(w http.ResponseWriter, r *http.Request) {
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
