package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/weastur/replacer/internal/config"
	"github.com/weastur/replacer/internal/generator"
)

const FailCode = 1

var version = "v0.0.0-dev1"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(
			os.Stderr,
			"\nFill the config, put `//go:generate replacer` at the top of your source file.\n"+
				"Refer to `go help generate` for more information about go code generators\n\nDefaults:\n",
		)
		flag.PrintDefaults()
	}
}

func main() {
	versionFlag := flag.Bool("version", false, "Print the version and exit")
	configFlag := flag.String(
		"config",
		"",
		"Path to the configuration file.\n"+
			"If not provided (default), the generator will look for a config file (\".replacer.yml\") "+
			"in the current directory, then move up to each parent directory until it reaches the root (the directory "+
			"containing go.mod).\nIf no config file is found, the generator will do nothing and exit with 0 code",
	)

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	configPath, err := config.Lookup(*configFlag)
	if errors.Is(err, config.ErrNotFound) {
		os.Exit(0)
	} else if err != nil {
		fmt.Printf("Error looking up the config file: %s\n", err)
		os.Exit(FailCode)
	}

	fmt.Printf("Using config file: %s\n", configPath)

	fmt.Println("Running my generator...")
	generator.Run()
}
