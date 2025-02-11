package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var exitCode int = 0
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		os.Setenv(k, v.Value)
	}
	comd := exec.Command(fmt.Sprintf("%s %s %s", cmd[2], cmd[3], cmd[4]))
	comd.Stdout = os.Stdout
	comd.Stderr = os.Stderr
	err := comd.Start()
	if err != nil {
	}
	err = comd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Команда завершилась с ненулевым кодом завершения
			exitCode = exitError.ExitCode()
		}
	}
	return exitCode
}
