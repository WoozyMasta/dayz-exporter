package main

import (
	_ "embed"
	"fmt"
	"internal/vars"
	"os"
	"path/filepath"
	"strings"
)

//go:embed example.config.yaml
var exampleConfig []byte

//go:embed example.env
var exampleEnv []byte

// arguments parser
func parseArgs() {
	if len(os.Args) < 2 || !strings.HasPrefix(os.Args[1], "-") {
		return
	}

	switch os.Args[1] {
	case "--help", "-h":
		printHelp()
	case "--version", "-v":
		printVersion()
	case "--get-yaml", "-y":
		printExampleConfig("YAML config", exampleConfig)
	case "--get-env", "-e":
		printExampleConfig(".env file", exampleEnv)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command. Use --help for a list of available commands.")
		os.Exit(0)
	}
}

// printer for embedded files content
func printExampleConfig(name string, content []byte) {
	fmt.Printf("# Example %s\n", name)
	fmt.Println(string(content))
	os.Exit(0)
}

// just print help message and exit
func printHelp() {
	fmt.Printf(`%[1]s v%s
Collects and publishes Prometheus metrics from Battleye RCON and Steam A2S Query for DayZ server.

Usage:
  %[1]s [option] [config.yaml]

Available options:
  -y, --get-yaml    Prints an example YAML configuration file.
  -e, --get-env     Prints an example .env file.
  -v, --version     Show version, commit, and build time.
  -h, --help        Prints this help message.

Configuration File Lookup:
  - The program first checks for a configuration file specified by the variable 'DAYZ_EXPORTER_CONFIG_PATH'.
  - If the config variable is not set, it will attempt to use a configuration file passed as a command-line argument.
  - If no argument is provided, it will attempt to read from the default file 'config.yaml' in the current directory.

Examples:
  Save an example YAML configuration to file:
    %[1]s --get-yaml > config.yaml

  Print an example environment variables:
    %[1]s --get-env

  Run the program with a custom configuration file:
    DAYZ_EXPORTER_RCON_PASSWORD=strong %[1]s config.yaml

  Run the program normally (without any options):
    %[1]s

YAML config options have higher priority, for mixed use with variables comment out the overridden options.
For more information on configuration parameters, refer to the example configuration files (YAML and .env).
`, filepath.Base(os.Args[0]), vars.Version)
	os.Exit(0)
}

// print version information message and exit
func printVersion() {
	fmt.Printf(`
file:     %s
version:  %s
commit:   %s
built:    %s
project:  %s
`, os.Args[0], vars.Version, vars.Commit, vars.BuildTime, vars.URL)
	os.Exit(0)
}
