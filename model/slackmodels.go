package model

import "github.com/jinzhu/gorm"

type BotMessage struct {
	UserID string
	Message string
	ChannelID string
}
type Rights int

const (
	USER Rights = 1 + iota
	MODERATOR
	ADMINISTRATOR
	SUPER
	OWNER
)
type SlackUser struct {
	ID           string
	ChannelID    int
	StatusID     int
	LevelID      int
	ReferenceID  string
	Name         string
	Deleted      bool
	RealName     string
	Phone        string
	DisplayName  string
	StatusText   string
	Email        string
	FirstName    string
	LastName     string
	IsAdmin      bool
	IsOwner      bool
	IsBot        bool
	AccessRights Rights
}
type SlackMessage struct {
	gorm.Model
	RegionId int
	ChannelId string
	Ts  float64
}
type AddLocationCmd struct {
	Latitude float64
	Longitude float64
	Radius float64
}
type SlackUserFilter struct {
	gorm.Model
	UserId int
	Filters string
}