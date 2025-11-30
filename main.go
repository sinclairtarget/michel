package main

import (
	"flag"
	"fmt"
	"os"
)

var Version string  // Semantic version
var BuildTag string // Optional tag

func main() {
	mainFlagSet := flag.NewFlagSet("michel", flag.ExitOnError)

	versionFlag := mainFlagSet.Bool("version", false, "Print version and exit")

	mainFlagSet.Usage = func() {
		fmt.Println("Usage: michel [subcommand] [subcommand options...]")
		fmt.Println("       michel --version")
		fmt.Println("michel builds websites from MyST markdown")

		fmt.Println()
		fmt.Println("Top-level options:")
		mainFlagSet.PrintDefaults()

		fmt.Println()
		fmt.Println("Subcommands:")
		fmt.Println("  build")
		fmt.Printf("\tBuilds site from content\n")
		fmt.Println("  serve")
		fmt.Printf("\tRuns local server for site\n")
	}

	mainFlagSet.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Println(getVersionString())
		return
	}

	subcommand := mainFlagSet.Arg(0)
	if subcommand == "serve" {
		fmt.Println("Run serve!")
	} else if subcommand == "build" || subcommand == "" {
		fmt.Println("Run build!")
	} else {
		fmt.Fprintf(os.Stderr, "Unrecognized subcommand: \"%s\"\n", subcommand)
		mainFlagSet.Usage()
		os.Exit(1)
	}
}

func getVersionString() string {
	if Version == "" {
		return "unknown"
	}

	if BuildTag != "" {
		return fmt.Sprintf("%s %s", Version, BuildTag)
	}

	return Version;
}
