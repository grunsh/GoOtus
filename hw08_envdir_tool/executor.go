package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var exitCode int
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		os.Setenv(k, v.Value)
	}
	fmt.Println("Команда с аргументами", cmd[2], cmd[3], cmd[4])
	comd := exec.Command(cmd[2], cmd[3], cmd[4]) //nolint
	comd.Stdout = os.Stdout
	comd.Stderr = os.Stderr
	err := comd.Start()
	if err != nil {
		panic("Error starting command: " + err.Error())
	}
	err = comd.Wait()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}
	return exitCode
}
