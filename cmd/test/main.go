package main

import (
	"fmt"
	"os"
)

func main() {
	entries, err := os.ReadDir("/tmp/fffffffffffff/test-cls")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, entry := range entries {
		fmt.Println(entry.Info())
	}
}
