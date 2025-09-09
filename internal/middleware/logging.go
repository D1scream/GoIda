package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем обертку для ResponseWriter чтобы получить статус код
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Определяем уровень логирования по статус коду
		level := logrus.InfoLevel
		if wrapped.statusCode >= 400 {
			level = logrus.WarnLevel
		}
		if wrapped.statusCode >= 500 {
			level = logrus.ErrorLevel
		}

		logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      wrapped.statusCode,
			"duration":    duration,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}).Log(level, "HTTP request")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
