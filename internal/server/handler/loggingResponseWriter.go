package handler

import "net/http"

type (

	// responseData - структура, хранящая данные о запросе
	responseData struct {
		status int
		size   int
	}

	// loggingResponseWriter реализует интерфейс http.ResponseWrite, поэтому подменяется в
	// middleware и получает необходимую информацию для responseData
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
