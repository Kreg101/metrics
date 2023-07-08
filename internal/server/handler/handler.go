package handler

import (
	"fmt"
	"github.com/Kreg101/metrics/internal/server/constants"
	"github.com/Kreg101/metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Mux struct {
	storage *storage.Storage
	router  chi.Router
}

func NewMux() *Mux {
	mux := &Mux{}
	mux.storage = storage.NewStorage()
	mux.router = chi.NewRouter()
	return mux
}

func (mux *Mux) Apply() chi.Router {
	mux.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "text/html")
		w.Write([]byte(mux.storage.GetAllString()))
	})

	mux.router.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		if v, ok := mux.storage.GetString(name); ok {
			if mux.storage.CheckType(name) == chi.URLParam(r, "type") {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("content-type", "text/html")
				w.Write([]byte(v))
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})

	mux.router.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		if chi.URLParam(r, "name") == constants.EmptyString {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch chi.URLParam(r, "type") {
		case constants.GaugeType:
			fmt.Println(chi.URLParam(r, "type"), " ", chi.URLParam(r, "name"), " ", chi.URLParam(r, "value"))
			res, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
			mux.storage.Add(chi.URLParam(r, "name"), storage.Gauge(res))
		case constants.CounterType:
			res, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
			mux.storage.Add(chi.URLParam(r, "name"), storage.Counter(res))
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	return mux.router
}
