package transport

import (
	"encoding/json"
	"github.com/Kreg101/metrics/internal/server/entity"
	"net/http"
)

// updatesMetrics позволяет обновлять хранилище не 1 метрикой, а целым батчем
func (mux *Mux) updateMetrics(w http.ResponseWriter, r *http.Request) {
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

		mux.storage.Add(r.Context(), m)
	}
}
