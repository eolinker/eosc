package eoscli

import "github.com/urfave/cli/v2"

type App struct {
	app *cli.App
}

func NewApp() *App {
	return &App{app: &cli.App{
		Name:     "eosctl",
		Usage:    "eosc controller",
		Commands: make([]*cli.Command, 0, 6),
	}}
}

func (a *App) AppendCommand(cmd *cli.Command) {
	a.app.Commands = append(a.app.Commands, cmd)
}

func (a *App) Run(args []string) error {
	return a.app.Run(args)
}
