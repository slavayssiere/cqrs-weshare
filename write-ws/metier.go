package main

import (
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

// Topic a topic struct
type Topic struct {
	gorm.Model
	Name string `json:"topicname"`
}

// Message a message struct
type Message struct {
	gorm.Model
	User    User   `json:"user" gorm:"foreignkey:ID"`
	UserID  uint   `json:"userid"`
	Topic   Topic  `json:"topic" gorm:"foreignkey:ID"`
	TopicID uint   `json:"topicid"`
	Data    string `json:"data" gorm:"size:255"`
}

func (s server) testCreateTables() {
	s.db.AutoMigrate(&User{})
	s.db.AutoMigrate(&Topic{})
	s.db.AutoMigrate(&Message{})

	if s.db.HasTable(&User{}) == false {
		s.db.CreateTable(&User{})
	}
	if s.db.HasTable(&Topic{}) == false {
		s.db.CreateTable(&Topic{})
	}
	if s.db.HasTable(&Message{}) == false {
		s.db.CreateTable(&Message{})
	}
	s.db.BlockGlobalUpdate(true)
}

func (s server) createUser(u *User) {
	s.db.Create(u)
	s.ec.Publish("users", u)
}

func (s server) updateUser(u *User) {
	s.db.Model(u).Updates(u)
	s.db.First(u, u.ID)
	s.ec.Publish("users", u)
}

func (s server) createTopic(t *Topic) {
	s.db.Create(t)
	s.ec.Publish("topics", t)
}

func (s server) updateTopic(t *Topic) {
	s.db.Model(t).Updates(t)
	s.db.First(t, t.ID)
	s.ec.Publish("topics", t)
}

func (s server) createMessage(t *Message) {
	s.db.Create(t)
	s.ec.Publish("messages", t)
}

func (s server) updateMessage(t *Message) {
	s.db.Model(t).Updates(t)
	s.db.First(t, t.ID)
	s.ec.Publish("messages", t)
}
