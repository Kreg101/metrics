package transport

import "net/http"

// getMetrics в теле ответа запишет все существующие в хранилище метрики
func (mux *Mux) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(mux.storage.GetAll(r.Context()).String()))
}
