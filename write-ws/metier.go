package main

import (
	"github.com/slavayssiere/cqrs-weshare/libmetier"
)

func (s server) testCreateTables() {
	s.db.AutoMigrate(&libmetier.User{})
	s.db.AutoMigrate(&libmetier.Topic{})
	s.db.AutoMigrate(&libmetier.Message{})

	if s.db.HasTable(&libmetier.User{}) == false {
		s.db.CreateTable(&libmetier.User{})
	}
	if s.db.HasTable(&libmetier.Topic{}) == false {
		s.db.CreateTable(&libmetier.Topic{})
	}
	if s.db.HasTable(&libmetier.Message{}) == false {
		s.db.CreateTable(&libmetier.Message{})
	}
	s.db.BlockGlobalUpdate(true)
}

func (s server) createUser(u *libmetier.User) {
	s.db.Create(u)
	s.ec.Publish("users", u)
}

func (s server) updateUser(u *libmetier.User) {
	s.db.Model(u).Updates(u)
	s.db.First(u, u.ID)
	s.ec.Publish("users", u)
}

func (s server) createTopic(t *libmetier.Topic) {
	s.db.Create(t)
	s.ec.Publish("topics", t)
}

func (s server) updateTopic(t *libmetier.Topic) {
	s.db.Model(t).Updates(t)
	s.db.First(t, t.ID)
	s.ec.Publish("topics", t)
}

func (s server) createMessage(t *libmetier.Message) {
	s.db.Create(t)
	s.ec.Publish("messages", t)
}

func (s server) updateMessage(t *libmetier.Message) {
	s.db.Model(t).Updates(t)
	s.db.First(t, t.ID)
	s.ec.Publish("messages", t)
}
