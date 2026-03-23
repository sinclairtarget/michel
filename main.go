package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	atrus "github.com/sinclairtarget/libatrus-go"

	"github.com/sinclairtarget/michel/internal/build"
	"github.com/sinclairtarget/michel/internal/config"
	"github.com/sinclairtarget/michel/internal/server"
)

var Version string  // Semantic version
var BuildTag string // Optional tag

func main() {
	mainFlagSet := flag.NewFlagSet("michel", flag.ExitOnError)

	verboseFlag := mainFlagSet.Bool("verbose", false, "Enable verbose logging")

	mainFlagSet.Usage = func() {
		fmt.Println(
			"Usage: michel [-verbose] [SUBCOMMAND] [subcommand options...]",
		)
		fmt.Println("michel builds websites from MyST markdown")

		fmt.Println()
		fmt.Println("Top-level options:")
		mainFlagSet.PrintDefaults()

		fmt.Println()
		fmt.Println("Subcommands:")
		fmt.Println("  build")
		fmt.Printf("\tBuilds site (default)\n")
		fmt.Println("  serve")
		fmt.Printf("\tRuns local HTTP server for site\n")
		fmt.Println("  config")
		fmt.Printf("\tPrints site config\n")
		fmt.Println("  version")
		fmt.Printf("\tPrint version and exit\n")
	}

	mainFlagSet.Parse(os.Args[1:])

	if *verboseFlag {
		configureLogging(slog.LevelDebug)
	} else {
		configureLogging(slog.LevelInfo)
	}

	subcommand := mainFlagSet.Arg(0)
	switch subcommand {
	case "build", "":
		runBuild()
	case "config":
		runConfig()
	case "serve":
		runServer()
	case "version":
		fmt.Println(getVersionString())
		fmt.Printf("libatrus: %s\n", atrus.Version())
	default:
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

	return Version
}

func configureLogging(level slog.Level) {
	handler := slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{
			Level: level,
		},
	)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Debug("logging configured", "configuredLevel", level)
}

func runBuild() {
	err := build.Build()
	if err != nil {
		build.PrintBuildError(err)
		os.Exit(1)
	}
}

func runConfig() {
	c, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	s := c.Dump()
	fmt.Print(s)
}

func runServer() {
	// Build before running server
	err := build.Build()
	if err != nil {
		build.PrintBuildError(err)
		os.Exit(1)
	}

	// Run server
	err = server.Run("./public", 8080)
	fmt.Fprintf(os.Stderr, "Server exited: %v\n", err)
}
