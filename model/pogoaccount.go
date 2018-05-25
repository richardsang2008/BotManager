package model

import (
	"github.com/jinzhu/gorm"
	"time"
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
	Xp                 int        `json:"xp,int"`
	Encounters         int        `json:"encounters,int"`
	BallsThrown        int        `json:"balls_thrown,int"`
	Captures           int        `json:"captures,int"`
	Spins              int        `json:"spins,int"`
	Walked             float32    `json:"walked"`
	Team               string     `json:"team"`
	Coins              int        `json:"coins int"`
	Stardust           int        `json:"stardust"`
	Warn               bool       `json:"warn"`
	Banned             bool       `json:"banned"`
	BanFlag            bool       `json:"ban_flag"`
	TutorialState      string     `json:"tutorial_state"`
	Captcha            bool       `json:"captcha"`
	RarelessScans      int        `json:"rareless_scans"`
	Shadowbanned       bool       `json:"shadowbanned"`
	Balls              int        `json:"balls"`
	TotalItems         int        `json:"total_items"`
	Pokemon            int        `json:"pokemon"`
	Eggs               int        `json:"eggs"`
	Incubators         int        `json:"incubators"`
	Lures              int        `json:"lures"`
}

type BaseMon struct {
	Type string
	gorm.Model
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Address      string
	State        string
	County       string
	Country      string
	PokemonLevel int `json:"pokemon_level"`
	ChannelId    int
	GymName      string
	TeamId       int
}

type BaseStats struct {
	Attack     int     `json:"attack"`
	Defense    int     `json:"defense"`
	Stamina    int     `json:"stamina"`
	Type1      int     `json:"type1"`
	Type2      int     `json:"type2"`
	Legendary  bool    `json:"legendary"`
	Generation int     `json:"generation"`
	Weight     float64 `json:"weight"`
	Height     float64 `json:"height"`
}

type PokemonData struct {
	BaseMon
	EncounterID         string  `json:"encounter_id"`
	SpawnpointID        string  `json:"spawnpoint_id"`
	PokemonID           int     `json:"pokemon_id"`
	PlayerLevel         int     `json:"player_level"`
	DisappearTime       int     `json:"disappear_time"`
	LastModifiedTime    int64   `json:"last_modified_time"`
	TimeUntilHiddenMs   int     `json:"time_until_hidden_ms"`
	SecondsUntilDespawn int     `json:"seconds_until_despawn"`
	SpawnStart          int     `json:"spawn_start"`
	SpawnEnd            int     `json:"spawn_end"`
	Verified            bool    `json:"verified"`
	CpMultiplier        float64 `json:"cp_multiplier"`
	Form                int     `json:"form"`
	Cp                  int     `json:"cp"`
	Iv                  float64 `json:iv`
	IndividualAttack    int     `json:"individual_attack"`
	IndividualDefense   int     `json:"individual_defense"`
	IndividualStamina   int     `json:"individual_stamina"`
	Move1               int     `json:"move_1"`
	Move2               int     `json:"move_2"`
	Height              int     `json:"height"`
	Weight              int     `json:"weight"`
	Gender              int     `json:"gender"`
	Size                string  `json:"size"`
	TinyRat             string  `json:tinyrat`
	BigKarp             string  `json:bigkarp`
	gmaps               string  `json:gmaps`
	applemaps           string  `json:applemaps`
}
