package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return -1
	}
	commd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	envs := updateEnv(env)
	commd.Env = append(os.Environ(), envs...)
	commd.Stdin = os.Stdin
	commd.Stdout = os.Stdout
	commd.Stderr = os.Stderr
	if err := commd.Run(); err != nil {
		log.Fatal(err)
	}
	return commd.ProcessState.ExitCode()
}

func updateEnv(e Environment) []string {
	var env []string
	for k, v := range e {
		if _, ok := os.LookupEnv(k); ok {
			os.Unsetenv(k)
		}
		if !v.NeedRemove {
			s := fmt.Sprintf("%s=%s", k, v.Value)
			env = append(env, s)
		}
	}
	return env
}
