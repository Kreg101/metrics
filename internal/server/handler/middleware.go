package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Kreg101/metrics/internal/server/logger"
	"io"
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

		duration := time.Since(start)

		log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration.String(),
			"size", responseData.size,
		)
	}
}

func (mux *Mux) check(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Проверяем, что у нас есть ключ
		if mux.key == "" || r.Header.Get("HashSHA256") == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Считываем тело зароса
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			mux.log.Errorf("can't read body of http request: %v\n", err)
			return
		}

		// создаем дополнительные копии ReadCloser
		copy1 := io.NopCloser(bytes.NewBuffer(buf))
		copy2 := io.NopCloser(bytes.NewBuffer(buf))

		// Еще раз прочитываем тело
		body, err := io.ReadAll(copy1)
		if err != nil {
			mux.log.Errorf("can't read body of request: %v\n", err)
			return
		}

		// Генерируем хэш с помощью ключа
		h := hmac.New(sha256.New, []byte(mux.key))
		h.Write(body)
		src := h.Sum(nil)

		dst := make([]byte, hex.EncodedLen(len(src)))
		hex.Encode(dst, src)

		// проверяем хэши на равенство
		if string(dst) == r.Header.Get("HashSHA256") {
			r.Body = copy2
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

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
