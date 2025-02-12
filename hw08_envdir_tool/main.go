package main

import (
	"fmt"
	"os"
)

func main() {
	Args := os.Args
	env, er := ReadDir(Args[1])
	if er != nil {
		fmt.Println(er)
	}

	d := RunCmd(Args, env)
	_ = d
}
