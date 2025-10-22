package health

import (
	"encoding/json"
	"net/http"
)

type Config struct {
	ServiceName string
    Version     string
}

func NewHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status":  "ok",
			"service": cfg.ServiceName,
			"version": cfg.Version,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(response); err != nil {
            // Nothing to do: response writer may be closed by client; best effort only
        }
	}
}