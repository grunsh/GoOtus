package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

//func progressBar(current, total, width int) string {
//	if total == 0 {
//		return strings.Repeat("=", width) // Если total == 0, показываем полный прогресс
//	}
//	progress := current * width / total
//	return strings.Repeat("=", progress)
//}

func main() {
	flag.Parse()

	from = ".\\testdata\\out_offset0_limit0.txt"
	to = "c:\\temp\\123456789"

	fmt.Println(from)
	fmt.Println(to)
	err := Copy(from, to, 0, 0)
	if err != nil {
		fmt.Println(err)
	}

	// Place your code here.
	//total := 237   // Любое значение, не только 100
	//barWidth := 50 // Ширина прогресс-бара в символах

	//for i := 0; i <= total; i++ {
	//	// Очищаем строку и возвращаем каретку в начало
	//	percent := float64(i) / float64(total) * 100
	//	fmt.Printf("\r[%-*s] %d/%d (%.2f%%)", barWidth, progressBar(i, total, barWidth), i, total, percent)
	//	time.Sleep(30 * time.Millisecond) // Имитация работы
	//}
	//fmt.Println() // Переход на новую строку после завершения

	//var limit int64 = 1024 * 1024 * 1024 * 500
	//
	//// we will copy 500 MiB from /dev/rand to /dev/null
	//reader := io.LimitReader(rand.Reader, limit)
	//writer := ioutil.Discard
	//
	//// start new bar
	//bar := pb.Full.Start64(limit)
	//
	//// create proxy reader
	//barReader := bar.NewProxyReader(reader)
	//
	//// copy from proxy reader
	//io.Copy(writer, barReader)
	//
	//// finish bar
	//bar.Finish()
}
