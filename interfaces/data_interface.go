package interfaces

import "github.com/richardsang2008/BotManager/model"

type DataInterface interface {
	AddAccount(account model.PogoAccount) (*string, error)
	GetAccount(id uint) (*[]model.PogoAccount, error)
	GetAccountByUserName(username string) (*[]model.PogoAccount, error)
	GetAccountByLevel(minlevel, maxlevel int) (*[]model.PogoAccount, error)
	UpdateAccount(account model.PogoAccount) (*string, error)
}
