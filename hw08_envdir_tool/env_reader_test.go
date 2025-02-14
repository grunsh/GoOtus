package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Создаем временную директорию
	tmpDir := t.TempDir()

	// Создаем временные файлы в директории
	validFile, err := os.CreateTemp(tmpDir, "VALID_*.env")
	require.NoError(t, err)
	defer validFile.Close()

	_, err = validFile.WriteString("TEST_VALUE")
	require.NoError(t, err)

	emptyFile, err := os.CreateTemp(tmpDir, "EMPTY_*.env")
	require.NoError(t, err)
	defer emptyFile.Close()

	invalidFile, err := os.CreateTemp(tmpDir, "INVALID=_*.env")
	require.NoError(t, err)
	defer invalidFile.Close()

	// Подкаталог, чтобы проверить, что игнорируется
	subDirPath := tmpDir + "/subdir"
	err = os.Mkdir(subDirPath, 0o755)
	require.NoError(t, err)

	t.Run("Успешное чтение каталога", func(t *testing.T) {
		env, err := ReadDir(tmpDir)
		require.NoError(t, err)

		validFileName := validFile.Name()[len(tmpDir)+1:] // Получаем имя файла без пути
		require.Equal(t, "TEST_VALUE", env[validFileName].Value)
		require.False(t, env[validFileName].NeedRemove)

		emptyFileName := emptyFile.Name()[len(tmpDir)+1:] // Получаем имя файла без пути
		require.Equal(t, "", env[emptyFileName].Value)
		require.True(t, env[emptyFileName].NeedRemove)

		_, exists := env["subdir"]
		require.False(t, exists)
	})

	// Тест 2: Обработка ошибки при чтении некорректного файла
	t.Run("Обработка ошибки при чтении некорректного файла", func(t *testing.T) {
		env, err := ReadDir(tmpDir)
		require.NoError(t, err, t.Name())

		invalidFileName := invalidFile.Name()[len(tmpDir)+1:] // Получаем имя файла без пути
		_, exists := env[invalidFileName]
		require.False(t, exists)
	})

	// Тест 3: Обработка ошибки при чтении несуществующей директории
	t.Run("Обработка ошибки при чтении несуществующей директории", func(t *testing.T) {
		_, err := ReadDir("/nonexistent/directory")
		require.Error(t, err)
	})
}
