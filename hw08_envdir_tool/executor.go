package main

import (
	"errors"
	"os"
	"os/exec"
)

// CommandRunner — интерфейс для запуска команд.
type CommandRunner interface {
	Run() error
	Wait() error
	SetStdout(stdout *os.File)
	SetStderr(stderr *os.File)
	SetStdin(stdin *os.File)
}

// RealCommandRunner — реализация CommandRunner, которая использует exec.Command.
type RealCommandRunner struct {
	cmd *exec.Cmd
}

func (r RealCommandRunner) Run() error {
	return r.cmd.Run()
}

func (r RealCommandRunner) Wait() error {
	return r.cmd.Wait()
}

func (r RealCommandRunner) SetStdout(stdout *os.File) {
	r.cmd.Stdout = stdout
}

func (r RealCommandRunner) SetStderr(stderr *os.File) {
	r.cmd.Stderr = stderr
}

func (r RealCommandRunner) SetStdin(stdin *os.File) {
	r.cmd.Stdin = stdin
}

// NewCommand создает новую команду с использованием exec.Command.
func NewCommand(name string, arg ...string) CommandRunner {
	return RealCommandRunner{cmd: exec.Command(name, arg...)}
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, commandRunner CommandRunner) (returnCode int) {
	var exitCode int
	for k, v := range env {
		if v.NeedRemove {
			er := os.Unsetenv(k)
			if er != nil {
				exitCode = 1
			}
			continue
		}
		er := os.Setenv(k, v.Value)
		if er != nil {
			exitCode = 1
		}
	}

	// Используем методы интерфейса для настройки потоков ввода-вывода.
	commandRunner.SetStdout(os.Stdout)
	commandRunner.SetStderr(os.Stderr)
	commandRunner.SetStdin(os.Stdin)

	err := commandRunner.Run()
	if err != nil {
		panic("Error starting command: " + err.Error())
	}
	err = commandRunner.Wait()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}
	return exitCode
}
