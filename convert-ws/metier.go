package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/slavayssiere/cqrs-weshare/libmetier"
)

// // User a user struct
// type User struct {
// 	gorm.Model
// 	Name         string    `json:"username"`
// 	Email        string    `json:"email"`
// 	Address      string    `json:"address"`
// 	Age          int       `json:"age"`
// 	CreationTime int64     `json:"creation_time"`
// 	CreateTime   time.Time `json:"create_at"`
// }

// // Topic a topic struct
// type Topic struct {
// 	gorm.Model
// 	Name string `json:"topicname"`
// }

// // Message a message struct
// type Message struct {
// 	gorm.Model
// 	UserID  uint   `json:"userid"`
// 	TopicID uint   `json:"topicid"`
// 	Data    string `json:"data" gorm:"size:255"`
// }

func (s server) eventUserReceive(m *libmetier.User) {
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

func (s server) eventTopicReceive(m *libmetier.Topic) {
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

// // MessageComplete a message in redis struct
// type MessageComplete struct {
// 	User    User   `json:"user"`
// 	UserID  uint   `json:"userid"`
// 	TopicID uint   `json:"topicid"`
// 	Data    string `json:"data"`
// }

func (s server) eventMessageReceive(m *libmetier.Message) {
	log.Println(m)
	log.Println("user_" + fmt.Sprint(m.ID))

	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.client.HSet("messages", "message_"+fmt.Sprint(m.ID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}

	var u libmetier.User
	val, err := s.client.HGet("users", "user_"+fmt.Sprint(m.UserID)).Result()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(val), &u)
	if err != nil {
		panic(err)
	}

	var mc libmetier.MessageComplete
	mc.TopicID = m.TopicID
	mc.Data = m.Data
	mc.User = u
	mc.UserID = u.ID
	b, err = json.Marshal(mc)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.client.RPush("topic_complete_"+fmt.Sprint(m.TopicID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}
}
