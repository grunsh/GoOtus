package main

import (
	"fmt"
	"testing"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("simple test", func(t *testing.T) {
		fmt.Println(t.Name())
	})
}
