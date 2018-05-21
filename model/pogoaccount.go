package model

import (
	"time"
	"github.com/jinzhu/gorm"
)

type PogoAccount struct {
	gorm.Model

	AuthService        string     `json:"auth_service"`
	Username           string     `json:"username"  gorm:"type:varchar(100);unique_index"`
	Password           string     `json:"password"`
	Email              string     `json:"email"`
	LastModified       *time.Time `json:"last_modified,string"`
	ReachLvl30Datetime *time.Time `json:"reach_lvl30_datetime,string"`
	SystemId           string     `json:"system_id"`
	AssignedAt         *time.Time `json:"assigned_at,string"`
	Latitude           float32    `json:"latitude"`
	Longitude          float32    `json:"longitude"`
	Level              int        `json:"level,int"`
	Xp                 int       `json:"xp,int"`
	Encounters         int       `json:"encounters,int"`
	BallsThrown        int       `json:"balls_thrown,int"`
	Captures           int       `json:"captures,int"`
	Spins              int       `json:"spins,int"`
	Walked             float32   `json:"walked"`
	Team               string    `json:"team"`
	Coins              int       `json:"coins int"`
	Stardust           int       `json:"stardust"`
	Warn               bool      `json:"warn"`
	Banned             bool      `json:"banned"`
	BanFlag            bool      `json:"ban_flag"`
	TutorialState      string    `json:"tutorial_state"`
	Captcha            bool      `json:"captcha"`
	RarelessScans      int       `json:"rareless_scans"`
	Shadowbanned       bool      `json:"shadowbanned"`
	Balls              int       `json:"balls"`
	TotalItems         int       `json:"total_items"`
	Pokemon            int       `json:"pokemon"`
	Eggs               int       `json:"eggs"`
	Incubators         int       `json:"incubators"`
	Lures              int       `json:"lures"`
}