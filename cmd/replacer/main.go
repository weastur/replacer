package main

import (
	"fmt"

	"github.com/weastur/go-gen-replacer/internal/generator"
)

var version = "v0.0.0-dev0"

func main() {
	fmt.Println("Running my generator...")
	generator.Run()
}
