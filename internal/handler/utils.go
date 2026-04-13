package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const DEFAULT_LIMIT = 20
const DEFAULT_OFFSET = 0
const MAX_LIMIT = 100

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

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	if parsed < 0 {
		return 0, fmt.Errorf("%s must be >= 0", key)
	}
	if key == "limit" && parsed > MAX_LIMIT {
		return 0, fmt.Errorf("%s must be <= %d", key, MAX_LIMIT)
	}

	return parsed, nil
}
