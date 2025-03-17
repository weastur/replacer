package main

import (
	"fmt"

	"github.com/weastur/replacer/internal/generator"
)

var version = "v0.0.0-dev1" //nolint:unused

func main() {
	fmt.Println("Running my generator...")
	generator.Run()
}
