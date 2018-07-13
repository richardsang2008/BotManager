package model

import "github.com/jinzhu/gorm"

type BotMessage struct {
	UserID    string
	Message   string
	ChannelID string
}



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
	DBId uint
}
type AddLocationCmd struct {
	Latitude  float64
	Longitude float64
	Radius    float64
}
type AddAllMonsCmd struct {
	Lvl *Range
	IV  *Range
}
type AddMonCmd struct {
	Name string
	CP   *Range
	Lvl  *Range
	IV   *Range
	//Move1 *string
	//Move2 *string
}
type AddRaidCmd struct {
	Name string
	AddAllRaidCmd
}
type AddAllRaidCmd struct {
	Sponsored bool
	Lvl       *Range
	Boosted   bool
	Team      string
	GymName   string
}
type AddEggCmd struct {
	Sponsored bool
	Boosted   bool
	Lvl       *Range
	Team      string
	GymName   string
}
type AddGymCmd struct {
	Team    string
	GymName string
}
//the follwoing structs are used for db
type SlackDBMessage struct {
	gorm.Model
	RegionId  int
	ChannelId string
	Ts        float64
}

type SlackDBUserFilter struct {
	gorm.Model
	UserId  uint
	Filters string
}

type SlackDBUser struct {
	gorm.Model
	Referenceid  string
	ChannelId    int
	Fname    string
	Lname     string
	Notifyname         string
	StatusId     int
	Email        string
	Phone        string
	Isadmin      bool
	Isowner      bool
	Isbot        bool
	Realname     string
	AccessRights      string
}
type SlackRegion struct {
	gorm.Model
	Mode string
	RegionId int
	MasterslackToken string
	BotslackToken string
	Lisaslacktoken string
}
