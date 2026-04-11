package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const DEFAULT_LIMIT = 20
const DEFAULT_OFFSET = 0

func WriteResponseWithEncoder(w http.ResponseWriter, value any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
func ParseQueryInt(r *http.Request, key string, defaultValue int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(val)
}
