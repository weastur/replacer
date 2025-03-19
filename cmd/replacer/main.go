package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/weastur/replacer/internal/generator"
)

var version = "v0.0.0-dev1"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	versionFlag := flag.Bool("version", false, "Print the version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	fmt.Println("Running my generator...")
	generator.Run()
}
