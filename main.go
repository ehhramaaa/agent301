package main

import (
	"agent301/core"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {
	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config.yml")
	if err != nil {
		panic(err)
	}

	queryPath := config.String("query-file")
	apiUrl := config.String("bot.api-url")
	referUrl := config.String("bot.refer-url")
	refId := config.String("bot.ref-Id")
	thread := config.Int("thread")

	core.ProcessAccount(thread, queryPath, apiUrl, referUrl, refId)
}
