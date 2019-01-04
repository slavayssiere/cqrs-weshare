package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/slavayssiere/cqrs-weshare/libmetier"
)

func (s server) setUser(m *libmetier.User) {
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

func (s server) getUser(id uint) libmetier.User {
	var u libmetier.User
	val, err := s.client.HGet("users", "user_"+fmt.Sprint(id)).Result()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(val), &u)
	if err != nil {
		panic(err)
	}
	return u
}

func (s server) setTopic(m *libmetier.Topic) {
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

func (s server) getTopic(id uint) libmetier.Topic {
	var u libmetier.Topic
	val, err := s.client.HGet("topics", "topic_"+fmt.Sprint(id)).Result()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(val), &u)
	if err != nil {
		panic(err)
	}
	return u
}

func (s server) setMessage(m *libmetier.Message) {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.client.HSet("messages", "message_"+fmt.Sprint(m.ID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (s server) setMessageCompleteByTopic(m *libmetier.MessageComplete, topicid uint) {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.client.HSet("messages_topics", "topic_"+fmt.Sprint(topicid)+"_"+fmt.Sprint(m.ID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (s server) getMessageCompleteByTopic(topicid uint) []libmetier.MessageComplete {
	var ml []libmetier.MessageComplete
	keys, _ := s.client.HKeys("messages_topics").Result()
	for i := range keys {
		if strings.Contains(fmt.Sprint(keys[i]), "topic_"+fmt.Sprint(topicid)) {
			log.Printf("key: %s\n", fmt.Sprint(keys[i]))
			var m libmetier.MessageComplete
			val, err := s.client.HGet("messages_topics", fmt.Sprint(keys[i])).Result()
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal([]byte(val), &m)
			if err != nil {
				panic(err)
			}
			ml = append(ml, m)
		}
	}
	return ml
}

func (s server) setTopicComplete(m *libmetier.TopicComplete) {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.client.HSet("topics_complete", "topic_"+fmt.Sprint(m.ID), string(b)).Err()
	if err != nil {
		log.Fatal(err)
	}
}
