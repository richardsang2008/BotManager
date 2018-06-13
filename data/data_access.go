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
	DataBase, err = gorm.Open("mysql", con)
	if err != nil {
		utility.MLog.Panic("Error creating connection pool: " + err.Error())
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == "pogo_accounts" {
			defaultTableName = "account"
		}
		return defaultTableName
	}
	DataBase.CreateTable(&model.PogoAccount{}, &model.PokemonData{})
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
