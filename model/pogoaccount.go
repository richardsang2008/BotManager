package model

import (
	"github.com/jinzhu/gorm"
	"github.com/weilunwu/go-geofence"
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
	Type         string
	Latitude     interface{} `json:"latitude"`  //float64
	Longitude    interface{} `json:"longitude"` //float64
	Address      interface{} //string
	State        interface{} //string
	County       interface{} //string
	Country      interface{} //string
	PokemonLevel interface{} `json:"pokemon_level"` //int
	ChannelId    interface{} //int
	GymName      interface{} //string
	TeamId       interface{} //int
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

type GeoLocation struct {
	Region    *string `json:region`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type RaidData struct {
	GeoLocation
	Move1       *int     `json:"move_1"`
	Move2       *int     `json:"move_2"`
	Cp          *int     `json:"cp"`
	PokemonID   *int     `json:"pokemon_id"` //int
	GymURL      *string  `json:"gym_url"`
	Level       *float64 `json:"level"`         //int
	Base64GymID *string  `json:"base64_gym_id"` //string
	Team        *int     `json:"team"`          //int
	GymID       *string  `json:"gym_id"`        //string
	RaidBegin   *int64   `json:"raid_begin"`    //int64
	RaidSeed    *int64   `json:"raid_seed"`     //int64
	GymName     *string  `json:"gym_name"`      //string
	RaidEnd     *int64   `json:"raid_end"`      //int
	Park        *string  `json:"park"`          //string
	Sponsor     *bool    `json:"sponsor"`
}

type GymData struct {
	GeoLocation
	ID          *string      `json:"id"`          //string
	URL         *string      `json:"url"`         //string
	Name        *string      `json:"name"`        //string
	Description *string      `json:"description"` //string
	Team        *int         `json:"team"`        //int
	Sponsor     *bool        `json:"sponsor"`
	Park        *string      `json:"park"` //string
	Guards      []GymPokemon `json:"pokemon"`
}

type GymPokemon struct {
	NumUpgrades            *int     `json:"num_upgrades"`             //int
	Move1                  *int     `json:"move_1"`                   //int
	Move2                  *int     `json:"move_2"`                   //int
	AdditionalCpMultiplier *int     `json:"additional_cp_multiplier"` //int
	IvDefense              *int     `json:"iv_defense"`               //int
	Weight                 *float64 `json:"weight"`                   //float64
	PokemonID              *int     `json:"pokemon_id"`               //int
	StaminaMax             *int     `json:"stamina_max"`              //int
	CpMultiplier           *float64 `json:"cp_multiplier"`            //float64
	Height                 *float64 `json:"height"`                   //float64
	Stamina                *int     `json:"stamina"`                  //int
	PokemonUID             *int64   `json:"pokemon_uid"`              //int64
	DeploymentTime         *int     `json:"deployment_time"`          //int
	IvAttack               *int     `json:"iv_attack"`                //int
	TrainerName            *string  `json:"trainer_name"`             //string
	TrainerLevel           *int     `json:"trainer_level"`            //int
	Cp                     *int     `json:"cp"`                       //int
	IvStamina              *int     `json:"iv_stamina"`               //int
	CpDecayed              *int     `json:"cp_decayed"`               //int
}
type PokemonData struct {
	GeoLocation
	Move1               *int        `json:"move_1"`
	Move2               *int        `json:"move_2"`
	Cp                  *float64    `json:"cp"`
	PokemonID           *int        `json:"pokemon_id"`            //int
	EncounterID         *uint64     `json:"encounter_id"`          //string
	SpawnpointID        *int64      `json:"spawnpoint_id"`         //string
	PokemonLevel        *float64    `json:"pokemon_level"`         //int
	PlayerLevel         *int        `json:"player_level"`          // int
	Iv                  *float64    `json:iv`                      //float64
	Size                *string     `json:"size"`                  //string
	TinyRat             *string     `json:"tinyrat"`               //string
	BigKarp             *string     `json:"bigkarp"`               //string
	DisappearTime       *int64      `json:"disappear_time"`        //int
	LastModifiedTime    *int64      `json:"last_modified_time"`    //int64
	TimeUntilHiddenMs   *float64    `json:"time_until_hidden_ms"`  //int
	SecondsUntilDespawn *int        `json:"seconds_until_despawn"` //int
	SpawnStart          *int        `json:"spawn_start"`           //int
	SpawnEnd            *int        `json:"spawn_end"`             //int
	Verified            *bool       `json:"verified"`
	CpMultiplier        *float64    `json:"cp_multiplier"` //float64
	Form                interface{} `json:"form"`
	IndividualAttack    *float64    `json:"individual_attack"`
	IndividualDefense   *float64    `json:"individual_defense"`
	IndividualStamina   *float64    `json:"individual_stamina"`
	Height              *float64    `json:"height"`
	Weight              *float64    `json:"weight"`
	Gender              *int        `json:"gender"`
	GMaps               *string     `json:"gmaps"`     //string
	AppleMaps           *string     `json:"applemaps"` //string
}
type PokeMinerMonMessage struct {
	MessageType string      `json:"type"`
	Message     PokemonData `json:"message"`
}
type PokeMinerRaidMessage struct {
	MessageType string   `json:"type"`
	Message     RaidData `json:"message"`
}
type PokeMinerGymMessage struct {
	MessageType string  `json:"type"`
	Message     GymData `json:"message"`
}

//The following is for the User filters
type Range struct {
	Min float64 `json:"Min"`
	Max float64 `json:"Max"`
}
type AddNotifyAll struct {
	Level   *Range `json:"level"`
	Iv      *Range `json:"Iv"`
	Boosted *bool  `json:"Boosted"`
}
type AddNotify struct {
	Mon               *NameAndID  `json:"Mon"`
	Cp                *Range      `json:"Cp"`
	Iv                *Range      `json:"Iv"`
	Level             *Range      `json:"level"`
	Boosted           *bool       `json:"Boosted"`
	Move1             []NameAndID `json:"Move_1"`
	Move2             []NameAndID `json:"Move_2"`
	IndividualStamina *Range      `json:"individual_stamina"`
	IndividualDefense *Range      `json:"individual_defense"`
	IndividualAttack  *Range      `json:"individual_attack"`
	Gender            *NameAndID  `json:"gender"`
}
type AddLocation struct {
	GeoLocation
	Radius *float64 `json:"radius"` //float64
}
type NameAndID struct {
	Name string `json:"Name"`
	Id   int    `json:"Id"`
}
type MonMovesFilter struct {
	Mon   *NameAndID  `json:"Mon"` //string
	Move1 []NameAndID `json:"Move_1"`
	Move2 []NameAndID `json:"Move_2"`
}
type AddNotifyRaid struct {
	EggOrRaid
	MonMovesFilters []MonMovesFilter `json:"MonMoves"`
}
type EggOrRaid struct {
	GymName *string    `json:"gym_name"` //string
	Team    *NameAndID `json:"team"`     //string
	Level   *Range     `json:"level"`
	Sponsor *bool      `json:"sponsor"`
}
type AddNotifyEgg struct {
	EggOrRaid
}
type AddNotifyGym struct {
	Sponsor *bool   `json:"sponsor"`
	GymName *string `json:"gym_name"` //string
}
type UserRegion struct {
	UserId string `json:"UserId"`
	Region string `json:"Region"`
}
type Region struct {
	Region string        `json:"Region"`
	Zone   []GeoLocation `json:"zone"`
}
type Regions struct {
	Regions []Region `json:"Regions"`
}
type Filters struct {
	UserRegion    *UserRegion    `json:"UserRegion"`
	AddNotifyAll  *AddNotifyAll  `json:"AddNotifyAll"`  //AddNotifyAll
	AddNotifies   []AddNotify    `json:"AddNotifies"`   //slice of AddNotify
	AddLocation   *AddLocation   `json:"AddLocation"`   //AddLocation
	AddNotifyRaid *AddNotifyRaid `json:"AddNotifyRaid"` //AddNotifyRaid
	AddNotifyEgg  *AddNotifyEgg  `json:"AddNotifyEgg"`  //AddNotifyEgg
	AddNotifyGym  *AddNotifyGym  `json:"AddNotifyGym"`  //AddNotifyGym
}
type GeoFences struct {
	Region   string
	Geofence *geofence.Geofence
}
