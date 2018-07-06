package data

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/utility"
)

var (
	err      error
	DataBase *gorm.DB
)

type DataAccessLay struct {
}

func (s *DataAccessLay) New(user, pass, host, dbname string) {
	utility.MLog.Info("Open database")
	con := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", user, pass, host, dbname)
	//con :=fmt.Sprintf("%v:%v@/%v?charset=utf8&parseTime=True&loc=Local", user, pass,  dbname)
	DataBase, err = gorm.Open("mysql", con)
	if err != nil {
		utility.MLog.Panic("Error creating connection pool: " + err.Error())
	}
	//create tables
	DataBase.AutoMigrate(&model.PogoAccount{}, &model.SlackMessage{}, &model.SlackUserFilter{})
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == "pogo_accounts" {
			defaultTableName = "account"
		}
		return defaultTableName
	}
	if !DataBase.HasTable(&model.PogoAccount{}) {
		DataBase.CreateTable(&model.PogoAccount{})
	}
	if !DataBase.HasTable(&model.SlackMessage{}) {
		DataBase.CreateTable(&model.SlackMessage{})
	}
	if !DataBase.HasTable(&model.SlackUserFilter{}) {
		DataBase.CreateTable(&model.SlackUserFilter{})
	}
}
func (s *DataAccessLay) Close() {
	utility.MLog.Info("Closing database")
	DataBase.Close()
}
func (s *DataAccessLay) AddAccount(account model.PogoAccount) (*string, error) {
	utility.MLog.Debug("DataAccessLay AddAccount starting")
	DataBase.Create(&account)
	utility.MLog.Debug("DataAccessLay AddAccount end")
	ret := fmt.Sprint(account.ID)
	return &ret, nil
}
func (s *DataAccessLay) GetAccount(id uint) (*[]model.PogoAccount, error) {
	utility.MLog.Debug("DataAccessLay GetAccount starting")
	var accounts []model.PogoAccount
	if err := DataBase.Where("id=?", id).Find(&accounts).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetAccount failed " + err.Error())
		return nil, err
	} else {
		utility.MLog.Debug("DataAccessLay GetAccount end")
		return &accounts, nil
	}
}
func (s *DataAccessLay) GetAccountByUserName(username string) (*[]model.PogoAccount, error) {
	utility.MLog.Debug("DataAccessLay GetAccountByUserName starting")
	var accounts []model.PogoAccount
	if err := DataBase.Where("username=?", username).Find(&accounts).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetAccountByUserName failed " + err.Error())
		return nil, err
	} else {
		utility.MLog.Debug("DataAccessLay GetAccountByUserName end")
		return &accounts, nil
	}
}
func (s *DataAccessLay) GetAccountByLevel(minlevel, maxlevel int) (*[]model.PogoAccount, error) {
	utility.MLog.Debug("DataAccessLay GetAccountByLevel starting")
	var accounts []model.PogoAccount
	if err := DataBase.Find(&accounts, "level>=? and level<=?", minlevel, maxlevel).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetAccountByUserName failed " + err.Error())
		return nil, err
	} else {
		utility.MLog.Debug("DataAccessLay GetAccountByUserName end")
		return &accounts, nil
	}
}

func (s *DataAccessLay) UpdateAccount(account model.PogoAccount) (*string, error) {
	utility.MLog.Debug("DataAccessLay UpdateAccount starting")
	DataBase.Model(&account).Updates(account)
	ret := fmt.Sprint(account.ID)
	utility.MLog.Debug("DataAccessLay UpdateAccount end")
	return &ret, nil
}
func (s *DataAccessLay) UpdateAccountSetSystemIdToNull(account model.PogoAccount) {
	DataBase.Model(&account).Update("system_id", gorm.Expr("NULL"))
}
func (s *DataAccessLay) InsertSlackMessage(regionId int, channelId string, ts float64) error {
	utility.MLog.Debug("DataAccessLay InsertSlackMessage starting")
	if ts != 0 {
		slackMessage := model.SlackMessage{ChannelId: channelId, RegionId: regionId, Ts: ts}
		DataBase.Create(&slackMessage)
		utility.MLog.Debug("DataAccessLay InsertSlackMessage inserted ts is ", ts)
		utility.MLog.Debug("DataAccessLay InsertSlackMessage end")
		return nil
	}
	return nil
}
func (s *DataAccessLay) GetSlackUserFilter(userId int) (*model.SlackUserFilter, error) {
	var userfilters []model.SlackUserFilter
	utility.MLog.Debug("DataAccessLay GetSlackUserFilter starting")
	if err := DataBase.Where("user_id=?", userId).Find(&userfilters).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetSlackUserFilter failed " + err.Error())
		return nil, err
	} else {
		if len(userfilters) > 0 {
			utility.MLog.Debug("DataAccessLay GetSlackUserFilter end")
			return &userfilters[0], nil
		} else {
			return nil, nil
		}
	}
}
func (s *DataAccessLay) InsertSlackUserFilter(userId int, filters string) error {
	utility.MLog.Debug("DataAccessLay InsertSlackUserFilter starting")
	var userfilters []model.SlackUserFilter
	if err := DataBase.Where("user_id=?", userId).Find(&userfilters).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetSlackUserFiltreByUserId failed " + err.Error())
		return err
	} else {
		//find the userfilter by userid but no data found
		if len(userfilters) == 0 {
			record := model.SlackUserFilter{UserId: userId, Filters: filters}
			DataBase.Create(&record)
			utility.MLog.Debug("DataAccessLay InsertSlackUserFilter end")
		} else {
			userfilters[0].Filters = filters
			DataBase.Model(&userfilters[0]).Update("filters", filters)
			utility.MLog.Debug("DataAccessLay InsertSlackUserFilter update end")
		}
		return nil
	}
}
