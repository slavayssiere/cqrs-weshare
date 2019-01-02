package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Topic a topic struct
type Topic struct {
	gorm.Model
	Name string `json:"topicname"`
}

func handlerTopicsGetFunc(w http.ResponseWriter, r *http.Request) {

	var us []Topic
	val, err := client.HGetAll("topics").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	for v := range val {
		var u Topic
		err = json.Unmarshal([]byte(val[v]), &u)
		if err != nil {
			panic(err)
		}
		us = append(us, u)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

func handlerTopicGetFunc(w http.ResponseWriter, r *http.Request) {

	var u Topic
	vars := mux.Vars(r)

	topicID := vars["id"]
	val, err := client.HGet("topics", "topic_"+topicID).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	err = json.Unmarshal([]byte(val), &u)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
}
