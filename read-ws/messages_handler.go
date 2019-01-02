package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Message a message struct
type Message struct {
	gorm.Model
	User    User   `json:"user" gorm:"foreignkey:ID"`
	UserID  uint   `json:"userid"`
	Topic   Topic  `json:"topic" gorm:"foreignkey:ID"`
	TopicID uint   `json:"topicid"`
	Data    string `json:"data" gorm:"size:255"`
}

func handlerMessagesGetFunc(w http.ResponseWriter, r *http.Request) {

	var us []Message
	val, err := client.HGetAll("messages").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	for v := range val {
		var u Message
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

func handlerMessageGetFunc(w http.ResponseWriter, r *http.Request) {

	var u Message
	vars := mux.Vars(r)

	messageID := vars["id"]
	val, err := client.HGet("messages", "message_"+messageID).Result()
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
