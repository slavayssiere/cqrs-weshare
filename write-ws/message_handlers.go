package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s server) handlerMessageCreateFunc(w http.ResponseWriter, r *http.Request) {

	var u Message
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}

	s.createMessage(&u)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
}

func (s server) handlerMessageUpdateFunc(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var u Message
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	uID, _ := strconv.ParseUint(id, 10, 64)
	u.ID = uint(uID)

	s.updateMessage(&u)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
}
