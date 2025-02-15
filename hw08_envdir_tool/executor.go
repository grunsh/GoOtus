package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// CommandRunner — интерфейс для запуска команд.
// Весь этот балаган с интерфейсами ради того, чтобы можно было замокать cmd в тестах.
type CommandRunner interface {
	Run() error
	Wait() error
	SetStdout(stdout *os.File)
	SetStderr(stderr *os.File)
	SetStdin(stdin *os.File)
	SetCommand(cmd []string) // Добавляем метод для установки команды.
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

func (r RealCommandRunner) SetCommand(cmd []string) {
	// Устанавливаем команду и аргументы.
	r.cmd = exec.Command(cmd[0], cmd[1:]...) //nolint
}

// NewCommand создает новую команду с использованием exec.Command.
func NewCommand(name string, arg ...string) CommandRunner {
	return RealCommandRunner{cmd: exec.Command(name, arg...)}
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, commandRunner CommandRunner) (returnCode int) { //nolint
	returnCode = 0
	for k, v := range env {
		if v.NeedRemove {
			er := os.Unsetenv(k)
			if er != nil {
				returnCode = 1
			}
			continue
		}
		v.Value = strings.TrimRight(v.Value, " \t")
		v.Value = strings.ReplaceAll(v.Value, "\x00", "\n")
		er := os.Setenv(k, v.Value)
		if er != nil {
			returnCode = 1
		}
	}
	commandRunner.SetCommand(cmd)

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
			returnCode = exitError.ExitCode()
		}
	}
	return
}
