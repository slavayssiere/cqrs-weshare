package main

import (
	"github.com/slavayssiere/cqrs-weshare/libmetier"
)

func (s server) eventUserReceive(m *libmetier.User) {
	s.setUser(m)
}

func (s server) eventTopicReceive(m *libmetier.Topic) {
	s.setTopic(m)
}

func (s server) eventMessageReceive(m *libmetier.Message) {
	s.setMessage(m)

	u := s.getUser(m.UserID)
	var mc libmetier.MessageComplete
	mc.TopicID = m.TopicID
	mc.Data = m.Data
	mc.User = u
	mc.UserID = u.ID
	mc.ID = m.ID
	s.setMessageCompleteByTopic(&mc, m.TopicID)

	t := s.getTopic(m.TopicID)
	ml := s.getMessageCompleteByTopic(m.TopicID)
	var tc libmetier.TopicComplete
	tc.ID = m.TopicID
	tc.Conversation = ml
	tc.Name = t.Name

	s.setTopicComplete(&tc)
}
