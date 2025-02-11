package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readEnvFile(filename string) (EnvValue, error) {

	var retValue EnvValue

	file, err := os.Open(filename)
	if err != nil {
		return EnvValue{
			Value:      "",
			NeedRemove: false,
		}, fmt.Errorf("Ошибка открытия файла %s переменной", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		retValue.Value = scanner.Text() // Получаем текст строки
		retValue.Value = strings.TrimRight(retValue.Value, " \t")
		retValue.Value = strings.ReplaceAll(retValue.Value, "\x00", "\n")
	}
	if err := scanner.Err(); err != nil {
		return EnvValue{
			Value:      "",
			NeedRemove: false,
		}, fmt.Errorf("Ошибка чтения переменной из файла: %s", filename, err)
	}

	if retValue.Value == "" {
		retValue.NeedRemove = true
	}

	return retValue, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {

	RetVal := make(Environment)

	files, err := os.ReadDir(dir)

	if err != nil {
		return nil, fmt.Errorf("Не удалось прочитать каталог %s файлов переменных: %w", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		env, err := readEnvFile(dir + fileName)
		if err != nil {
			return nil, fmt.Errorf("Не удалось прочитать файл %s переменных: %w", dir+fileName, err)
		}
		RetVal[fileName] = env
	}

	return RetVal, nil
}
