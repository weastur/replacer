package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const ExitFailureCode = 1

var (
	version     = "v1.0.0"
	verboseFlag bool
)

func PrintVerbose(format string, a ...any) {
	if verboseFlag {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

func Print(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func init() {
	flag.Usage = func() {
		Print("Usage:\n")
		Print(
			"\nFill the config, put `//go:generate replacer` at the top of your source file.\n" +
				"Refer to `go help generate` for more information about go code generators\n\nDefaults:\n",
		)
		flag.PrintDefaults()
	}
}

func main() {
	versionFlag := flag.Bool("version", false, "Print the version and exit")
	flag.BoolVar(&verboseFlag, "verbose", false, "Enable verbose output")
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
		Print("%s\n", version)
		os.Exit(0)
	}

	configPath, err := LookupConfig(*configFlag)
	if errors.Is(err, ErrNotFound) {
		PrintVerbose("No config file found, exiting without error\n")
		os.Exit(0)
	} else if err != nil {
		Print("Error looking up the config file: %s\n", err)
		os.Exit(ExitFailureCode)
	}

	PrintVerbose("Using config file: %s\n", configPath)

	cfg, err := LoadConfig(configPath)
	if err != nil {
		Print("Error loading the config file: %s\n", err)
		os.Exit(ExitFailureCode)
	}

	if err := Run(cfg); err != nil {
		Print("Error running the generator: %s\n", err)
		os.Exit(ExitFailureCode)
	}
}
