package main

import (
	"encoding/json"
	"net/http"
)

const (
	TypeJSON = "application/json"
)

func sendJSON(w http.ResponseWriter, result interface{}) error {
	w.Header().Set("Content-Type", TypeJSON)

	marshaller, ok := result.(json.Marshaler)
	if ok {
		b, err := marshaller.MarshalJSON()
		if err != nil {
			// TODO check whether it is possible to send error at this phase
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		w.Write(b)
		return nil
	}

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		// TODO check whether it is possible to send error at this phase
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func sendError(w http.ResponseWriter, err error) {
	sendJSON(w, &struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
