package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
		}, fmt.Errorf("ошибка открытия файла %s переменной: %w", filename, err)
	}
	defer func() { // Чтобы не должбал ворнинг при каждом коммите
		err := file.Close()
		if err != nil {
			panic("Это конечно ужасно и странно, но не удалось заклрыть файл: " + err.Error())
		}
	}()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		retValue.Value = scanner.Text() // Получаем текст строки
	}
	if err := scanner.Err(); err != nil {
		return EnvValue{
			Value:      "",
			NeedRemove: false,
		}, fmt.Errorf("ошибка чтения переменной из файла: %s: %w", filename, err)
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
	re := regexp.MustCompile(`^[^=]*$`)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать каталог %s файлов переменных: %w", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !re.MatchString(file.Name()) {
			continue
		}
		fileName := file.Name()
		env, err := readEnvFile(dir + "/" + fileName)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать файл %s переменных: %w", dir+"/"+fileName, err)
		}
		RetVal[fileName] = env
	}

	return RetVal, nil
}
