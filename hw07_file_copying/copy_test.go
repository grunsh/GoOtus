package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func ReadOneString(filename string) (string, error) {
	var Line string
	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close() // Закрываем файл после завершения

	// Создаем сканер для чтения файла построчно
	scanner := bufio.NewScanner(file)

	// Читаем файл построчно
	for scanner.Scan() {
		Line = scanner.Text() // Получаем текущую строку
	}

	// Проверяем, не возникла ли ошибка при сканировании
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return Line, nil
}

func DeleteFile(filename string) error {
	// Проверяем, существует ли файл
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}

	// Удаляем файл
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func GetfileSize(filename string) (int64, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func TestCopy(t *testing.T) {
	// С помощью смещения и лимита читаем слово "Pack" из input.txt.
	t.Run("Read word \"Pack\" from input.txt", func(t *testing.T) {
		fmt.Println(t.Name())
		err := Copy("./testdata/input.txt", "out_test.txt", 13, 4)
		if err != nil {
			fmt.Println(err)
		}
		s, erSt := ReadOneString("out_test.txt")
		erDel := DeleteFile("out_test.txt")
		require.NoError(t, err)
		require.NoError(t, erSt)
		require.NoError(t, erDel)
		require.Equal(t, "Pack", s)
		require.Len(t, s, 4)
	})

	// Проверяем нулевой целевой файл.
	t.Run("Нулевой целевой файл из-за лимита и смещения", func(t *testing.T) {
		inFileSize, erInFileSize := GetfileSize("./testdata/input.txt")
		err := Copy("./testdata/input.txt", "out_test.txt", inFileSize, 0)
		if err != nil {
			fmt.Println(err)
		}
		OutFileSize, erOutFileSize := GetfileSize("out_test.txt")
		erDel := DeleteFile("out_test.txt")
		require.Zero(t, OutFileSize)
		require.NoError(t, err)
		require.NoError(t, erInFileSize)
		require.NoError(t, erOutFileSize)
		require.NoError(t, erDel)
	})
}
