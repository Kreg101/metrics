package transport

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// getMetric1 проверит, что существует запрашиваемая метрика и в теле ответа запишет ее
func (mux *Mux) getMetric1(w http.ResponseWriter, r *http.Request) {
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
