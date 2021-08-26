package main

import (
	"os"

	"github.com/eolinker/eosc/log"
)

func main() {
	app := newApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}

}
