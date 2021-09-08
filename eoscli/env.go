package eoscli

import (
	"fmt"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	"github.com/urfave/cli/v2"
)

func Env(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "list env",

		Action: x,
	}
}

func EnvFunc(c *cli.Context) error {
	for _, name := range eosc_args.Envs() {
		fmt.Printf("%s = %s\n", name, eosc_args.GetDefault(name, ""))
	}
	return nil
}
