package handler

import (
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
		w.Write([]byte(mux.storage.GetAll()))
	})

	mux.router.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "type")
		if v, ok := mux.storage.Get(name); ok {
			if mux.storage.CheckType(name) == chi.URLParam(r, "type") {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("content-type", "text/html")
				w.Write([]byte(name + ":" + v.(string)))
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
			//fmt.Println(chi.URLParam(r, "type"), " ", chi.URLParam(r, "name"), " ", chi.URLParam(r, "value"))
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

	//split := strings.Split(r.RequestURI, "/")
	//split = removeFirstLastEmptyElements(split)
	//
	//res, err := requestValidation(split)
	//
	//if err == constants.InvalidRequestError || err == constants.InvalidMetricTypeError || err == constants.InvalidValueError {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//} else if err == constants.NoMetricNameError {
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//}
	//
	//mux.storage.Add(split[2], res)
	//w.Header().Set("content-type", "application/json")
	//w.WriteHeader(http.StatusOK)
}

//func requestValidation(split []string) (interface{}, constants.Error) {
//	if len(split) < 3 {
//		return nil, constants.NoMetricNameError
//	}
//	if split[0] != constants.Update {
//		return nil, constants.InvalidRequestError
//	}
//	if split[2] == constants.EmptyString {
//		return nil, constants.NoMetricNameError
//	}
//	if len(split) != 4 {
//		return nil, constants.InvalidRequestError
//	}
//	if split[1] != constants.CounterType && split[1] != constants.GaugeType {
//		return nil, constants.InvalidMetricTypeError
//	}
//	if split[1] == constants.CounterType {
//		res, err := strconv.ParseInt(split[3], 10, 64)
//		if err != nil {
//			return nil, constants.InvalidValueError
//		}
//		return storage.Counter(res), constants.NoError
//	} else {
//		res, err := strconv.ParseFloat(split[3], 64)
//		if err != nil {
//			return nil, constants.InvalidValueError
//		}
//		return storage.Gauge(res), constants.NoError
//	}
//}
//
//func removeFirstLastEmptyElements(split []string) []string {
//	if (split)[0] == "" {
//		split = (split)[1:]
//	}
//	if (split)[len(split)-1] == "" {
//		split = (split)[:len(split)-1]
//	}
//	return split
//}
