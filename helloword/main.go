package main

import (
	"github.com/derekyu332/goii-examples/common/controller"
	"github.com/derekyu332/goii/frame"
	"github.com/derekyu332/goii/frame/base"
	"github.com/derekyu332/goii/frame/behaviors"
	"github.com/derekyu332/goii/frame/validators"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

type DemoController struct {
	controller.WebController
}

func (this *DemoController) Group() string {
	return "demo"
}

func (this *DemoController) RoutesMap() []base.Route {
	return []base.Route{
		{"GET", "index/", this.ActionIndex},
	}
}

func (this *DemoController) ActionIndex(c *gin.Context) map[string]interface{} {
	name := this.GetOrPost(c, "name")

	if name == "" {
		return this.ErrorOutput(100)
	}

	return gin.H{"ret": 0, "message": "hello " + name}
}

func (this *DemoController) Behaviors() []base.IActionFilter {
	return []base.IActionFilter{
		&behaviors.QueryFilter{map[string][]base.IValidator{
			"/demo/index/": {
				&validators.RequiredValidator{Values: []string{"name"}},
			},
		}},
	}
}

func main() {
	app := &frame.App{
		LogLevel: logging.INFO,
		WebInit: &frame.WebServerConfig{
			Address: ":80",
		},
		MessageSource: []string{"../messages/err.en.toml", "../messages/err.zh.toml"},
	}
	app.PrepareToRun()
	app.RegisterControllers([]base.IController{
		&DemoController{},
	})
	app.Run()
}
