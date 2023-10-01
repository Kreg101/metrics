package transport

import "net/http"

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
