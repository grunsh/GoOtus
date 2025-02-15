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
	// Создаем реальный CommandRunner.
	commandRunner := NewCommand(Args[2], Args[3:]...)
	d := RunCmd(Args, env, commandRunner)
	_ = d // чтоб не ругалось :)
}
