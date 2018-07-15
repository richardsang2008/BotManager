package controller

import (
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/utility"
	"encoding/json"
	"fmt"
)

func FilterPokeMinerInputForAllUsers(data []byte, usersfilters []string, genfence_zones []model.GeoFences, pokemonMap map[int]string, moveMap map[int]string, teamsMap map[int]string) {
	u := utility.PokeUtility{}
	//input pokeminer message
	inputData, isWithinTime, regionstr, err := u.ParsePokeMinerInput(data, genfence_zones, true)
	if regionstr == nil || *regionstr == ""{
		utility.MLog.Debug("RegionName can not be determined or not within the region so no filter")
		return
	}
	//make sure the data disappear time is > now
	if isWithinTime != nil && *isWithinTime == false {
		return
	}
	if err != nil {
		utility.MLog.Error(err)
	}
	for _, userfiltersline := range usersfilters {
		go tryHandleMsgAndFilters(inputData, userfiltersline, regionstr, pokemonMap, moveMap, teamsMap)

	}
}

func tryHandleMsgAndFilters(inputData interface{}, userfiltersline string, regionstr *string, pokemonMap map[int]string, moveMap map[int]string, teamsMap map[int]string) {
	//load the userfilters into filters
	u := utility.PokeUtility{}
	filters := &model.Filters{}
	userfilters := []byte(userfiltersline)
	if err := json.Unmarshal(userfilters, &filters); err != nil {
		utility.MLog.Error(err)
	}
	filters = u.BackFillIdForFilters(filters, pokemonMap, moveMap, teamsMap)
	//now let's try to filter it. if the region does not match do not filter
	if filters.UserRegion.Region == *(regionstr) {
		isOkToNotify := u.ApplyFiltersToMessage(inputData, filters)
		if isOkToNotify {
			msg := fmt.Sprintf("Sending message %s  to user %s", inputData, filters.UserRegion.UserId)
			utility.MLog.Debug(msg)
		} else {
			msg := fmt.Sprintf("No message to user %s", filters.UserRegion.UserId)
			utility.MLog.Debug(msg)
		}

	} else {
		utility.MLog.Debug("RegionName not match so no data will be sent")
	}
}
