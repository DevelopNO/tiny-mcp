package main

import (
	"fmt"
	"os"
)

const version = "v0.0.1"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-version" {
		fmt.Println(version)
		return
	}
}
