// ─── Uniform JSON Response Helpers ────────────────────────────
package handler

import (
	"encoding/json"
	"net/http"
)

// JSONResponse is the standard envelope for all API responses.
type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
	})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"message": msg})
}

func writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, data)
}

// parseInt converts a string to int with a fallback default.
func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return def
		}
		n = n*10 + int(c-'0')
	}
	return n
}
