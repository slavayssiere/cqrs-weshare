package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// User a user struct
type User struct {
	gorm.Model
	Name         string    `json:"username"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	Age          int       `json:"age"`
	CreationTime int64     `json:"creation_time"`
	CreateTime   time.Time `json:"create_at"`
}

// Topic a topic struct
type Topic struct {
	gorm.Model
	Name string `json:"topicname"`
}

// Message a message struct
type Message struct {
	gorm.Model
	UserID  uint   `json:"userid"`
	TopicID uint   `json:"topicid"`
	Data    string `json:"data" gorm:"size:255"`
}

func (s server) eventUserReceive(m *User) {
	log.Println(m)
	log.Println("user_" + fmt.Sprint(m.ID))
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.client.HSet("users", "user_"+fmt.Sprint(m.ID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (s server) eventTopicReceive(m *Topic) {
	log.Println(m)
	log.Println("topic_" + fmt.Sprint(m.ID))
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.client.HSet("topics", "topic_"+fmt.Sprint(m.ID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (s server) eventMessageReceive(m *Message) {
	log.Println(m)
	log.Println("user_" + fmt.Sprint(m.ID))
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: something here

	err = s.client.HSet("messages", "message_"+fmt.Sprint(m.ID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}
}
