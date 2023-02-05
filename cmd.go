package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/microctar/licorice/app/config"
	"github.com/microctar/licorice/app/facade"
	"github.com/microctar/licorice/app/utils"
)

func RunCMD() {

	var generator facade.Generator

	if confdir == "" {
		confdir = config.GetDefaultConfigDirectory()
	}

	enc_subcribtion, read_err := utils.ReadAll(inputfile)

	if read_err != nil {
		log.Fatal(read_err)
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

	collect_err := generator.Collect(enc_subcribtion, confdir, rule)

	if collect_err != nil {
		log.Fatal(collect_err)
	}

	data, err := generator.Export()

	if err != nil {
		log.Fatal(err)
	}

	utils.WriteContent(outputfile, data)

}
