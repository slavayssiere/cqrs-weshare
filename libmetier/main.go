package libmetier

import "github.com/jinzhu/gorm"

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
	UserID  uint   `json:"userid"`
	TopicID uint   `json:"topicid"`
	Data    string `json:"data" gorm:"size:255"`
}

// MessageComplete is MessageComplete struct
type MessageComplete struct {
	User    User   `json:"user"`
	UserID  uint   `json:"userid"`
	TopicID uint   `json:"topicid"`
	Data    string `json:"data"`
}

// TopicComplete is TopicComplete struct
type TopicComplete struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	Conversation []MessageComplete `json:"conversation"`
}
