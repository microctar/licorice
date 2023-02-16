package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/microctar/licorice/app/config"
	"github.com/microctar/licorice/app/facade"
	"github.com/microctar/licorice/app/utils"
)

func runCMD() {

	var generator facade.Generator

	if confdir == "" {
		confdir = config.GetDefaultConfigDirectory()
	}

	encSubscription, readErr := utils.ReadAll(inputfile)

	if readErr != nil {
		log.Fatal(readErr)
	}

	switch target {
	case "clash":
		generator = facade.NewGenerator("clash")

		if rule == "" {
			rule = config.DefaultClashRule
		} else {
			rule = fmt.Sprintf("%s/%s", config.DefaultClashConfigPath, rule)
		}

	default:
		log.Fatal(errors.New("Unknown target"))
	}

	if collectErr := generator.Collect(encSubscription, confdir, rule); collectErr != nil {
		log.Fatal(collectErr)
	}

	data, err := generator.Export()

	if err != nil {
		log.Fatal(err)
	}

	utils.WriteContent(outputfile, data)

}
