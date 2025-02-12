package main

import (
	"fmt"
	"testing"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	t.Run("base test", func(t *testing.T) {
		fmt.Println(t.Name())
	})
}
