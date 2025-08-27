package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, v any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, msg string){
	writeJSON(w, status, map[string]any{"error": msg})
}

func withLogging(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("request", "method", r.Method, "path", r.URL.Path, "took", time.Since(start).String())
	})
}
