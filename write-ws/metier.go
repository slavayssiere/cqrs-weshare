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

func (s server) testCreateTables() {
	s.db.AutoMigrate(&User{})
	if s.db.HasTable(&User{}) == false {
		s.db.CreateTable(&User{})
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
