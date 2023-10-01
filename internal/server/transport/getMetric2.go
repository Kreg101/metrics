package transport

import (
	"encoding/json"
	"fmt"
	"github.com/Kreg101/metrics/internal/server/entity"
	"net/http"
)

// getMetric2 вернет метрику и StatusOk, если метрика с указанным именем и типом существует в хранилище
// иначе вернет StatusNotFound
func (mux *Mux) getMetric2(w http.ResponseWriter, r *http.Request) {
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
