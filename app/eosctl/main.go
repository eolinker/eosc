package main

import (
	"os"

	"github.com/eolinker/eosc/eoscli"

	"github.com/eolinker/eosc/log"
)

func main() {
	app := eoscli.NewApp()
	app.AppendCommand(eoscli.Join(nil))
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}

}
