package main

import (
	"fmt"
	"os"
)

func main() {
	env, er := ReadDir(".\\testdata\\env\\")
	if er != nil {
		fmt.Println(er)
	}

	d := RunCmd(os.Args, env)
	_ = d
}
