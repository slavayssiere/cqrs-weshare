package main

import (
	"encoding/json"
	"net/http"
)

// Status a test status struct
type Status struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}

func handlerHealthFunc(w http.ResponseWriter, r *http.Request) {
	stt := Status{Name: "OK", Code: 200}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stt); err != nil {
		panic(err)
	}
}

func handlerStatusFunc(w http.ResponseWriter, r *http.Request) {
	var stt Status
	stt.Name = "OK"
	stt.Code = 200
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stt); err != nil {
		panic(err)
	}
}
