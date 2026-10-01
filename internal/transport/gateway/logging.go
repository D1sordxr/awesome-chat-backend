package gateway

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"awesome-chat/internal/domain/app/ports"
)

type responseWriter struct {
	http.ResponseWriter
	requestID string
	status    int
}

func (w *responseWriter) WriteHeader(status int) {
	if w.Header().Get(requestIDHeader) == "" {
		w.Header().Set(requestIDHeader, w.requestID)
	}

	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(data)
}

func withLogging(next http.Handler, log ports.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()

		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
			r.Header.Set(requestIDHeader, requestID)
		}

		writer := &responseWriter{ResponseWriter: w, requestID: requestID}

		next.ServeHTTP(writer, r)

		if writer.status == 0 {
			writer.status = http.StatusOK
		}

		fields := []any{
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", writer.status,
			"duration", time.Since(started).String(),
		}

		switch {
		case writer.status >= http.StatusInternalServerError:
			log.Error("HTTP request failed", fields...)
		case writer.status >= http.StatusBadRequest:
			log.Warn("HTTP request rejected", fields...)
		default:
			log.Debug("HTTP request handled", fields...)
		}
	})
}
