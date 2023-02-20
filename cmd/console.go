package cmd

import (
	"log"

	"github.com/microctar/licorice/app/facade"
	"github.com/microctar/licorice/app/utils"
)

func runConsole() {
	var generator facade.Generator

	encSubs, rerr := utils.ReadAll(inputFile)

	if rerr != nil {
		log.Fatalln(rerr)
	}

	{
		var gerr error
		switch client {
		case "clash":
			generator = facade.NewGenerator("clash")
			gerr = generator.Collect(encSubs, confDir, clashRulePath+"/"+clashRule)
		default:
			log.Panicf("unknown client: %v\n", client)
		}

		if gerr != nil {
			log.Fatalln(gerr)
		}
	}

	data, gerr := generator.Export()

	if gerr != nil {
		log.Fatal(gerr)
	}

	_, werr := utils.WriteContent(outputFile, data)

	if werr != nil {
		log.Fatal(werr)
	}

}
