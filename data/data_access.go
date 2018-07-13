package data

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/utility"
	"github.com/pkg/errors"
	"strings"
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
	DataBase.LogMode(true)
	DataBase.SetLogger(utility.MLog.GetLogger())
	//create tables
	DataBase.AutoMigrate(&model.PogoAccount{}, &model.SlackDBMessage{}, &model.SlackDBUserFilter{},
	&model.SlackDBUser{}, &model.SlackRegion{})
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == "pogo_accounts" {
			defaultTableName = "account"
		}
		return defaultTableName
	}
	if !DataBase.HasTable(&model.PogoAccount{}) {
		DataBase.CreateTable(&model.PogoAccount{})
	}
	if !DataBase.HasTable(&model.SlackDBMessage{}) {
		DataBase.CreateTable(&model.SlackDBMessage{})
	}
	if !DataBase.HasTable(&model.SlackDBUserFilter{}) {
		DataBase.CreateTable(&model.SlackDBUserFilter{})
	}
	if !DataBase.HasTable(&model.SlackDBUser{}) {
		DataBase.CreateTable(&model.SlackDBUser{})
	}
	if !DataBase.HasTable(&model.SlackRegion{}){
		DataBase.CreateTable(&model.SlackRegion{})
	}
		//DataBase.Create()


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
	if err := DataBase.Where("level>=? and level<=?", minlevel, maxlevel).Find(&accounts ).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetAccountByUserName failed " + err.Error())
		return nil, err
	} else {
		utility.MLog.Debug("DataAccessLay GetAccountByUserName end")
		return &accounts, nil
	}
}
func (s *DataAccessLay) GetSlackRegions(mode string, regionid int) (*model.SlackRegion,error){
	var regions []model.SlackRegion
	//and region_id=?",mode,regionid)
	if err :=DataBase.Where("mode =? and region_id=?",mode,regionid).Find(&regions).Error;err !=nil {
		utility.MLog.Error("DataAccessLay GetSlackRegions failed" + err.Error())
		return nil, err
	} else {
		if len(regions) >0 {
			return &regions[0],nil
		} else {
			return nil,nil
		}
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
		slackMessage := model.SlackDBMessage{ChannelId: channelId, RegionId: regionId, Ts: ts}
		DataBase.Create(&slackMessage)
		utility.MLog.Debug("DataAccessLay InsertSlackMessage inserted ts is ", ts)
		utility.MLog.Debug("DataAccessLay InsertSlackMessage end")
		return nil
	}
	return nil
}
func (s *DataAccessLay) GetSlackUserFilter(userId int) (*model.SlackDBUserFilter, error) {
	var userfilters []model.SlackDBUserFilter
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
func (s *DataAccessLay) InsertSlackUserFilter(userId uint, filters string) error {
	utility.MLog.Debug("DataAccessLay InsertSlackUserFilter starting")
	var userfilters []model.SlackDBUserFilter
	if err := DataBase.Where("user_id=?", userId).Find(&userfilters).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetSlackUserFiltreByUserId failed " + err.Error())
		return err
	} else {
		//find the userfilter by userid but no data found
		if len(userfilters) == 0 {
			record := model.SlackDBUserFilter{UserId: userId, Filters: filters}
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
func (s *DataAccessLay) AddSlackDBUser(user *model.SlackUser, regionid int ) (uint ,error) {
	//if user is not found
	foundUser, err:=s.GetSlackDBUserByEmail(user.Email,regionid)
	if err != nil {
		utility.MLog.Error(err)
		return 0,err
	} else {
		if foundUser != nil {
			//not null means there are record found then check to see if needs update
			if strings.EqualFold(user.DisplayName, foundUser.Notifyname) {
				DataBase.Model(&foundUser).Updates(model.SlackDBUser{Notifyname:user.Name,Fname:user.FirstName,Lname:user.LastName})
				return foundUser.ID,nil
			}
		} else {
			record:= model.SlackDBUser{ChannelId:regionid, Referenceid:user.ReferenceID,Fname:user.FirstName,Lname:user.LastName,
			Notifyname:user.Name, StatusId:user.StatusID, Email:user.Email, Phone:user.Phone,
			Isadmin:user.IsAdmin, Isowner:user.IsOwner,	Isbot:user.IsBot, Realname:user.Name,AccessRights: model.RightsUSER.String()}
			DataBase.Create(&record)
			return record.ID,nil
		}
	}
	return 0,nil
}
func (s *DataAccessLay) GetSlackDBUserByUserId(userId string,regionId int) (*model.SlackDBUser, error) {
	utility.MLog.Debug("DataAccessLay GetSlackDBUserByUserId starting")
	var dbusers []model.SlackDBUser
	if err := DataBase.Where("reference_id=? and channel_id=?", userId, regionId).Find(&dbusers).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetSlackDBUserByUserId failed "+ err.Error())
		return nil,err
	} else {
		//find the record
		if len(dbusers) ==0{
			return nil,nil
		} else if len(dbusers) >1{
			utility.MLog.Error("DataAccessLay GetSlackDBUserByUserId failed due to more records for the id")
			return nil,errors.New("More records are returned for the same id")
		} else {
			return &dbusers[0],nil
		}
	}
}
func (s *DataAccessLay) GetSlackDBUserByEmail(email string, regionId int )(*model.SlackDBUser, error) {
	utility.MLog.Debug("DataAccessLay GetSlackDBUserByUserId starting")
	var dbusers []model.SlackDBUser
	if err := DataBase.Where("email=? and channel_id=?", email,regionId).Find(&dbusers).Error; err != nil {
		utility.MLog.Error("DataAccessLay GetSlackDBUserByUserId failed "+ err.Error())
		return nil,err
	} else {
		//find the record
		if len(dbusers) ==0{
			return nil,nil
		} else if len(dbusers) >1{
			utility.MLog.Error("DataAccessLay GetSlackDBUserByUserId failed due to more records for the email")
			return nil,errors.New("More records are returned for the same id")
		} else {
			return &dbusers[0],nil
		}
	}
}