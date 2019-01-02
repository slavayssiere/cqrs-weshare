package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// User a user struct
type User struct {
	Name         string    `json:"username"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	Age          int       `json:"age"`
	CreationTime int64     `json:"creation_time"`
	CreateTime   time.Time `json:"create_at"`
	ID           int       `json:"id"`
}

func (s server) eventReceive(m *User) {
	log.Println(m)
	log.Println("user_" + strconv.Itoa(m.ID))
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.client.HSet("users", "user_"+strconv.Itoa(m.ID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}
}
