package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// MockCommandRunner — моковая реализация CommandRunner.
type MockCommandRunner struct {
	Environments map[string]string
}

func (m MockCommandRunner) Run() error {
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
	Tests := []struct {
		name          string
		envName       string
		envValue      string
		expectedName  string
		expectedValue string
		remove        bool
	}{
		{
			name:          "Однобуквенная переменная с однобуквенным значением",
			envName:       "A",
			envValue:      "A",
			expectedName:  "A",
			expectedValue: "A",
			remove:        false,
		},
		{
			name:          "После однобуквенной переменной, удалим её",
			envName:       "A",
			envValue:      "A",
			expectedName:  "",
			expectedValue: "",
			remove:        true,
		},
		{
			name:          "Переменная с подчёркиванием",
			envName:       "TEST_A",
			envValue:      "WITH_UNDERLINE",
			expectedName:  "TEST_A",
			expectedValue: "WITH_UNDERLINE",
			remove:        false,
		},
		{
			name:          "Переменная с тремя терминальными нулями подряд",
			envName:       "TRIPPLE_ZERO_BYTES",
			envValue:      "STRING WITH\x00ONE\x00TWO\x00THREE ZERO BYTES",
			expectedName:  "TRIPPLE_ZERO_BYTES",
			expectedValue: "STRING WITH\nONE\nTWO\nTHREE ZERO BYTES",
			remove:        false,
		},
	}

	mockRunner := MockCommandRunner{}
	mockRunner.SetCommand([]string{})

	// Вызываем тестируемую функцию с моковым CommandRunner.
	for _, test := range Tests {
		t.Run(test.name, func(t *testing.T) {
			env := Environment{
				test.envName: EnvValue{
					Value:      test.envValue,
					NeedRemove: test.remove,
				},
			}
			_ = RunCmd([]string{""}, env, mockRunner)
			require.Equal(t, test.expectedValue, os.Getenv(test.envName))
			if test.remove {
				require.Equal(t, test.expectedName, os.Getenv(test.envName))
			}
		})
	}
}
