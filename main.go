package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/sinclairtarget/michel/internal/build"
	"github.com/sinclairtarget/michel/internal/server"
	"github.com/sinclairtarget/michel/internal/watch"
)

var Version string  // Semantic version
var BuildTag string // Optional tag

func main() {
	mainFlagSet := flag.NewFlagSet("michel", flag.ExitOnError)

	versionFlag := mainFlagSet.Bool("version", false, "Print version and exit")
	verboseFlag := mainFlagSet.Bool("v", false, "Enable verbose logging")

	mainFlagSet.Usage = func() {
		fmt.Println("Usage: michel [-v] SUBCOMMAND [subcommand options...]")
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

	var logger *slog.Logger
	if *verboseFlag {
		logger = configureLogging(slog.LevelDebug)
	} else {
		logger = configureLogging(slog.LevelInfo)
	}

	subcommand := mainFlagSet.Arg(0)
	if subcommand == "serve" {
		runServer(logger)
	} else if subcommand == "build" || subcommand == "" {
		runBuild(logger)
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

	return Version
}

func configureLogging(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{
			Level: level,
		},
	)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	logger.Debug("logging configured", "configuredLevel", level)
	return logger
}

func runServer(logger *slog.Logger) {
	serverBuild := func() {
		start := time.Now()
		err := build.Build(logger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during build: %v\n", err)
		}

		elapsed := time.Now().Sub(start)
		fmt.Printf("Site rebuilt in %dms.\n", elapsed.Milliseconds())
	}
	serverBuild() // Make sure `public` is up-to-date

	watcher := watch.NewWatcher("site")
	defer watcher.Close()

	go func() {
		for event := range watcher.Events {
			logger.Debug("got file modified event", "path", event.Path)
			serverBuild()
		}

		logger.Debug("goroutine exiting; watch events channel closed")
	}()

	err := watcher.Start(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start file watcher: %v", err)
		os.Exit(1)
	}

	err = server.Run(logger, "./public", 8080)
	fmt.Fprintf(os.Stderr, "Server exited: %v", err)
}

func runBuild(logger *slog.Logger) {
	err := build.Build(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during build: %v\n", err)
		os.Exit(1)
	}
}
