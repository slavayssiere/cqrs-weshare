package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s server) handlerUserCreateFunc(w http.ResponseWriter, r *http.Request) {

	var u User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}

	requestAt := time.Now()
	s.createUser(&u)
	duration := time.Since(requestAt)
	u.CreationTime = duration.Nanoseconds()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
}

func (s server) handlerUserUpdateFunc(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var u User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	uID, _ := strconv.ParseUint(id, 10, 64)
	u.ID = uint(uID)

	requestAt := time.Now()
	s.updateUser(&u)
	duration := time.Since(requestAt)
	u.CreationTime = duration.Nanoseconds()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
}
