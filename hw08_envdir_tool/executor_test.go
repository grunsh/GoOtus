package main

import (
	"os"
	"testing"
)

// MockCommandRunner — моковая реализация CommandRunner.
type MockCommandRunner struct{}

func (m MockCommandRunner) Run() error {
	// Выводим значение переменной окружения Test12.
	println("Test12:", os.Getenv("Test12"))
	return nil
}

func (m MockCommandRunner) Wait() error {
	return nil
}

func (m MockCommandRunner) SetStdout(stdout *os.File) {
	_ = stdout.Name()
}

func (m MockCommandRunner) SetStderr(stderr *os.File) {
	_ = stderr.Name()
}

func (m MockCommandRunner) SetStdin(stdin *os.File) {
	_ = stdin.Name()
}

func (m MockCommandRunner) SetCommand(cmd []string) {
	// Ничего не делаем, так как это мок.
	_ = cmd
}

func TestRunCmd(t *testing.T) {
	// Устанавливаем переменную окружения Test12.
	env := Environment{
		"Test12": EnvValue{
			Value:      "MockedValue",
			NeedRemove: false,
		},
	}

	// Создаем моковый CommandRunner.
	mockRunner := MockCommandRunner{}
	mockRunner.SetCommand([]string{})

	// Вызываем тестируемую функцию с моковым CommandRunner.
	exitCode := RunCmd([]string{"dummy", "arg1", "echo", "hello"}, env, mockRunner)

	// Проверяем, что функция завершилась с кодом 0.
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}
