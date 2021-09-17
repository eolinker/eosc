package eoscli

import (
	"fmt"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	"github.com/urfave/cli/v2"
)

type App struct {
	app *cli.App
}

func NewApp() *App {
	return &App{app: &cli.App{
		Name:     eosc_args.AppName(),
		Usage:    fmt.Sprintf("%s controller", eosc_args.AppName()),
		Commands: make([]*cli.Command, 0, 6),
	}}
}

func (a *App) AppendCommand(cmd ...*cli.Command) {
	a.app.Commands = append(a.app.Commands, cmd...)
}

func (a *App) Run(args []string) error {
	return a.app.Run(args)
}
