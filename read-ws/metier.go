package main

import (
	"encoding/json"
	"fmt"

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

// Topic a topic struct
type Topic struct {
	gorm.Model
	Name string `json:"topicname"`
}

// User a user struct
type User struct {
	gorm.Model
	Name         string `json:"username"`
	Email        string `json:"email"`
	Address      string `json:"address" gorm:"size:255"`
	Age          int    `json:"age"`
	CreationTime int64  `json:"creation_time" gorm:"-"`
}

func getMessages() []Message {
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
	return us
}

func getTopics() []Topic {
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
	return us
}

func getUsers() []User {
	var us []User
	val, err := client.HGetAll("users").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	for v := range val {
		var u User
		err = json.Unmarshal([]byte(val[v]), &u)
		if err != nil {
			panic(err)
		}
		us = append(us, u)
	}
	return us
}

func getMessage(mID string) Message {
	var m Message
	val, err := client.HGet("messages", "message_"+mID).Result()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(val), &m)
	if err != nil {
		panic(err)
	}
	return m
}

func getTopic(mID string) Topic {
	var m Topic
	val, err := client.HGet("topics", "topic_"+mID).Result()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(val), &m)
	if err != nil {
		panic(err)
	}
	return m
}

func getUser(mID string) User {
	var m User
	val, err := client.HGet("users", "user_"+mID).Result()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(val), &m)
	if err != nil {
		panic(err)
	}
	return m
}

// TopicComplete is TopicComplete struct
type TopicComplete struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Conversation []Message `json:"conversation"`
}

func getTopicComplete(mID string) TopicComplete {
	var tc TopicComplete
	t := getTopic(mID)
	tc.Name = t.Name
	tc.ID = t.ID

	val, err := client.LRange("topic_complete_"+mID, 0, 10).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	for v := range val {
		var m Message
		err = json.Unmarshal([]byte(val[v]), &m)
		if err != nil {
			panic(err)
		}
		tc.Conversation = append(tc.Conversation, m)
	}

	return tc
}
