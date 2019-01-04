package main

import (
	"encoding/json"
	"fmt"

	"github.com/slavayssiere/cqrs-weshare/libmetier"
)

func getMessages() []libmetier.Message {
	var us []libmetier.Message
	val, err := client.HGetAll("messages").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	for v := range val {
		var u libmetier.Message
		err = json.Unmarshal([]byte(val[v]), &u)
		if err != nil {
			panic(err)
		}
		us = append(us, u)
	}
	return us
}

func getTopics() []libmetier.Topic {
	var us []libmetier.Topic
	val, err := client.HGetAll("topics").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	for v := range val {
		var u libmetier.Topic
		err = json.Unmarshal([]byte(val[v]), &u)
		if err != nil {
			panic(err)
		}
		us = append(us, u)
	}
	return us
}

func getUsers() []libmetier.User {
	var us []libmetier.User
	val, err := client.HGetAll("users").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	for v := range val {
		var u libmetier.User
		err = json.Unmarshal([]byte(val[v]), &u)
		if err != nil {
			panic(err)
		}
		us = append(us, u)
	}
	return us
}

func getMessage(mID string) libmetier.Message {
	var m libmetier.Message
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

func getTopic(mID string) libmetier.Topic {
	var m libmetier.Topic
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

func getUser(mID string) libmetier.User {
	var m libmetier.User
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

func getTopicComplete(mID string) libmetier.TopicComplete {
	var tc libmetier.TopicComplete
	val, err := client.HGet("topics_complete", "topic_"+mID).Result()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(val), &tc)
	if err != nil {
		panic(err)
	}
	return tc
}
