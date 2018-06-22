package controller

import (
	"encoding/json"
	"fmt"
	"github.com/richardsang2008/BotManager/data"
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/utility"
)

var (
	Data data.DataAccessLay
)

func AddAccount(account model.PogoAccount) (*string, error) {
	utility.MLog.Debug("Controller AddAccount starting")
	newid, err := Data.AddAccount(account)
	if err != nil {
		utility.MLog.Error("Controller AddAccount error " + err.Error())
		return nil, err
	} else {
		utility.MLog.Debug("Controller AddAccount end")
		return newid, nil
	}
}
func GetAccount(id uint) (*[]model.PogoAccount, error) {
	utility.MLog.Debug("Controller GetAccount starting")
	accounts, err := Data.GetAccount(id)
	if err != nil {
		utility.MLog.Error("Controller GetAccount error " + err.Error())
		return nil, err
	} else {
		utility.MLog.Debug("Controller GetAccount end")
		return accounts, nil
	}
}
func GetAccountByUserName(username string) (*[]model.PogoAccount, error) {
	utility.MLog.Debug("Controller GetAccountByUserName starting")
	accounts, err := Data.GetAccountByUserName(username)
	if err != nil {
		utility.MLog.Error("Controller GetAccountByUserName error " + err.Error())
		return nil, err
	} else {
		utility.MLog.Debug("Controller GetAccountByUserName end")
		return accounts, nil
	}
}
func GetNextUseableAccountByLevel(minlevel, maxlevel int) (*[]model.PogoAccount, error) {
	utility.MLog.Debug("Controller GetNextUseableAccountByLevel starting")
	accounts, err := Data.GetAccountByLevel(minlevel, maxlevel)
	if err != nil {
		utility.MLog.Error("Controller GetAccountByUserName error " + err.Error())
		return nil, err
	} else {
		//filter the accounts not usable
		ret := []model.PogoAccount{}
		for _, account := range *accounts {
			if account.Banned == false && (account.SystemId == "") {
				ret = append(ret, account)
			}
		}
		utility.MLog.Debug("Controller GetAccountByUserName end")
		if len(ret) > 0 {
			return &ret, nil
		} else {
			return nil, nil
		}
	}
}
func UpdateAccountBySpecialFields(account model.PogoAccount) (*string, error) {
	utility.MLog.Debug("Controller UpdateAccountBySpecialFields starting")
	idptr, err := Data.UpdateAccount(account)
	if err != nil {
		utility.MLog.Error("Controller UpdateAccountBySpecialFields error " + err.Error())
		return nil, err
	} else {
		utility.MLog.Debug("Controller UpdateAccountBySpecialFields end")
		return idptr, nil
	}
}
func UpdateAccountSetSystemIdToNull(account model.PogoAccount) {
	Data.UpdateAccountSetSystemIdToNull(account)
}

func FilterPokeMinerInputForAllUsers(data []byte, usersfilters []string, genfence_zones []model.GeoFences, pokemonMap map[int]string, moveMap map[int]string, teamsMap map[int]string)  {
	u := utility.PokeUtility{}
	//input pokeminer message
	inputData, isWithinTime, regionstr, err := u.ParsePokeMinerInput(data, genfence_zones, true)
	if regionstr == nil || *regionstr == "" {
		utility.MLog.Debug("Region can not be determined or not within the region so no filter")
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
		utility.MLog.Debug("Region not match so no data will be sent")
	}
}
