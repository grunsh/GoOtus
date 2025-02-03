package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrGettingFileSize       = errors.New("cant get file size")
	ErrOpeningFile           = errors.New("opening file failed")
	ErrCreatiingFile         = errors.New("creating target file failed")
	ErrCopingFile            = errors.New("copy file failed")
	ErrFileSeek              = errors.New("file seek failed")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Получаем информацию о файле
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnsupportedFile, err)
	}

	// Проверяем, является ли файл устройством
	if fileInfo.Mode()&os.ModeDevice != 0 {
		return fmt.Errorf("%w: %w", ErrUnsupportedFile, err)
	}

	// Открываем файл. Выходим, если ошибка.
	inFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOpeningFile, err)
	}
	defer func() {
		er := inFile.Close()
		if er != nil {
			fmt.Printf("Error closing file: %v\n", er)
		}
	}()

	// Получаем размер файла.
	inFileInfo, err := inFile.Stat()
	if err != nil {
		return fmt.Errorf("%w (%s): %w", ErrGettingFileSize, fromPath, err)
	}
	fmt.Println("Размер входного файла: ", inFileInfo.Size())

	// Смещение больше размера файла? Выходим
	if offset > inFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Создание целевого файла
	outFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("%w(%s): %w", ErrCreatiingFile, toPath, err)
	}
	defer func() {
		er := outFile.Close()
		if er != nil {
			fmt.Printf("Error closing file: %v\n", er)
		}
	}()

	// Сместим указатель, если его задали
	if offset > 0 {
		_, err := inFile.Seek(offset, io.SeekStart)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrFileSeek, err)
		}
	}

	var Bar *pb.ProgressBar
	var WrCount int64
	var CopyErr error
	if offset+limit >= inFileInfo.Size() || (offset == 0 && limit == 0) {
		Bar = pb.Full.Start64(inFileInfo.Size())
		barReader := Bar.NewProxyReader(inFile)
		WrCount, CopyErr = io.Copy(outFile, barReader)
		if CopyErr != nil {
			return fmt.Errorf("%w (from %s to %s): %w", ErrCopingFile, inFile.Name(), toPath, CopyErr)
		}
	} else {
		fmt.Println("Для отладки: ", inFileInfo.Size(), offset, limit)
		Bar := pb.Full.Start64(limit)
		barReader := Bar.NewProxyReader(inFile)
		WrCount, CopyErr = io.CopyN(outFile, barReader, limit)
		if CopyErr != nil {
			return fmt.Errorf("%w (from %s to %s): %w", ErrCopingFile, inFile.Name(), toPath, CopyErr)
		}
	}

	fmt.Println("Скопировано: ", WrCount)
	// finish bar
	Bar.Finish()
	return nil
}
