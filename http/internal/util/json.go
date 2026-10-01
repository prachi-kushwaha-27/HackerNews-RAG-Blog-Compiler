package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func EncodeHttpResponseToJson[T any](w http.ResponseWriter, statusCode int, val T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(val); err != nil {
		return errors.Join(
			fmt.Errorf("failed to encode http response to json: %v", val),
			err)
	}
	return nil
}

func DecodeJsonHttpRequest[T any](req *http.Request) (*T, error) {
	var val T
	if err := json.NewDecoder(req.Body).Decode(&val); err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to decode http request to json: %v", req.Body),
			err)
	}
	return &val, nil
}
