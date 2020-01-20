package main

import (
	"github.com/derekyu332/goii-examples/common/controller"
	"github.com/derekyu332/goii-examples/common/models"
	"github.com/derekyu332/goii/frame/base"
	"github.com/derekyu332/goii/frame/behaviors"
	"github.com/derekyu332/goii/frame/validators"
	"github.com/derekyu332/goii/helper/logger"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"regexp"
	"time"
)

type AccountController struct {
	controller.WebController
}

func (this *AccountController) Group() string {
	return "account"
}

func (this *AccountController) RoutesMap() []base.Route {
	return []base.Route{
		{"POST", "register/", this.ActionRegister},
		{"POST", "login/", this.ActionLogin},
	}
}

func (this *AccountController) Behaviors() []base.IActionFilter {
	return []base.IActionFilter{
		&behaviors.QueryFilter{map[string][]base.IValidator{
			"/account/register/": {
				&validators.RequiredValidator{Values: []string{"platform", "account", "password", "email"}},
				&validators.StringValidator{Values: []string{"account", "password"}, Min: 6, Max: 20},
				&validators.StringValidator{Values: []string{"email"}, Min: 6, Max: 18},
			},
			"/account/login/": {
				&validators.RequiredValidator{Values: []string{"platform", "account", "password"}},
				&validators.StringValidator{Values: []string{"account", "password"}, Min: 6, Max: 20},
			},
		}},
		&behaviors.ActionReporter{},
	}
}

func (this *AccountController) ActionRegister(c *gin.Context) map[string]interface{} {
	platform := this.GetOrPostInt(c, "platform")
	account := this.GetOrPost(c, "account")
	password := this.GetOrPost(c, "password")
	email := this.GetOrPost(c, "email")
	pattern := `^\w+$`
	reg := regexp.MustCompile(pattern)

	if !reg.MatchString(account) {
		return this.ErrorOutput(100)
	}

	accountModel := this.NewAccountModel()
	_, err := accountModel.FindOne(bson.M{"platform": platform, "account": account})

	if err != nil {
		return this.ErrorOutput(106)
	} else if accountModel.Exists {
		return this.ErrorOutput(1000)
	}

	accountData := accountModel.Document()

	switch platform {
	case models.PLATFORM_PASSWORD:
		accountData.Account = account
		accountData.Platform = platform
		accountData.Email = email
		accountModel.SetPassword(password)
		logger.Info("[%v] Create new account %v", this.RequestID,
			accountModel.Data)
		var max_uid int64
		max_uid, err = accountModel.GetAutoIncrement()

		if err != nil {
			return this.ErrorOutput(106)
		}

		accountData.Uid = models.ACCOUNT_BASE_UID + max_uid
		accountData.Status = models.ACCOUNT_STATUS_REGISTER

		if err = accountModel.Save(); err != nil {
			return this.ErrorOutput(106)
		}

	default:
		return this.ErrorOutput(101)
	}

	accountModel.SetScenario(accountModel.Scenarios(models.ACCOUNT_SCENARIO_REGISTER))

	return gin.H{"ret": 0, "data": accountModel.Fields()}
}

func (this *AccountController) ActionLogin(c *gin.Context) map[string]interface{} {
	platform := this.GetOrPostInt(c, "platform")
	account := this.GetOrPost(c, "account")
	password := this.GetOrPost(c, "password")
	accountModel := this.NewAccountModel()
	_, err := accountModel.FindOne(bson.M{"platform": platform, "account": account})

	if err != nil {
		return this.ErrorOutput(106)
	} else if !accountModel.Exists {
		return this.ErrorOutput(1009)
	}

	switch platform {
	case models.PLATFORM_PASSWORD:
		if !accountModel.ValidatePassword(password) {
			return this.ErrorOutput(1004)
		}

	default:
		return this.ErrorOutput(101)
	}

	accountModel.Document().LastLogin = time.Now().Unix()
	accountModel.Document().LoginIp = c.ClientIP()

	if err = accountModel.Save(); err != nil {
		return this.ErrorOutput(106)
	}

	accountModel.SetScenario(accountModel.Scenarios(models.ACCOUNT_SCENARIO_LOGIN))

	return gin.H{"ret": 0, "data": accountModel.Fields()}
}
