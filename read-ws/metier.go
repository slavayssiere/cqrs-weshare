package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// User a user struct
type User struct {
	gorm.Model
	Name         string `json:"username"`
	Email        string `json:"email"`
	Address      string `json:"address" gorm:"size:255"`
	Age          int    `json:"age"`
	CreationTime int64  `json:"creation_time" gorm:"-"`
}

func handlerUsersFunc(w http.ResponseWriter, r *http.Request) {

	var us []User
	val, err := client.HGetAll("users").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	for v := range val {
		var u User
		fmt.Printf("key[%s] value[%s]\n", v, val[v])
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

func handlerUserFunc(w http.ResponseWriter, r *http.Request) {

	var u User
	vars := mux.Vars(r)

	userID := vars["id"]
	val, err := client.HGet("users", "user_"+userID).Result()
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
