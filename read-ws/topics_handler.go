package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func handlerTopicsGetFunc(w http.ResponseWriter, r *http.Request) {
	us := getTopics()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

func handlerTopicGetFunc(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	topicID := vars["id"]

	u := getTopic(topicID)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
}

func handlerTopicCompleteGetFunc(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	topicID := vars["id"]

	us := getTopicComplete(topicID)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}
