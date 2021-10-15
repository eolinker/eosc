package eoscli

import (
	"fmt"

	"github.com/eolinker/eosc/env"

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
	for _, name := range env.Envs() {
		fmt.Println(env.GenEnv(name, env.GetDefault(name, "")))
	}
	return nil
}
