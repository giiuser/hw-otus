package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	ProcSuccessful = 0
	ProcFailed     = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return ProcFailed
	}
	app, args := cmd[0], cmd[1:]
	proc := exec.Command(app, args...)

	for k, v := range env {
		if v.NeedRemove {
			if err := os.Unsetenv(k); err != nil {
				return 1
			}
			delete(env, k)
		}
	}

	proc.Env = os.Environ()
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin

	for k, v := range env {
		proc.Env = append(proc.Env, fmt.Sprintf("%s=%s", k, v.Value))
	}

	if err := proc.Run(); err != nil {
		return ProcFailed
	}

	return ProcSuccessful
}
