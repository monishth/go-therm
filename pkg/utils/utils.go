package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func encodeHTTPResponse[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decodeHTTPResponse[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("failed to decode request body: %w", err)
	}
	return v, nil
}

func decodeValidResponse[T Validator](r *http.Request) (T, map[string]string, error) {
	v, err := decodeHTTPResponse[T](r)

	if err != nil {
		return v, nil, err
	}

	if problems := v.Valid(r.Context()); problems != nil {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}

func Decode[T any](bytes []byte) (T, error) {
	var v T
	if err := json.Unmarshal(bytes, &v); err != nil {
		return v, fmt.Errorf("failed to decode request body: %w", err)
	}
	return v, nil
}

func ConvertOnOffToInt(value string) int {
	switch value {
	case "ON":
		return 1
	case "OFF":
		return 0
	default:
		panic("invalid value for ON/OFF")
	}
}
