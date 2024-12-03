package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/woozymasta/dayz-exporter/pkg/config"
)

//go:embed example.config.yaml
var exampleConfig []byte

//go:embed example.env
var exampleEnv []byte

func parseArgs() {
	if len(os.Args) < 2 || !strings.HasPrefix(os.Args[1], "-") {
		return
	}

	switch os.Args[1] {
	case "--help", "-h":
		printHelp(filepath.Base(os.Args[0]))
	case "--version", "-v":
		printVersion(filepath.Base(os.Args[0]))
	case "--get-yaml", "-y":
		printExampleConfig("YAML config", exampleConfig)
	case "--get-env", "-e":
		printExampleConfig(".env file", exampleEnv)
	default:
		fmt.Println("Unknown command. Use --help for a list of available commands.")
		os.Exit(0)
	}
}

// Функция для печати содержимого встроенного файла
func printExampleConfig(name string, content []byte) {
	fmt.Printf("# Example %s\n", name)
	fmt.Println(string(content))
	os.Exit(0)
}

// Функция для вывода справки
func printHelp(binary string) {
	helpText := fmt.Sprintf(`%[1]s v%s
Collects and publishes Prometheus metrics from Battleye RCON and Steam A2S Query for DayZ server.

Usage:
  %[1]s [option] [config.yaml]

Available options:
  -y, --get-yaml    Prints an example YAML configuration file.
  -e, --get-env     Prints an example .env file.
  -v, --version     Show version, commit, and build time.
  -h, --help        Prints this help message.

Configuration File Lookup:
  - The program first checks for a configuration file specified by the variable 'BERCON_EXPORTER_CONFIG_PATH'.
  - If the config variable is not set, it will attempt to use a configuration file passed as a command-line argument.
  - If no argument is provided, it will attempt to read from the default file 'config.yaml' in the current directory.

Examples:
  Save an example YAML configuration to file:
    %[1]s --get-yaml > config.yaml

  Print an example environment variables:
    %[1]s --get-env

  Run the program with a custom configuration file:
    BERCON_EXPORTER_RCON_PASSWORD=strong %[1]s config.yaml

  Run the program normally (without any options):
    %[1]s

YAML config options have higher priority, for mixed use with variables comment out the overridden options.
For more information on configuration parameters, refer to the example configuration files (YAML and .env).
`, binary, config.Version)

	fmt.Println(helpText)
	os.Exit(0)
}

// Функция для вывода версии и сборки
func printVersion(binary string) {
	fmt.Printf("%s\nversion=%s\ncommit=%s\nbuilt=%s\n", binary, config.Version, config.Commit, config.BuildTime)
	os.Exit(0)
}
