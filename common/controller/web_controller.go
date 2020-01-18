package controller

import (
	"fmt"
	"github.com/derekyu332/goii/frame/base"
	"github.com/gin-gonic/gin"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/yuxinyun/goii/frame/i18n"
)

type WebController struct {
	base.WebController
	localizer *goi18n.Localizer
}

func (this *WebController) PreparedForUse(c *gin.Context) {
	this.WebController.PreparedForUse(c)
	this.localizer = i18n.NewLocalizer(c, nil)
}

func (this *WebController) ErrorOutput(errorno int) map[string]interface{} {
	message := i18n.L(this.localizer, &goi18n.LocalizeConfig{
		MessageID: fmt.Sprintf("%v", errorno),
	})

	if message == "" {
		message = fmt.Sprintf("Unknown error %v", errorno)
	}

	return gin.H{"ret": errorno, "message": message}
}

func (this *WebController) ErrorCustomOutput(errorno int, errormsg string) map[string]interface{} {
	return gin.H{"ret": errorno, "message": errormsg}
}
