package env

import (
	"bytes"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Home() string {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		panic("not support windows")
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() string {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return ""
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "/"
	}

	return result
}
