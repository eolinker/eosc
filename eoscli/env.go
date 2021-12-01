package eoscli

import (
	"fmt"

	"github.com/eolinker/eosc/env"

	"github.com/urfave/cli/v2"
)

func Env() *cli.Command {
	return &cli.Command{
		Name:   "env",
		Usage:  "list env",
		Action: EnvFunc,
	}
}

func EnvFunc(c *cli.Context) error {
	for _, name := range env.Envs() {
		fmt.Println(name, ":\t", env.GenEnv(name, env.GetDefault(name, "")))
	}
	for k, v := range env.GetConfig() {
		fmt.Println(k, ":\t", v)
	}
	return nil
}
