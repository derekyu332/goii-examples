package models

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/derekyu332/goii/frame/base"
	"github.com/derekyu332/goii/frame/mongo"
	"github.com/derekyu332/goii/helper/extend"
	"github.com/globalsign/mgo/bson"
	"time"
)

const (
	ACCOUNT_SCENARIO_REGISTER = "register"
	ACCOUNT_SCENARIO_LOGIN    = "login"
	ERROR_PASSWORD_TRY_TIMES  = 10
	PASSWORD_LOCK_TIME        = 86400
	ACCOUNT_BASE_UID          = 10000 //最小的UID
	ACCOUNT_STATUS_REGISTER   = 1     //已经注册
)

const (
	PLATFORM_GUEST    = 0
	PLATFORM_PASSWORD = 1
)

type SecurityDocument struct {
	LastTry    int64 `bson:"last_try" attr:"last_try"`
	LockExpire int64 `bson:"lock_expire" attr:"lock_expire"`
	Today      int
}

type AccountCollection struct {
	Id          bson.ObjectId `bson:"_id,omitempty" attr:",omitempty"`
	CreatedAt   int64         `bson:"created_at" attr:"created_at"`
	UpdatedAt   int64         `bson:"updated_at" attr:"updated_at"`
	Platform    int
	Account     string
	Seq         int64 `attr:",omitempty"`
	Uid         int64
	Status      int
	Email       string
	Password    string           `attr:",omitempty"`
	SecurityDoc SecurityDocument `bson:"security_doc" attr:",omitempty"`
	LastLogin   int64            `bson:"last_login" attr:"last_login"`
	LoginIp     string           `bson:"login_ip" attr:",omitempty"`
	//以下字段不写入数据库
	modelProxy *ModelProxy `bson:",omitempty" attr:",omitempty"`
}

func (this *AccountCollection) TableName() string {
	return "col_example_account"
}

func (this *AccountCollection) GetId() interface{} {
	return this.Id
}

func (this *AccountCollection) OptimisticLock() string {
	return "seq"
}

func (this *AccountCollection) SetModified(now time.Time) {
	this.UpdatedAt = now.Unix()
}

func (this *AccountCollection) Attr(name string) interface{} {
	if this.modelProxy == nil {
		this.modelProxy = &ModelProxy{}
	}

	switch name {
	case "last_login":
		return extend.UnixTime2TimeStr(this.LastLogin)

	default:
		return nil
	}
}

type AccountModel struct {
	mongo.MongoModel
}

func (this *AccountModel) Scenarios(scenario string) (string, []string) {
	switch scenario {
	case base.DEFAULT_SCENARIO:
		{
			return this.Model.Scenarios(scenario)
		}

	case ACCOUNT_SCENARIO_REGISTER:
		{
			return scenario, []string{"platform", "account", "email"}
		}

	case ACCOUNT_SCENARIO_LOGIN:
		{
			return scenario, []string{"platform", "account", "uid", "status", "last_login"}
		}
	}

	return "", nil
}

func (this *AccountModel) Document() *AccountCollection {
	if this.Data == nil {
		this.Data = &AccountCollection{}
	}

	doc, ok := this.Data.(*AccountCollection)

	if !ok {
		return nil
	} else {
		return doc
	}
}

func (this *AccountModel) SetPassword(password string) error {
	accountData := this.Document()
	temp_string := accountData.Account + password
	src := sha1.Sum([]byte(temp_string))
	accountData.Password = hex.EncodeToString(src[:])

	return nil
}

func (this *AccountModel) ValidatePassword(password string) bool {
	accountData := this.Document()
	now := time.Now().Unix()

	if accountData.SecurityDoc.LockExpire != 0 &&
		now > accountData.SecurityDoc.LockExpire {
		accountData.SecurityDoc.Today = 0
		accountData.SecurityDoc.LockExpire = 0
	} else if now < accountData.SecurityDoc.LockExpire {
		return false
	}

	if accountData.SecurityDoc.Today >= ERROR_PASSWORD_TRY_TIMES {
		accountData.SecurityDoc.LockExpire = now + PASSWORD_LOCK_TIME
		this.Save()
		return false
	}

	temp_string := accountData.Account + password
	src := sha1.Sum([]byte(temp_string))

	if accountData.Password == hex.EncodeToString(src[:]) {
		accountData.SecurityDoc.Today = 0
		accountData.SecurityDoc.LockExpire = 0
		return true
	} else {
		accountData.SecurityDoc.LastTry = now
		accountData.SecurityDoc.Today += 1
		this.Save()
		return false
	}
}
