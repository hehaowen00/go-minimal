package gominimal

import (
	"encoding/json"
	"net/http"
)

type JSON map[string]any

func MarshalJSON[T any](w http.ResponseWriter, statusCode int, data T) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

func UnmarshalJSON[T any](r http.Request, data *T) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(data)
}
