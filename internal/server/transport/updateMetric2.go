package transport

import (
	"encoding/json"
	"fmt"
	"github.com/Kreg101/metrics/internal/server/entity"
	"net/http"
)

// updateMetric2 в формате json передается метрика, проверяется ее корректность и возвращается она же обновленная
func (mux *Mux) updateMetric2(w http.ResponseWriter, r *http.Request) {
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
