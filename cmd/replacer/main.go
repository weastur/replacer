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
	flag.String(
		"config",
		"",
		"Path to the configuration file (empty by default).\n"+
			"If not provided, the generator will look for a config file (\".replacer.yml\") "+
			"in the current directory, then move up to each parent directory until it reaches the root (the directory "+
			"containing go.mod).\nIf no config file is found, the generator will do nothing and exit with 0 code",
	)

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	fmt.Println("Running my generator...")
	generator.Run()
}
