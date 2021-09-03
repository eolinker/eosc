package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/eolinker/eosc/process"
)

func Start(name string, args []string, extra []*os.File) (*exec.Cmd, error) {

	cmd, err := process.Cmd(name, args)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	cmd.ExtraFiles = extra

	e := cmd.Start()
	if e != nil {
		log.Println(e)
		return nil, e
	}

	return cmd, nil
}
