package controller

import (
	"encoding/json"
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

func FilterPokeMinerInput(data []byte, userfilters []byte, genfence_zones []model.GeoFences, pokemonMap map[int]string, moveMap map[int]string, teamsMap map[int]string) (interface{}, bool, error) {
	u := utility.PokeUtility{}
	//input pokeminer message
	inputData, isWithinTime, regionstr,err := u.ParsePokeMinerInput(data, genfence_zones,true)
	if regionstr == nil || *regionstr == "" {
		utility.MLog.Debug("Region can not be determined or not within the region so no filter" )
		return inputData, false, nil
	}
	//make sure the data disappear time is > now
	if isWithinTime != nil && *isWithinTime == false {
		return inputData, false, nil
	}
	if err != nil {
		utility.MLog.Error(err)
	}
	//load the userfilters into filters
	filters := &model.Filters{}
	if err := json.Unmarshal(userfilters, &filters); err != nil {
		utility.MLog.Error(err)
	}
	filters = u.BackFillIdForFilters(filters, pokemonMap, moveMap, teamsMap)
	//utility.MLog.Debug(filters)
	//now let's try to filter it
	//if the region does not match do not filter
	if  filters.UserRegion.Region == *(regionstr) {
		isOkToNotify := u.ApplyFiltersToMessage(inputData, filters)
		utility.MLog.Debug(isOkToNotify)
		return inputData, isOkToNotify, nil
	} else {
		utility.MLog.Debug("Region not match so no data will be sent")
		return inputData, false, nil
	}

}
