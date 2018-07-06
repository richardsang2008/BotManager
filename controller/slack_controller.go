package controller

import (
	"strings"
	"github.com/richardsang2008/BotManager/utility"
	"github.com/nlopes/slack"
	"encoding/json"
	"sync"
	"github.com/richardsang2008/BotManager/model"
	"strconv"
)
var (

	// Create a cache with a default expiration time of 12 hours, and which
	// purges expired items every 24 hours
	//localcache     = cache.New(12*time.Hour, 24*time.Hour)
	//lisawordParser = utility.LisaWordParser{}
	//messages       = make(chan slack.Msg, 2000)
	//message = utility.MNSQUtility.
	//lisaChannel    = "D5SSTG73R"
	//configuration  = model.Configuration{}


)
func doesNeedingHandleMessage(ev *slack.MessageEvent) bool {
	if (ev.SubType !=""){
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
func (c *SlackController) SlackSelfHost(lisaslacktoken string, masterslackToken string, botslackToken string, produceraddress string, consumeraddress string, topic string, channel string, wg *sync.WaitGroup) {
	wg.Add(1)
	c.Lisaapi          = slack.New(lisaslacktoken)
	c.Masterapi        = slack.New(masterslackToken)
	c.Botapi           = slack.New(botslackToken)
	c.SlackUtility    =  utility.SlackUtility{}
	c.NSQController = NSQController{}
	go c.NSQController.InitNSQ(produceraddress,consumeraddress, topic, channel, wg)
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
				slackMessage:=model.SlackMessage{RegionId:1,ChannelId:ev.Channel}
				tsfloat,_:=strconv.ParseFloat(ev.Timestamp,64)
				slackMessage.Ts = tsfloat
				byteArray, _:= json.Marshal(slackMessage)
				c.NSQController.ProducerPublishMessage(byteArray,topic)
			case *slack.GroupMarkedEvent:
				slackMessage:=model.SlackMessage{RegionId:1,ChannelId:ev.Channel}
				tsfloat,_:=strconv.ParseFloat(ev.Timestamp,64)
				slackMessage.Ts = tsfloat
				byteArray, _:= json.Marshal(slackMessage)
				c.NSQController.ProducerPublishMessage(byteArray,topic)
			case *slack.ChannelMarkedEvent:
				slackMessage:=model.SlackMessage{RegionId:1,ChannelId:ev.Channel}
				tsfloat,_:=strconv.ParseFloat(ev.Timestamp,64)
				slackMessage.Ts = tsfloat
				byteArray, _:= json.Marshal(slackMessage)
				c.NSQController.ProducerPublishMessage(byteArray,topic)
			case *slack.MessageEvent:
				if len(ev.User) ==0 {
					continue
				}

				slackMessage:=model.SlackMessage{RegionId:1,ChannelId:ev.Msg.Channel}
				tsfloat,_:=strconv.ParseFloat(ev.Msg.Timestamp,64)
				slackMessage.Ts = tsfloat
				byteArray, _:= json.Marshal(slackMessage)
				//byteArray, _:= json.Marshal(ev.Msg)
				c.NSQController.ProducerPublishMessage(byteArray,topic)
				slackUser,err:=c.SlackUtility.GetUserInfo(ev.User,c.Lisaapi)
				//this is the user information
				if err != nil {
					utility.MLog.Error(err)
				} else {
					b,_:=json.Marshal(*slackUser)
					utility.MLog.Debug(string(b))
				}
				//if the message is not starting with ! then nothing
				if strings.HasPrefix(ev.Msg.Text,"!") {
					utility.MLog.Info("I need to do something to make this done!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1")
					ParseSlackUserInput(ev.Msg.Text)
				}
				//b, err := json.Marshal(ev)
				//if err != nil {
				//	log.Println(err)
				//	return
				//}
				//str := fmt.Sprintf("%s", b)
				//utility.MLog.Debug(str)
				//getuserinfo from slack
				//slackUserInfo,err:=controller.GetUserInfo(ev.User,lisaapi,configuration.SlackTeams.ChannelKey)
				//if err !=nil{
				//	utility.MLog.Error(err)
				//}
				//log.Println(slackUserInfo)
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
func getParamIndex(parts []string, word string) int {
	if len(parts)>0{
		for i,value :=range(parts) {
			if strings.EqualFold(value,word){
				return i
			}
		}
	}
	return -1
}
func ParseSlackUserInput(userInput string) {
	//parse the userinput by delimiter white space
	if strings.HasPrefix(userInput,"!") {
		parts:=strings.Split(userInput," ")
		//if the parts is more than 1, check what is it start with
		if len(parts) >0 {
			//check the first part, and see what it is
			isUserInGoodStatus:=true
			switch a:=strings.ToLower(parts[0]); a{

			case "!addlocation":
				//make sure the user is in subscription status, balance is greater than 0
				if isUserInGoodStatus {
					addlocationcmd := model.AddLocationCmd{}
					isOkCount :=0
					index:=getParamIndex(parts,"lan")
					if index >0 {
						//get the value which is the next on the slice
						fl,err := strconv.ParseFloat(parts[index+1],64)
						if err != nil  {
							utility.MLog.Error(err)
						}else {
							addlocationcmd.Latitude = fl
							isOkCount+=1
						}
					}
					index=getParamIndex(parts,"lng")
					if index >0 {
						//get the value which is the next on the slice
						fl,err := strconv.ParseFloat(parts[index+1],64)
						if err != nil  {
							utility.MLog.Error(err)
						}else {
							addlocationcmd.Longitude = fl
							isOkCount+=1
						}
					}
					index=getParamIndex(parts, "radius")
					if index >0 {
						//get the value which is the next on the slice
						fl,err := strconv.ParseFloat(parts[index+1],64)
						if err != nil  {
							utility.MLog.Error(err)
						}else {
							addlocationcmd.Radius = fl
							isOkCount+=1
						}
					}
					if isOkCount ==3 {
						//addlocation is successful then further handle is required
						//check db to see if user already has the location alter, if not add it else update it
						//load the user filter from db
						userfilters:=model.Filters{}
						userfilters.AddLocation = &model.AddLocation{}
						userfilters.AddLocation.Radius =&(addlocationcmd.Radius)
						userfilters.AddLocation.Longitude = addlocationcmd.Longitude
						userfilters.AddLocation.Latitude = addlocationcmd.Latitude
						byteArray,_:= json.Marshal(userfilters)
						//save to db
						Data.InsertSlackUserFilter(1,string(byteArray))
					}
				}

			case "!addallmons":
			case "!addmon":
			case "!addallraid":
			case "!addraid":
			case "!addegg":
			case "!addgym":
			case "!showlocation":
				if isUserInGoodStatus {
					slackUserFilter,err:=Data.GetSlackUserFilter(1)
					if err != nil  {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err:=json.Unmarshal([]byte(slackUserFilter.Filters),&filters); err !=nil {
							utility.MLog.Error(err)
						}
						location,_:=json.Marshal(filters.AddLocation)
						utility.MLog.Debug(string(location))
					}
				}
			case "!showmons":
				if isUserInGoodStatus {
					slackUserFilter,err:=Data.GetSlackUserFilter(1)
					if err != nil  {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err:=json.Unmarshal([]byte(slackUserFilter.Filters),&filters); err !=nil {
							utility.MLog.Error(err)
						}
						allfilter,_:=json.Marshal(filters.AddNotifyAll)
						utility.MLog.Debug(string(allfilter))
						filtermon,_:=json.Marshal(filters.AddNotifies)
						utility.MLog.Debug(string(filtermon))
					}
				}
			case "!showraid":
				if isUserInGoodStatus {
					slackUserFilter,err:=Data.GetSlackUserFilter(1)
					if err != nil  {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err:=json.Unmarshal([]byte(slackUserFilter.Filters),&filters); err !=nil {
							utility.MLog.Error(err)
						}
						item,_:=json.Marshal(filters.AddNotifyRaid)
						utility.MLog.Debug(string(item))
					}
				}
			case "!showegg":
				if isUserInGoodStatus {
					slackUserFilter,err:=Data.GetSlackUserFilter(1)
					if err != nil  {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err:=json.Unmarshal([]byte(slackUserFilter.Filters),&filters); err !=nil {
							utility.MLog.Error(err)
						}
						item,_:=json.Marshal(filters.AddNotifyEgg)
						utility.MLog.Debug(string(item))
					}
				}
			case "!showgym":
				if isUserInGoodStatus {
					slackUserFilter,err:=Data.GetSlackUserFilter(1)
					if err != nil  {
						utility.MLog.Error(err)
					} else {
						var filters model.Filters
						if err:=json.Unmarshal([]byte(slackUserFilter.Filters),&filters); err !=nil {
							utility.MLog.Error(err)
						}
						item,_:=json.Marshal(filters.AddNotifyGym)
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