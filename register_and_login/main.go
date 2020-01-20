package main

import (
	"github.com/derekyu332/goii/frame"
	"github.com/derekyu332/goii/frame/base"
	"github.com/op/go-logging"
)

func main() {
	app := &frame.App{
		LogLevel: logging.INFO,
		MongoInit: &frame.MongoConfig{
			Url:    "input mongo url here",
			DbName: "input mongo dbname here",
		},
		WebInit: &frame.WebServerConfig{
			Address: ":80",
		},
		MessageSource: []string{"../messages/err.en.toml", "../messages/err.zh.toml"},
	}
	app.PrepareToRun()
	app.RegisterControllers([]base.IController{
		&AccountController{},
	})
	app.Run()
}
