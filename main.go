package main

import (
	"agent301/core"
	"agent301/helper"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {

	helper.PrintLogo()

	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config.yml")
	if err != nil {
		panic(err)
	}

	core.LaunchBot()
}
