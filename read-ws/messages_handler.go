package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func handlerMessagesGetFunc(w http.ResponseWriter, r *http.Request) {

	us := getMessages()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

func handlerMessageGetFunc(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	messageID := vars["id"]

	m := getMessage(messageID)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(m); err != nil {
		panic(err)
	}
}
