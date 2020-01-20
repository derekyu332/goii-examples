package models

import (
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"time"
)

type ModelProxy struct {
	RequestID int64
	Local     *goi18n.Localizer
}

func (this *ModelProxy) NewAccountModel() *AccountModel {
	newModel := &AccountModel{}
	newModel.Data = &AccountCollection{
		CreatedAt:  time.Now().Unix(),
		modelProxy: this,
	}
	newModel.RequestID = this.RequestID

	return newModel
}
