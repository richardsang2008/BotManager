package controller

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/utility"
	"strconv"
	"strings"
	"sync"
	"math"
)
const MaxFloat= math.MaxFloat64
func doesNeedingHandleMessage(ev *slack.MessageEvent) bool {
	if ev.SubType != "" {
		return false
	}
	return true
	//return ev.BotID == "" && ev.SubType != "message_deleted" && ev.SubMessage == nil && ev.Hidden == false && ev.ThreadTimestamp == "" && ev.Msg.ItemType == "" && !strings.HasPrefix(ev.Msg.Text, "<")
}
type SlackController struct {
	Lisaapi       *slack.Client
	Masterapi     *slack.Client
	Botapi        *slack.Client
	SlackUtility  utility.SlackUtility
	NSQController NSQController
}

func (c *SlackController) SlackSelfHost(env string,lisaslacktoken string, masterslackToken string, botslackToken string, produceraddress string, consumeraddress string, topic string, channel string, wg *sync.WaitGroup) {
	wg.Add(1)
	//default to one mpk_region
	slackregion,_:=Data.GetSlackRegionsByModeAndRegionName(env,"mpk_region")

	lisaslacktoken =slackregion.Lisaslacktoken
	masterslackToken= slackregion.MasterslackToken
	botslackToken = slackregion.MasterslackToken
	c.Lisaapi = slack.New(lisaslacktoken)
	c.Masterapi = slack.New(masterslackToken)
	c.Botapi = slack.New(botslackToken)
	c.SlackUtility = utility.SlackUtility{}
	c.NSQController = NSQController{}
	go c.NSQController.InitNSQ(produceraddress, consumeraddress, topic, channel, wg)
	rtm := c.Lisaapi.NewRTM()
	go rtm.ManageConnection()
	utility.MLog.Info("Lisa chat box is running..")
Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				//IMMarketedEvent, GroupMarkedEvent, ChannelMarkedEvent, MessageEvent needs to be tracked so they can be deleted later
			case *slack.IMMarkedEvent:
				slackMessage := model.SlackDBMessage{RegionId: 1, ChannelId: ev.Channel}
				tsfloat, _ := strconv.ParseFloat(ev.Timestamp, 64)
				slackMessage.Ts = tsfloat
				byteArray, _ := json.Marshal(slackMessage)
				c.NSQController.ProducerPublishMessage(byteArray, topic)
			case *slack.GroupMarkedEvent:
				slackMessage := model.SlackDBMessage{RegionId: 1, ChannelId: ev.Channel}
				tsfloat, _ := strconv.ParseFloat(ev.Timestamp, 64)
				slackMessage.Ts = tsfloat
				byteArray, _ := json.Marshal(slackMessage)
				c.NSQController.ProducerPublishMessage(byteArray, topic)
			case *slack.ChannelMarkedEvent:
				slackMessage := model.SlackDBMessage{RegionId: 1, ChannelId: ev.Channel}
				tsfloat, _ := strconv.ParseFloat(ev.Timestamp, 64)
				slackMessage.Ts = tsfloat
				byteArray, _ := json.Marshal(slackMessage)
				c.NSQController.ProducerPublishMessage(byteArray, topic)
			case *slack.MessageEvent:
				if len(ev.User) == 0 {
					continue
				}
				slackMessage := model.SlackDBMessage{RegionId: slackregion.ID, ChannelId: ev.Msg.Channel}
				tsfloat, _ := strconv.ParseFloat(ev.Msg.Timestamp, 64)
				slackMessage.Ts = tsfloat
				byteArray, _ := json.Marshal(slackMessage)
				//byteArray, _:= json.Marshal(ev.Msg)
				c.NSQController.ProducerPublishMessage(byteArray, topic)

				slackUser, err := getUserInfo(ev.User, c.Lisaapi,slackregion.ID)
				//this is the user information
				if err != nil {
					utility.MLog.Error(err)
				} else {
					b, _ := json.Marshal(*slackUser)
					utility.MLog.Debug(string(b))
				}
				//if the message is not starting with ! then nothing
				if strings.HasPrefix(ev.Msg.Text, "!") {
					utility.MLog.Info("I need to do something to make this done!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1")
					ParseSlackUserInput(slackUser.DBId,ev.Msg.Text,slackregion.ID)
				}
				if doesNeedingHandleMessage(ev) {
					subType := ev.Msg.SubType
					switch subType {
					case "channel_join":
						utility.MLog.Debug("%s join channel %s", ev.Msg.User, ev.Msg.Channel)
					case "file_share":
						utility.MLog.Debug(ev.Msg.Text)
					default:
						if strings.HasPrefix(ev.Msg.Text, "You have been removed") {
							utility.MLog.Debug(ev.Msg.Text)
						} else {
							//go handleEventMessage(ev.Text, ev.User, ev.Channel)
						}
					}
				}
			case *slack.RTMError:
				utility.MLog.Error("Message: %v\n", ev.Error())

			case *slack.InvalidAuthEvent:
				utility.MLog.Debug("Invalid credentials")
				break Loop
			default:
				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
	wg.Wait()
}
func getUserInfo(userid string, api *slack.Client, regionId uint) (*model.SlackUser, error){
	c:=utility.SlackUtility{}
	user,err := c.GetUserInfo(userid,api)
	if err != nil {
		utility.MLog.Error(err)
		return nil, err
	}else {
		//check user in db
		dbuser,err:= Data.GetSlackDBUserByEmail(user.Email,regionId)
		if err!=nil {
			utility.MLog.Error(err)
			return nil,err
		} else{
			if dbuser!= nil {
				//populate the db stuff
				if !strings.EqualFold(dbuser.Notifyname, user.DisplayName) {
					id,err:=Data.AddSlackDBUser(user, regionId)
					if err!=nil {
						utility.MLog.Error(err)
						return nil, err
					}
					dbuser.ID =id
				}
				user.DBId = dbuser.ID
				return user,nil
			} else {
				//add user into db
				id,err:=Data.AddSlackDBUser(user,regionId)
				if err!=nil {
					utility.MLog.Error(err)
					return nil, err
				}
				user.AccessRights = model.RightsUSER
				user.DBId =id
				return user,nil
			}
		}
	}
}
func getParamIndex(parts []string, word string) int {
	if len(parts) > 0 {
		for i, value := range parts {
			if strings.EqualFold(value, word) {
				return i
			}
		}
	}
	return -1
}
func setUserInputSingleStringValue(parts []string, paramword string) *string {
	isOkCount := 0
	totalkeyparams := 0
	index := getParamIndex(parts, paramword)
	if index > 0 {
		totalkeyparams += 1
		//get the value which is the next on the slice
		value := parts[index+1]
		isOkCount += 1
		return &value
	}
	return nil
}
func setUserInputSingleBoolValue(parts []string, paramword string) bool {
	strv := setUserInputSingleStringValue(parts, paramword)
	if strv != nil {
		isBool, _ := strconv.ParseBool(*strv)
		return isBool
	}
	return false
}
func setUserInputSingleFloatValue(parts []string, paramword string) float64 {
	strv := setUserInputSingleStringValue(parts, paramword)
	//get the value which is the next on the slice
	if strv != nil {
		fl, err := strconv.ParseFloat(*strv, 64)
		if err != nil {
			utility.MLog.Error(err)
		}
		return fl
	}
	return 0
}
func setUserinputAsRange(parts []string, paramword string) (*model.Range, bool) {
	isOkCount := 0
	totalkeyparams := 0
	index := getParamIndex(parts, paramword)
	if index > 0 {
		totalkeyparams += 1
		value := parts[index+1]
		var err1 error
		rangeValue := model.Range{}
		if strings.HasSuffix(value, "+") {
			//get the value before the +
			charArray := []rune(value)
			len := len(charArray)
			subvalue := string(charArray[0 : len-1])
			rangeValue.Min, err1 = strconv.ParseFloat(subvalue, 64)
			if err1 != nil {
				utility.MLog.Error(err1)
			} else {
				isOkCount += 1
			}
			//rangeValue.Max = MaxFloat
			rangeValue.Max = 200000
		} else if strings.HasSuffix(value, "-") {
			charArray := []rune(value)
			len := len(charArray)
			subvalue := string(charArray[0 : len-1])
			rangeValue.Max, err1 = strconv.ParseFloat(subvalue, 64)
			if err1 != nil {
				utility.MLog.Error(err1)
			} else {
				isOkCount += 1
			}
			rangeValue.Min = 0
		} else {
			rangeValue.Min, err1 = strconv.ParseFloat(value, 64)
			if err1 != nil {
				utility.MLog.Error(err1)
			} else {
				isOkCount += 1
			}
			rangeValue.Max = rangeValue.Min
		}
		return &rangeValue, isOkCount == totalkeyparams
	}
	return nil, false
}
func getUserFiltersByUserId(userId uint) (*model.Filters, error) {
	filter, _ := Data.GetSlackUserFilter(userId)
	if (filter != nil ) {
		//turn that into the object
		var filters model.Filters
		err := json.Unmarshal([]byte(filter.Filters), &filters)
		if err != nil {
			utility.MLog.Error(err)
			return nil, err
		} else {
			return &filters, nil
		}
	}
	return nil, nil
}
func ParseSlackUserInput(userid uint,userInput string, regionid uint) {
	//parse the userinput by delimiter white space
	if strings.HasPrefix(userInput, "!") {
		parts := strings.Split(userInput, " ")
		//if the parts is more than 1, check what is it start with
		if len(parts) > 0 {
			//check the first part, and see what it is
			isUserInGoodStatus := true
			region,err:=Data.GetSlackRegionsById(regionid)
			if err != nil {
				utility.MLog.Error(err)
			}
			_region :=model.UserRegion{Region:region.RegionName}
			filters:=model.Filters{UserRegion:&_region}
			switch a := strings.ToLower(parts[0]); a {
			case "!addlocation":
				//make sure the user is in subscription status, balance is greater than 0
				if isUserInGoodStatus {
					userfilters, err := getUserFiltersByUserId(userid)
					if err != nil {
						utility.MLog.Error(err)
					} else {
						if userfilters != nil {
							//userfilter does  exit
							addlocationcmd := model.AddLocationCmd{}
							addlocationcmd.Latitude = setUserInputSingleFloatValue(parts, "lan")
							addlocationcmd.Longitude = setUserInputSingleFloatValue(parts, "lng")
							addlocationcmd.Radius = setUserInputSingleFloatValue(parts, "radius")
							if addlocationcmd.Longitude != 0 && addlocationcmd.Latitude != 0 && addlocationcmd.Radius != 0 {
								//addlocation is successful then further handle is required
								//check db to see if user already has the location alter, if not add it else update it
								//load the user filter from db
								userfilters.AddLocation=&model.AddLocation{Radius:&(addlocationcmd.Radius),
								GeoLocation:model.GeoLocation{Region:&region.RegionName,
									Longitude:addlocationcmd.Longitude, Latitude:addlocationcmd.Latitude}}
								//userfilters.AddLocation.Radius = &(addlocationcmd.Radius)
								//userfilters.AddLocation.Longitude = addlocationcmd.Longitude
								//userfilters.AddLocation.Latitude = addlocationcmd.Latitude
								//save to db
								byteArray, _ := json.Marshal(userfilters)
								Data.InsertSlackUserFilter(userid, string(byteArray))
							}
						} else {
							//userfilter does  exit
							addlocationcmd := model.AddLocationCmd{}
							addlocationcmd.Latitude = setUserInputSingleFloatValue(parts, "lan")
							addlocationcmd.Longitude = setUserInputSingleFloatValue(parts, "lng")
							addlocationcmd.Radius = setUserInputSingleFloatValue(parts, "radius")
							if addlocationcmd.Longitude != 0 && addlocationcmd.Latitude != 0 && addlocationcmd.Radius != 0 {
								geolation:= model.GeoLocation{Latitude:addlocationcmd.Latitude,
									Longitude:addlocationcmd.Longitude,Region:&region.RegionName}
								filters.AddLocation = &model.AddLocation{GeoLocation:geolation,Radius:&addlocationcmd.Radius}
								//copy filters into db filters
								byteArray, _ := json.Marshal(filters)
								Data.InsertSlackUserFilter(userid, string(byteArray))
							}
						}
					}
				}
			case "!addallmons":
				if isUserInGoodStatus {
					userfilters, err := getUserFiltersByUserId(userid)
					if err != nil {
						utility.MLog.Error(err)
					} else{
						lvlranged, _ := setUserinputAsRange(parts, "lvl")
						ivranged, _ := setUserinputAsRange(parts, "iv")
						if userfilters != nil {
							//userfilter doe exit

							userfilters.AddNotifyAll = &model.AddNotifyAll{Level:lvlranged,Iv:ivranged}
							//save to db
							byteArray, _ := json.Marshal(userfilters)
							Data.InsertSlackUserFilter(userid, string(byteArray))
						} else {
							//userfilter does exit
							filters.AddNotifyAll= &model.AddNotifyAll{Level:lvlranged,Iv:ivranged}
								//save to db
							byteArray,_:=json.Marshal(filters)
							Data.InsertSlackUserFilter(userid,string(byteArray))
						}
					}
				}
			case "!addmon":
				if isUserInGoodStatus {
					userfilters, err := getUserFiltersByUserId(userid)
					if err != nil {

					} else{
						namestr := setUserInputSingleStringValue(parts, "name")
						lvlranged, _ := setUserinputAsRange(parts, "lvl")
						ivranged, _ := setUserinputAsRange(parts, "iv")
						cpranged, _ := setUserinputAsRange(parts, "cp")
						if namestr != nil {
							//check to see if the name is in the filters already
loop:
							for index, value :=range(userfilters.AddNotifies){
								if strings.EqualFold(value.Mon.Name,*namestr){
									if ivranged != nil{
										userfilters.AddNotifies[index].Iv = ivranged
									}
									if lvlranged !=nil {
										userfilters.AddNotifies[index].Level = lvlranged
									}
									if cpranged !=nil {
										userfilters.AddNotifies[index].Cp = cpranged
									}
									break loop
								}
							}
							//didn't find anything then
							//append a new record to the end of the existing
							id :=123
							mon :=model.NameAndID{Name:*namestr,Id:id}
							newone:=model.AddNotify{Level:lvlranged,Iv:ivranged,Cp:cpranged, Mon:&mon}
							userfilters.AddNotifies=append(userfilters.AddNotifies,newone)
							//save to db
							byteArray, _ := json.Marshal(userfilters)
							Data.InsertSlackUserFilter(regionid, string(byteArray))
						}
					}
				}
			case "!addallraid":
				cmd := model.AddAllRaidCmd{}
				lvlranged, _ := setUserinputAsRange(parts, "lvl")
				cmd.Lvl = lvlranged
				sponsorValue := setUserInputSingleBoolValue(parts, "sponsored")
				cmd.Sponsored = sponsorValue
				boostedValue := setUserInputSingleBoolValue(parts, "boosted")
				cmd.Boosted = boostedValue
				teamNamestr := setUserInputSingleStringValue(parts, "team")
				if teamNamestr != nil {
					cmd.Team = *teamNamestr
				}
				gymNamestr := setUserInputSingleStringValue(parts, "gym")
				if gymNamestr != nil {
					cmd.GymName = *gymNamestr
				}
			case "!addraid":
				cmd := model.AddRaidCmd{}
				lvlranged, _ := setUserinputAsRange(parts, "lvl")
				cmd.Lvl = lvlranged
				namestr := setUserInputSingleStringValue(parts, "name")
				if namestr != nil {
					cmd.Name = *namestr
				}
				sponsorValue := setUserInputSingleBoolValue(parts, "sponsored")
				cmd.Sponsored = sponsorValue
				boostedValue := setUserInputSingleBoolValue(parts, "boosted")
				cmd.Boosted = boostedValue
				teamNamestr := setUserInputSingleStringValue(parts, "team")
				if teamNamestr != nil {
					cmd.Team = *teamNamestr
				}
				gymNamestr := setUserInputSingleStringValue(parts, "gym")
				if gymNamestr != nil {
					cmd.GymName = *gymNamestr
				}
			case "!addegg":
				cmd := model.AddEggCmd{}
				lvlranged, _ := setUserinputAsRange(parts, "lvl")
				cmd.Lvl = lvlranged
				sponsorValue := setUserInputSingleBoolValue(parts, "sponsored")
				cmd.Sponsored = sponsorValue
				boostedValue := setUserInputSingleBoolValue(parts, "boosted")
				cmd.Boosted = boostedValue
				teamNamestr := setUserInputSingleStringValue(parts, "team")
				if teamNamestr != nil {
					cmd.Team = *teamNamestr
				}
				gymNamestr := setUserInputSingleStringValue(parts, "gym")
				if gymNamestr != nil {
					cmd.GymName = *gymNamestr
				}
			case "!addgym":
				cmd := model.AddGymCmd{}
				teamNamestr := setUserInputSingleStringValue(parts, "team")
				if teamNamestr != nil {
					cmd.Team = *teamNamestr
				}
				gymNamestr := setUserInputSingleStringValue(parts, "gym")
				if gymNamestr != nil {
					cmd.GymName = *gymNamestr
				}
			case "!showlocation":
				if isUserInGoodStatus {
					slackUserFilter, err := Data.GetSlackUserFilter(userid)
					if err != nil {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err := json.Unmarshal([]byte(slackUserFilter.Filters), &filters); err != nil {
							utility.MLog.Error(err)
						}
						location, _ := json.Marshal(filters.AddLocation)
						utility.MLog.Debug(string(location))
					}
				}
			case "!showmons":
				if isUserInGoodStatus {
					slackUserFilter, err := Data.GetSlackUserFilter(userid)
					if err != nil {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err := json.Unmarshal([]byte(slackUserFilter.Filters), &filters); err != nil {
							utility.MLog.Error(err)
						}
						allfilter, _ := json.Marshal(filters.AddNotifyAll)
						utility.MLog.Debug(string(allfilter))
						filtermon, _ := json.Marshal(filters.AddNotifies)
						utility.MLog.Debug(string(filtermon))
					}
				}
			case "!showraid":
				if isUserInGoodStatus {
					slackUserFilter, err := Data.GetSlackUserFilter(userid)
					if err != nil {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err := json.Unmarshal([]byte(slackUserFilter.Filters), &filters); err != nil {
							utility.MLog.Error(err)
						}
						item, _ := json.Marshal(filters.AddNotifyRaid)
						utility.MLog.Debug(string(item))
					}
				}
			case "!showegg":
				if isUserInGoodStatus {
					slackUserFilter, err := Data.GetSlackUserFilter(userid)
					if err != nil {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err := json.Unmarshal([]byte(slackUserFilter.Filters), &filters); err != nil {
							utility.MLog.Error(err)
						}
						item, _ := json.Marshal(filters.AddNotifyEgg)
						utility.MLog.Debug(string(item))
					}
				}
			case "!showgym":
				if isUserInGoodStatus {
					slackUserFilter, err := Data.GetSlackUserFilter(userid)
					if err != nil {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err := json.Unmarshal([]byte(slackUserFilter.Filters), &filters); err != nil {
							utility.MLog.Error(err)
						}
						item, _ := json.Marshal(filters.AddNotifyGym)
						utility.MLog.Debug(string(item))
					}
				}
			case "!removelocation":
			case "!removemons":
			case "!removemon":
			case "!removeraid":
			case "!removeegg":
			case "!removegym":
			case "!balance":
			case "!status":
			case "!emaillookup":
			case "!paid":
			case "!minuspay":

			default:
				//user does not need to be in the subscription, and it is ok to talk.
			}

		}
	}
}

/*
func handleEventMessage(text string, user string, channel string) {
	//get data from  the memeory cache
	if channel != lisaChannel || (strings.HasPrefix(text, ":") && strings.HasSuffix(text, ":")) {
		return
	}
	pokemonEn := model.PokemonEn{}
	if record, found := localcache.Get("pokemonEn"); found {
		pokemonEn = record.(model.PokemonEn)
	} else {
		pokemonEn1, err := slackUtility.LoadPokemonEn()
		if err != nil {
			log.Println(err)
		}
		pokemonEn = *pokemonEn1
		localcache.Set("pokemonEn", pokemonEn, cache.DefaultExpiration)
	}
	text = strings.TrimSpace(text)
	splittedStringParts := strings.Split(text, " ")
	size := len(splittedStringParts)
	botMessage := model.BotMessage{}
	botMessage.Message = "I don't understand"
	botMessage.UserID = user
	botMessage.ChannelID = channel
	if size == 0 {
	} else if size == 1 {
		if containsPokemonName(&pokemonEn, text) {
			botMessage.Message = fmt.Sprintf("find a %s", text)
		} else {
			wordrule := model.WordRules{}
			answer := wordrule.Answer_single_preset_question(text)
			botMessage.Message = answer
		}
		if botMessage.Message != "" {

		}
	} else {
		parsedInputDictionary := lisawordParser.Parse_user_input(text)
		if parsedInputDictionary == nil {

		}
		for k, v := range parsedInputDictionary {
			botMessage.Message = fmt.Sprintf("find %s %s", k, v.Param)
		}
	}
	go slackUtility.SendMessage(botMessage, lisaapi)
}*/
