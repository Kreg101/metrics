package handler

import (
	"github.com/Kreg101/metrics/internal/constants"
	"github.com/Kreg101/metrics/internal/memory"

	"net/http"
	"strconv"
	"strings"
)

type Mux struct {
	storage *memory.Storage
}

func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	split := strings.Split(r.RequestURI, "/")
	removeEmptyElements(&split)

	res, err := requestValidation(split)

	switch err {
	case constants.InvalidRequestError:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case constants.NoMetricNameError:
		w.WriteHeader(http.StatusNotFound)
		return
	case constants.InvalidMetricTypeError:
		w.WriteHeader(http.StatusBadRequest)
		return
	case constants.InvalidValueError:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mux.storage.Add(split[2], res)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func NewMux() *Mux {
	mux := &Mux{}
	mux.storage = memory.NewStorage()
	return mux
}

func requestValidation(split []string) (interface{}, constants.Error) {
	if len(split) < 3 {
		return nil, constants.InvalidRequestError
	}
	if (split)[2] == constants.EmptyString {
		return nil, constants.NoMetricNameError
	}
	if len(split) != 4 {
		return nil, constants.InvalidRequestError
	}
	if split[1] != constants.CounterType && split[1] != constants.GaugeType {
		return nil, constants.InvalidMetricTypeError
	}
	if split[1] == constants.CounterType {
		res, err := strconv.ParseInt(split[3], 10, 64)
		if err != nil {
			return nil, constants.InvalidValueError
		}
		return memory.Counter(res), constants.NoError
	} else {
		res, err := strconv.ParseFloat(split[3], 64)
		if err != nil {
			return nil, constants.InvalidValueError
		}
		return memory.Gauge(res), constants.NoError
	}
}

func removeEmptyElements(split *[]string) {
	if (*split)[0] == "" {
		*split = (*split)[1:]
	}
	if (*split)[len(*split)-1] == "" {
		*split = (*split)[:len(*split)-1]
	}
}
