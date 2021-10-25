package eoscli

import (
	"fmt"

	"github.com/eolinker/eosc/env"

	"github.com/urfave/cli/v2"
)

func Env() *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "list env",

		Action: EnvFunc,
	}
}

func EnvFunc(c *cli.Context) error {
	for _, name := range env.Envs() {
		fmt.Println(env.GenEnv(name, env.GetDefault(name, "")))
	}
	return nil
}
