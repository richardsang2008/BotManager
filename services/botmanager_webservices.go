package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/richardsang2008/BotManager/controller"
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/utility"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddAccount(c *gin.Context) {
	var account model.PogoAccount
	c.BindJSON(&account)
	utility.MLog.Debug("Services AddAccount starting")
	newid, err := controller.AddAccount(account)
	if err != nil {
		utility.MLog.Error("Services AddAccount error " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		utility.MLog.Debug("Services AddAccount end")
		c.JSON(http.StatusOK, gin.H{"id": newid})
	}
}
func GetAccountById(c *gin.Context) {
	id := c.Params.ByName("id")
	utility.MLog.Debug("Services GetAccount starting ")
	var accounts *[]model.PogoAccount
	var err error
	u, _ := strconv.ParseUint(id, 10, 32)
	accounts, err = controller.GetAccount(uint(u))
	if err != nil {
		utility.MLog.Error("Services GetAccount error " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	} else {
		utility.MLog.Debug("Services GetAccount end ")
		c.JSON(http.StatusOK, *accounts)
	}
}
func locateAccountByIdOrUserName(id uint, username string) []uint {
	var updateAccountIDs []uint
	if id == 0 {
		//get accounts by username
		accounts, err := controller.GetAccountByUserName(username)
		if err != nil {
			utility.MLog.Error("Services UpdateAccountBySpecificFields error " + err.Error())
			return nil
		} else {
			if len(*accounts) > 0 {
				for _, ac := range *accounts {
					updateAccountIDs = append(updateAccountIDs, ac.ID)
				}
			}
		}
	} else {
		updateAccountIDs = append(updateAccountIDs, id)
	}
	return updateAccountIDs
}
func UpdateAccountBySpecificFields(c *gin.Context) {
	utility.MLog.Debug("Services UpdateAccountBySpecificFields starting ")
	var account model.PogoAccount
	c.BindJSON(&account)
	updateAccountIDs := locateAccountByIdOrUserName(account.ID, account.Username)
	if updateAccountIDs == nil {
		c.JSON(http.StatusOK, gin.H{"username": account.Username})
	} else {
		allok := true
		for _, updateAccountID := range updateAccountIDs {
			account.ID = updateAccountID
			_, err := controller.UpdateAccountBySpecialFields(account)
			if err != nil {
				utility.MLog.Error("Services UpdateAccountBySpecificFields error " + err.Error())
				allok = false
				break
			} else {
				utility.MLog.Debug("Services UpdateAccountBySpecificFields end ")
			}
		}
		if len(updateAccountIDs) == 0 {
			utility.MLog.Debug("Services UpdateAccountBySpecificFields return no records ")
			utility.MLog.Debug("Services UpdateAccountBySpecificFields end ")
			c.JSON(http.StatusOK, gin.H{"username": account.Username})
		}
		if allok == false {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "some errors"})
		} else {
			c.JSON(http.StatusOK, gin.H{"username": account.Username})
		}
	}
}
func ReleaseAccount(c *gin.Context) {
	utility.MLog.Debug("Services ReleaseAccount starting ")
	var account model.PogoAccount
	c.BindJSON(&account)
	updateAccountIDs := locateAccountByIdOrUserName(account.ID, account.Username)
	if updateAccountIDs == nil {
		c.JSON(http.StatusOK, gin.H{"username": account.Username})
	} else {
		allok := true
		for _, updateAccountID := range updateAccountIDs {
			account.ID = updateAccountID
			controller.UpdateAccountSetSystemIdToNull(account)
		}
		if len(updateAccountIDs) == 0 {
			utility.MLog.Debug("Services ReleaseAccount return no records ")
			utility.MLog.Debug("Services ReleaseAccount end ")
			c.JSON(http.StatusOK, gin.H{"username": account.Username})
		} else {
			utility.MLog.Debug("Services ReleaseAccount end ")
			if allok == false {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "some errors"})
			} else {
				c.JSON(http.StatusOK, gin.H{"username": account.Username})
			}
		}
	}
}
func AddAccountWithLevelHandler(maxlevel int) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		utility.MLog.Debug("Services AddAccountWithLevel starting ")
		level := c.Param("level")
		var account model.PogoAccount
		//set these fields so pg scout will pickup
		account.Banned = false
		account.BanFlag = false
		account.Captcha = false
		account.Shadowbanned = false
		now := time.Now()
		account.LastModified = &now
		account.Warn = false
		c.BindJSON(&account)
		//make sure the account does not exit the account
		accountIDs := locateAccountByIdOrUserName(0, account.Username)
		lvl, erro := strconv.Atoi(level)
		if erro != nil {
			lvl = 1
		}
		if len(accountIDs) == 0 {
			//good now add

			account.Level = lvl
			idstr, err := controller.AddAccount(account)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				if lvl == maxlevel {
					//release account
					controller.UpdateAccountSetSystemIdToNull(account)
				}
				c.JSON(http.StatusOK, gin.H{"id": idstr})
			}

		} else {
			account.ID = accountIDs[0]
			account.Level = lvl
			controller.UpdateAccountBySpecialFields(account)
			if lvl == maxlevel {
				//release account
				controller.UpdateAccountSetSystemIdToNull(account)
			}
			c.JSON(http.StatusOK, gin.H{"id": account.ID})
		}
		utility.MLog.Debug("Services AddAccountWithLevel end ")
	}
	return gin.HandlerFunc(fn)
}

func GetAccountByUserName(c *gin.Context) {
	utility.MLog.Debug("Services GetAccountByUserName starting ")
	query := c.Request.URL.Query()
	username := query["username"][0]
	//make sure it is a single word
	parts := strings.Split(username, " ")
	if len(parts) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad input"})

	} else {
		ids := locateAccountByIdOrUserName(0, username)
		if len(ids) == 1 {
			accounts, _ := controller.GetAccount(ids[0])
			c.JSON(http.StatusOK, (*accounts)[0])

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad input"})
		}
	}
	utility.MLog.Debug("Services GetAccountByUserName end ")
}

func GetAccountBySystemIdAndLevelAndMark(c *gin.Context) {
	utility.MLog.Debug("Services GetAccountBySystemIdAndLevel starting ")
	//getting from request.url as ?system_id='mybox'&count=1&min_level=0&max_level=29&banned_or_new=true
	query := c.Request.URL.Query()
	var systemId, countstr, minLevelstr, maxLevelstr, bannedOrNewstr string
	if query != nil {
		systemId = query["system_id"][0]
		countstr = query["count"][0]
		minLevelstr = query["min_level"][0]
		maxLevelstr = query["max_level"][0]
		bannedOrNewstr = query["banned_or_new"][0]
	}
	if systemId == "" || countstr == "" || minLevelstr == "" || maxLevelstr == "" || bannedOrNewstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
	} else {
		count, err := strconv.Atoi(countstr)
		if err != nil {
			utility.MLog.Error("input count wrong format", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{})
		} else {
			min, err := strconv.Atoi(minLevelstr)
			if err != nil {
				utility.MLog.Error("input min_level wrong format", err.Error())
				c.JSON(http.StatusBadRequest, gin.H{})
			} else {
				max, err := strconv.Atoi(maxLevelstr)
				if err != nil {
					utility.MLog.Error("input max_level wrong format", err.Error())
					c.JSON(http.StatusBadRequest, gin.H{})
				} else {
					accounts, err := controller.GetNextUseableAccountByLevel(min, max)
					if err != nil {
						utility.MLog.Error("controller.GetNextUseableAccountByLevel error ", err.Error())
						c.JSON(http.StatusBadRequest, gin.H{})
					} else {
						msg := fmt.Sprintf("System ID [%s] requested %v accounts level %d-%d from server %s", systemId, count, min, max, c.Request.Host)
						utility.MLog.Info(msg)
						if accounts == nil {
							utility.MLog.Warning("Could only deliver 0 accounts")
							c.JSON(http.StatusOK, nil)
						} else {
							if count == 1 {
								account := (*accounts)[0]
								msg := fmt.Sprintf("Event for account %v: Got assigned to [%s]", account.Username, systemId)
								utility.MLog.Info(msg)
								utility.MLog.Debug("Services GetAccountBySystemIdAndLevel end ")
								//mark this account
								account.SystemId = systemId
								controller.UpdateAccountBySpecialFields(account)
								c.JSON(http.StatusOK, account)
							} else {
								for _, account := range *accounts {
									msg := fmt.Sprintf("Event for account %v: Got assigned to [%s]", account.Username, systemId)
									utility.MLog.Info(msg)
									account.SystemId = systemId
									controller.UpdateAccountBySpecialFields(account)
								}
								utility.MLog.Debug("Services GetAccountBySystemIdAndLevel end ")
								c.JSON(http.StatusOK, *accounts)
							}
						}
					}
				}
			}

		}
	}
}
