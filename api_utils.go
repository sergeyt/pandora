package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	TypeJSON = "application/json"
)

func sendJSON(w http.ResponseWriter, data interface{}, status ...int) error {
	w.Header().Set("Content-Type", TypeJSON)

	if len(status) > 0 {
		w.WriteHeader(status[0])
	}

	marshaller, ok := data.(json.Marshaler)
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

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// TODO check whether it is possible to send error at this phase
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func sendError(w http.ResponseWriter, err error, status ...int) {
	if len(status) == 0 {
		if strings.Contains(err.Error(), "not found") {
			status = []int{http.StatusNotFound}
		} else {
			status = []int{http.StatusInternalServerError}
		}
	}

	data := &struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	sendJSON(w, data, status...)
}