package eoscli

import (
	"github.com/urfave/cli/v2"
)

func Env() *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "list env",
		//Action: EnvFunc,
	}
}

//func EnvFunc(c *cli.Context) error {
//	for _, name := range env.Envs() {
//		log.Debug(env.GenEnv(name, env.GetDefault(name, "")))
//	}
//	for k, v := range env.GetConfig() {
//		if k == "" {
//			continue
//		}
//		log.Debug(fmt.Sprintf("%s_%s=%s", strings.ToUpper(env.AppName()), k, v))
//	}
//	return nil
//}
