package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kellydunn/golang-geo"
	"github.com/richardsang2008/BotManager/model"
	"math"
	"strings"
	"time"
)

type PokeUtility struct {
}

func (u *PokeUtility) LocateValueWithinMapIntString(str string, mapstr *map[int]string) (int, error) {
	if len(str) == 0 {
		MLog.Error("values is empty")
		return -1, nil
	} else {
		for k, v := range *mapstr {
			if strings.EqualFold(v, str) {
				return k, nil
			}
		}
		return -1, nil
	}
}
func (u *PokeUtility) LocateKeyWithinMapIntString(key int, mapstr map[int]string) (interface{}, error) {
	if val, ok := mapstr[key]; ok {
		return val, nil
	} else {
		MLog.Error("key does not exist ", key)
		return nil, errors.New("key does not exist")
	}
}
func (u *PokeUtility) CalculateTwoPointsDistanceInUnits(p1lan, p1lng, p2lan, p2lng float64, meansure model.MeansureUnit) float64 {
	p := geo.NewPoint(34.117671, -118.073250)
	p2 := geo.NewPoint(34.114826, -118.075295)
	// find the great circle distance between them in km
	dist := p.GreatCircleDistance(p2)
	ret := 0.0
	switch meansure {
	case model.Miles:
		ret = dist / 0.621371
	case model.Meters:
		ret = dist
	default:
		ret = dist
	}
	// change the km to miles

	return ret
}
func (u *PokeUtility) BackFillIdForFilters(filters *model.Filters, pokemonMap map[int]string, moveMap map[int]string, teamsMap map[int]string) *model.Filters {
	//now try to clean the filter object
	if filters.AddNotifies != nil {
		for _, element := range filters.AddNotifies {
			//check mon id
			if element.Mon != nil {
				//turn mon name into id
				id, err := u.LocateValueWithinMapIntString(element.Mon.Name, &pokemonMap)
				if err != nil {
					MLog.Error(err)
				}
				element.Mon.Id = id
			}
			//check move1 id
			if element.Move1 != nil {
				for i, move1 := range element.Move1 {
					id, err := u.LocateValueWithinMapIntString(move1.Name, &moveMap)
					if err != nil {
						MLog.Error(err)
					}
					element.Move1[i].Id = id
				}
			}
			//check move2 id
			if element.Move2 != nil {
				for i, move2 := range element.Move2 {
					id, err := u.LocateValueWithinMapIntString(move2.Name, &moveMap)
					if err != nil {
						MLog.Error(err)
					}
					element.Move2[i].Id = id
				}
			}
		}
	}
	//now try to clean the raid, egg filter
	if filters.AddNotifyRaid != nil {
		if filters.AddNotifyRaid.Team != nil {
			id, err := u.LocateValueWithinMapIntString(filters.AddNotifyRaid.Team.Name, &teamsMap)
			if err != nil {
				MLog.Error(err)
			}
			filters.AddNotifyRaid.Team.Id = id
		}
		if filters.AddNotifyRaid.MonMovesFilters != nil {
			for _, element := range filters.AddNotifyRaid.MonMovesFilters {
				if element.Move1 != nil {
					for i, move1 := range element.Move1 {
						id, err := u.LocateValueWithinMapIntString(move1.Name, &moveMap)
						if err != nil {
							MLog.Error(err)
						}
						element.Move1[i].Id = id
					}
				}
				//check move2 id
				if element.Move2 != nil {
					for i, move2 := range element.Move2 {
						id, err := u.LocateValueWithinMapIntString(move2.Name, &moveMap)
						if err != nil {
							MLog.Error(err)
						}
						element.Move2[i].Id = id
					}
				}
			}
		}
	}
	if filters.AddNotifyEgg != nil {
		if filters.AddNotifyEgg.Team != nil {
			id, err := u.LocateValueWithinMapIntString(filters.AddNotifyEgg.Team.Name, &teamsMap)
			if err != nil {
				MLog.Error(err)
			}
			filters.AddNotifyEgg.Team.Id = id
		}
	}
	return filters
}
func (u *PokeUtility) whichRegion (regions []model.GeoFences, lat float64, lng float64) *string {
	ret :=""
	for _,element :=range regions{
		//calculate the zone
		if element.Geofence.Inside(geo.NewPoint(lat,lng)) {
			ret = element.Region
			return &ret
		}
	}
	return nil
}
func (u *PokeUtility) ParsePokeMinerInput(data []byte, regions []model.GeoFences, isTest bool) (interface{}, *bool, *string, error) {
	//data is the json in byte array
	var ret interface{}
	//var region *string
	isWithinTime := false
	gen := make(map[string]interface{})
	now := time.Now()
	if err := json.Unmarshal(data, &gen); err != nil {
		MLog.Error(err)
		return nil, nil, nil, err
	}
	var regionstr *string
	if val, ok := gen["type"]; ok {
		pokemessage := &model.PokeMinerMonMessage{}
		raidmessage := &model.PokeMinerRaidMessage{}
		gymmessage := &model.PokeMinerGymMessage{}
		switch val {
		case "pokemon":
			if err := json.Unmarshal([]byte(data), &pokemessage); err != nil {
				MLog.Error(err)
			}
			ret = pokemessage
			expireTime := time.Unix(*(pokemessage.Message.DisappearTime), 0)
			if expireTime.Sub(now) > 0 {
				isWithinTime = true
			}
			//find out which region
			regionstr = u.whichRegion(regions, pokemessage.Message.Latitude,pokemessage.Message.Longitude)

		case "raid":
			if err := json.Unmarshal([]byte(data), &raidmessage); err != nil {
				MLog.Error(err)
			}
			ret = raidmessage
			expireTime := time.Unix(*(raidmessage.Message.RaidEnd), 0)
			if expireTime.Sub(now) > 0 {
				isWithinTime = true
			}
			//find out which region
			regionstr = u.whichRegion(regions, raidmessage.Message.Latitude,raidmessage.Message.Longitude)
		case "gym":
			if err := json.Unmarshal([]byte(data), &gymmessage); err != nil {
				MLog.Error(err)
			}
			ret = gymmessage
			//find out which region
			regionstr = u.whichRegion(regions, gymmessage.Message.Latitude,gymmessage.Message.Longitude)

		default:
			MLog.Error("input is not supported input type is ", val)
			return nil, nil,nil, errors.New("input is not supported")
		}
	}
	if isTest {
		isWithinTime = true
		return ret, &isWithinTime, regionstr,nil
	}
	return ret, &isWithinTime, regionstr,nil
}

// Calculates the Haversine distance between two points in kilometers.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS = 6371
	KMTOMILES    = 0.62137
)

func (u *PokeUtility) GreatCircleDistance(p2 *model.GeoLocation, p *model.GeoLocation) float64 {
	dLat := (p2.Latitude - p.Latitude) * (math.Pi / 180.0)
	dLon := (p2.Longitude - p.Longitude) * (math.Pi / 180.0)
	lat1 := p.Latitude * (math.Pi / 180.0)
	lat2 := p2.Latitude * (math.Pi / 180.0)
	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)
	a := a1 + a2
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EARTH_RADIUS * c
}
func (u *PokeUtility) ApplyFiltersToGymMessage(gym *model.PokeMinerGymMessage, filters *model.Filters) bool {
	isOkToAlter := false
	locationbool := false
	if filters.AddNotifyGym != nil {
		radius := u.GreatCircleDistance(&(gym.Message.GeoLocation), &(filters.AddLocation.GeoLocation)) * KMTOMILES
		d := *(filters.AddLocation.Radius)
		if radius > d {
			return false
		} else {
			locationbool = true
			enabledSum := 0
			boolSum := 0
			if filters.AddNotifyGym.GymName != nil {
				enabledSum += 1
				if strings.Contains(*(gym.Message.Name), *(filters.AddNotifyRaid.GymName)) {
					boolSum += 1
				}
			}
			if filters.AddNotifyGym.Sponsor != nil {
				enabledSum += 1
				if *(gym.Message.Sponsor) == *(filters.AddNotifyGym.Sponsor) {
					boolSum += 1
				}
			}
			if locationbool && (boolSum == enabledSum && boolSum != 0) {
				isOkToAlter = true
			}
		}
	}
	return isOkToAlter
}
func (u *PokeUtility) ApplyFiltersToRaidOrEggMessage(raidoregg *model.PokeMinerRaidMessage, filters *model.Filters) bool {
	isOkToAlter := false
	locationbool := false

	if filters.AddLocation != nil {
		radius := u.GreatCircleDistance(&(raidoregg.Message.GeoLocation), &(filters.AddLocation.GeoLocation)) * KMTOMILES
		d := *(filters.AddLocation.Radius)
		if radius > d {
			return false
		} else {
			locationbool = true
			enabledSum := 0
			boolSum := 0

			if raidoregg.Message.Cp == nil || *(raidoregg.Message.Cp) == 0 {
				//this is egg
				if filters.AddNotifyEgg != nil {
					if filters.AddNotifyEgg.Level != nil {
						enabledSum += 1
						if *(raidoregg.Message.Level) <= filters.AddNotifyEgg.Level.Max && *(raidoregg.Message.Level) >= filters.AddNotifyEgg.Level.Min {
							boolSum += 1
						}
					}
					if filters.AddNotifyEgg.Team != nil {
						enabledSum += 1
						if *(raidoregg.Message.Team) == filters.AddNotifyEgg.Team.Id {
							boolSum += 1
						}
					}
					if filters.AddNotifyEgg.GymName != nil {
						enabledSum += 1
						if raidoregg.Message.GymName != nil && strings.Contains(*(raidoregg.Message.GymName), *(filters.AddNotifyEgg.GymName)) {
							boolSum += 1
						}
					}
					if filters.AddNotifyEgg.Sponsor != nil {
						enabledSum += 1
						if raidoregg.Message.Sponsor != nil && *(raidoregg.Message.Sponsor) == *(filters.AddNotifyEgg.Sponsor) {
							boolSum += 1
						}
					}
					if locationbool && (boolSum == enabledSum && boolSum != 0) {
						isOkToAlter = true
					}
				}
			} else {
				//this is the raid, check mon and move
				if filters.AddNotifyRaid != nil {
					if filters.AddNotifyRaid.Level != nil {
						enabledSum += 1
						if *(raidoregg.Message.Level) <= filters.AddNotifyRaid.Level.Max && *(raidoregg.Message.Level) >= filters.AddNotifyRaid.Level.Min {
							boolSum += 1
						}
					}
					if filters.AddNotifyRaid.Team != nil {
						enabledSum += 1
						if *(raidoregg.Message.Team) == filters.AddNotifyRaid.Team.Id {
							boolSum += 1
						}
					}
					if filters.AddNotifyRaid.GymName != nil {
						enabledSum += 1
						if raidoregg.Message.GymName != nil && strings.Contains(*(raidoregg.Message.GymName), *(filters.AddNotifyRaid.GymName)) {
							boolSum += 1
						}
					}
					if filters.AddNotifyRaid.MonMovesFilters != nil {
						for _, element := range filters.AddNotifyRaid.MonMovesFilters {
							if element.Mon != nil {
								enabledSum += 1
								if element.Mon.Id == *(raidoregg.Message.PokemonID) {
									boolSum += 1
								}
							}
							if element.Move1 != nil {
								//filterMove1Enable = true
								enabledSum += 1
								for _, element := range element.Move1 {
									if raidoregg.Message.Move1 != nil && *(raidoregg.Message.Move1) == element.Id {
										//filterMove1bool = true
										boolSum += 1
									}
								}
							}
							if element.Move2 != nil {
								//filterMove2Enable = true
								enabledSum += 1
								for _, element := range element.Move2 {
									if raidoregg.Message.Move2 != nil && *(raidoregg.Message.Move2) == element.Id {
										//filterMove2bool = true
										boolSum += 1
									}
								}
							}
						}

					}
					if filters.AddNotifyRaid.Sponsor != nil {
						enabledSum += 1
						if raidoregg.Message.Sponsor != nil && *(raidoregg.Message.Sponsor) == *(filters.AddNotifyRaid.Sponsor) {
							boolSum += 1
						}
					}
					if locationbool && (boolSum == enabledSum && boolSum != 0) {
						isOkToAlter = true
					}
				}
			}
		}
	}
	return isOkToAlter
}

func (u *PokeUtility) ApplyFiltersToPokemonMessage(mon *model.PokeMinerMonMessage, filters *model.Filters) bool {
	isOkToAlter := false
	locationbool := false
	//fill the iv
	if mon.Message.Iv == nil {
		iv := ((*(mon.Message.IndividualAttack) + *(mon.Message.IndividualDefense) + *(mon.Message.IndividualStamina)) / 45) * 100
		mon.Message.Iv = &iv
		//fmt.Sprintf("%.2f",iv)
	}
	//make sure the location is within the region

	//checking to apply filter
	if filters.AddLocation != nil {
		//check to see if the mon is within the location radius
		radius := u.GreatCircleDistance(&(mon.Message.GeoLocation), &(filters.AddLocation.GeoLocation)) * KMTOMILES
		d := *(filters.AddLocation.Radius)
		if radius > d {
			return false
		} else {
			locationbool = true
			enabledSum := 0
			boolSum := 0
			/*var filteralllevelbool,filteralllevelEnabled, filterallivbool,filterallivEnabled,filterlevelbool,
			filterlevelEnabled,filterivbool,filterivEnabled,filterCpEnabled,filterCpbool,filterMove1bool,
			filterMove1Enable,filterMove2bool,filterMove2Enable,filterAttackEnable,filterAttackbool,
			filterDefenseEnable, filterDefensebool,filterStaminaEnable,filterStaminabool,filterGenderEnable,
			filterGenderbool bool*/
			//check to see if notify all is up
			if filters.AddNotifyAll != nil {
				//check to see if level, iv
				if filters.AddNotifyAll.Level != nil {
					//filteralllevelEnabled = true
					enabledSum += 1
					if (mon.Message.PokemonLevel != nil) && (*(mon.Message.PokemonLevel) <= filters.AddNotifyAll.Level.Max && *(mon.Message.PokemonLevel) >= filters.AddNotifyAll.Level.Min) {
						//filteralllevelbool = true
						boolSum += 1
					}
				}
				if filters.AddNotifyAll.Iv != nil {
					if mon.Message.Iv != nil && *(mon.Message.Iv) != 0 {
						//filterallivEnabled = true
						enabledSum += 1
						if (mon.Message.Iv != nil) && (*(mon.Message.Iv) <= filters.AddNotifyAll.Iv.Max && *(mon.Message.Iv) >= filters.AddNotifyAll.Iv.Min) {
							//filterallivbool = true
							boolSum += 1
						}
					}
				}
			}
			if filters.AddNotifies != nil {
				for i, element := range filters.AddNotifies {
					//find the pokemon according to to filter
					if filters.AddNotifies[i].Mon.Id == *(mon.Message.PokemonID) {
						if element.Iv != nil {
							//filterivEnabled = true
							//if the message mon id within the addnotify mon id
							enabledSum += 1

							if *(mon.Message.Iv) <= element.Iv.Max && *(mon.Message.Iv) >= element.Iv.Min {
								//filterivbool = true
								boolSum += 1
							}
						}
						if element.Cp != nil {
							//filterCpEnabled = true
							enabledSum += 1
							if (mon.Message.Cp != nil) && (*(mon.Message.Cp) <= element.Cp.Max && *(mon.Message.Cp) >= element.Cp.Min) {
								//filterCpbool = true
								boolSum += 1
							}
						}
						if element.Level != nil {
							//filterlevelEnabled = true
							enabledSum += 1
							if (mon.Message.PokemonLevel != nil) && (*(mon.Message.PokemonLevel) <= element.Level.Max && *(mon.Message.PokemonLevel) >= element.Level.Min) {
								//filterlevelbool = true
								boolSum += 1
							}
						}
						if element.Move1 != nil {
							//filterMove1Enable = true
							enabledSum += 1
							for _, element := range element.Move1 {
								if mon.Message.Move1 != nil && *(mon.Message.Move1) == element.Id {
									//filterMove1bool = true
									boolSum += 1
								}
							}
						}
						if element.Move2 != nil {
							//filterMove2Enable = true
							enabledSum += 1
							for _, element := range element.Move2 {
								if mon.Message.Move2 != nil && *(mon.Message.Move2) == element.Id {
									//filterMove2bool = true
									boolSum += 1
								}
							}
						}
						if element.IndividualAttack != nil {
							//filterAttackEnable = true
							enabledSum += 1
							if mon.Message.IndividualAttack != nil && (*(mon.Message.IndividualAttack) <= element.IndividualAttack.Max && *(mon.Message.IndividualAttack) >= element.IndividualAttack.Min) {
								//filterAttackbool = true
								boolSum += 1
							}
						}
						if element.IndividualDefense != nil {
							//filterDefenseEnable = true
							enabledSum += 1
							if mon.Message.IndividualDefense != nil && (*(mon.Message.IndividualDefense) <= element.IndividualDefense.Max && *(mon.Message.IndividualDefense) >= element.IndividualDefense.Min) {
								//filterDefensebool = true
								boolSum += 1
							}
						}
						if element.IndividualStamina != nil {
							//filterStaminaEnable = true
							enabledSum += 1
							if mon.Message.IndividualStamina != nil && (*(mon.Message.IndividualStamina) <= element.IndividualStamina.Max && *(mon.Message.IndividualStamina) >= element.IndividualStamina.Min) {
								//filterStaminabool = true
								boolSum += 1
							}
						}
						if element.Gender != nil {
							//filterGenderEnable = true
							enabledSum += 1
							if mon.Message.Gender != nil && *(mon.Message.Gender) == element.Gender.Id {
								//filterGenderbool = true
								boolSum += 1
							}
						}
					}
				}
			}
			if locationbool && (boolSum == enabledSum && boolSum != 0) {
				isOkToAlter = true
			}
		}
	}
	return isOkToAlter
}

func (u *PokeUtility) ApplyFiltersToMessage(message interface{}, filters *model.Filters) bool {
	//check what kind of message this is
	isOkToAlter := false
	switch t := message.(type) {
	case *model.PokeMinerMonMessage:
		mon := message.(*model.PokeMinerMonMessage)
		isOkToAlter = u.ApplyFiltersToPokemonMessage(mon, filters)
		//fmt.Print("I am pokemessage", mon.Message.Latitude)
	case *model.PokeMinerRaidMessage:
		raidoregg := message.(*model.PokeMinerRaidMessage)
		isOkToAlter = u.ApplyFiltersToRaidOrEggMessage(raidoregg, filters)
	case *model.PokeMinerGymMessage:
		gym := message.(*model.PokeMinerGymMessage)
		isOkToAlter = u.ApplyFiltersToGymMessage(gym, filters)
	default:
		_ = t
		fmt.Print("nothing")
	}
	return isOkToAlter
}
