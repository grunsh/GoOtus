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
}
