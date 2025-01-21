package main

import (
	"encoding/json"
	"net/http"
)

func decodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func encodeJSON(w http.ResponseWriter, v interface{}, resCode int) error {
	dat, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resCode)
	w.Write(dat)
	return nil
}

func returnError(w http.ResponseWriter, message string, statusCode int) {
	type errorJSON struct {
		Error string `json:"error"`
	}
	errorResp := errorJSON{Error: message}
	dat, _ := json.Marshal(errorResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}
