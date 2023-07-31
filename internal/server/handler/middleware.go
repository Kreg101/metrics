package handler

import (
	"github.com/Kreg101/metrics/internal/server/logger"
	"net/http"
	"strings"
	"time"
)

// logging логирует запрос и ответ посредством middleware
func logging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// При вызове logger.Default() возращается единый на весь сервер логгер
		log := logger.Default()

		start := time.Now()

		responseData := &responseData{}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		// Подменяет w, на свой с логированием
		h.ServeHTTP(&lw, r)

		duration := time.Since(start).Milliseconds()

		log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
}

// compression позволяет разжимать запрос и сжимать ответ, если такое возможно
func compression(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		// Если ответ запрос допускает сжатие ответа форматом gzip,
		// мы подменим ответ на наш энкодер
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			ow.Header().Set("Content-Encoding", "gzip")
			defer cw.Close()
		}

		// Если тело запроса закодировано gzip, ты мы подменим тело на наш декодер
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	}
}
