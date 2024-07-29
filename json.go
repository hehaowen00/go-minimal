package gominimal

import (
	"encoding/json"
	"net/http"
)

func Marshal[T any](w http.ResponseWriter, statusCode int, data T) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

func Unmarshal[T any](r http.Request, data *T) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(data)
}
