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
	}

	mainFlagSet.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Println(getVersionString())
		return
	}
}

func getVersionString() string {
	if (Version == "") {
		return "unknown"
	}

	if (BuildTag != "") {
		return fmt.Sprintf("%s-%s", Version, BuildTag)
	}

	return Version;
}
