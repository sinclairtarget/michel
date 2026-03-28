package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/sinclairtarget/michel/internal/build"
	"github.com/sinclairtarget/michel/internal/config"
	"github.com/sinclairtarget/michel/internal/info"
	"github.com/sinclairtarget/michel/internal/server"
)

type command struct {
	flagSet     *flag.FlagSet
	description string
	run         func(args []string)
}

func main() {
	subcommands := map[string]command{
		"build":   buildCmd(),
		"serve":   serveCmd(),
		"config":  configCmd(),
		"version": versionCmd(),
	}

	// handle top-level flags
	mainFlagSet := flag.NewFlagSet("michel", flag.ExitOnError)

	verboseFlag := mainFlagSet.Bool("verbose", false, "Enable verbose logging")

	mainFlagSet.Usage = func() {
		fmt.Println(
			"Usage: michel [-verbose] [SUBCOMMAND] [OPTIONS...]",
		)
		fmt.Println("michel builds websites from MyST markdown")

		fmt.Println()
		fmt.Println("Top-level options:")
		mainFlagSet.PrintDefaults()

		fmt.Println()
		fmt.Println("Subcommands:")
		for _, name := range []string{"build", "serve", "config", "version"} {
			cmd := subcommands[name]

			if name == "build" {
				fmt.Printf("  %s (default)\n", name)
			} else {
				fmt.Printf("  %s\n", name)
			}

			fmt.Printf("\t%s\n", cmd.description)
		}
	}

	// Look for the index of the first arg not intended as a top-level flag.
	// We handle this manually so that specifying the default subcommand is
	// optional even when providing subcommand flags.
	subcmdIndex := 1
loop:
	for subcmdIndex < len(os.Args) {
		switch os.Args[subcmdIndex] {
		case "-verbose", "--verbose", "-h", "--help":
			subcmdIndex += 1
		default:
			break loop
		}
	}

	mainFlagSet.Parse(os.Args[1:subcmdIndex])

	if *verboseFlag {
		configureLogging(slog.LevelDebug)
	} else {
		configureLogging(slog.LevelInfo)
	}

	args := os.Args[subcmdIndex:]

	// handle subcommands
	cmd := subcommands["build"] // Default to "build"
	if len(args) > 0 {
		first := args[0]
		subcommand, ok := subcommands[first]
		if ok {
			cmd = subcommand
			args = args[1:]
		} else {
			fmt.Fprintf(
				os.Stderr,
				"Unrecognized subcommand: \"%s\"\n",
				subcommand,
			)
			mainFlagSet.Usage()
			os.Exit(1)
		}
	}

	cmd.flagSet.Parse(args)
	subargs := cmd.flagSet.Args()
	cmd.run(subargs)
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

func buildCmd() command {
	flagSet := flag.NewFlagSet("michel build", flag.ExitOnError)

	outdir := flagSet.String(
		"o",
		build.DefaultOutputDir,
		"Output directory for build",
	)

	description := "Build site"

	flagSet.Usage = func() {
		fmt.Println("Usage: michel build [OPTIONS...]")
		fmt.Println(description)
		fmt.Println()
		flagSet.PrintDefaults()
	}

	return command{
		flagSet:     flagSet,
		description: description,
		run: func(args []string) {
			err := build.Build(*outdir)
			if err != nil {
				build.PrintBuildError(err)
				os.Exit(1)
			}
		},
	}
}

func serveCmd() command {
	flagSet := flag.NewFlagSet("michel serve", flag.ExitOnError)

	outdir := flagSet.String(
		"o",
		build.DefaultOutputDir,
		"Output directory for build",
	)

	description := "Run local HTTP server for site"

	flagSet.Usage = func() {
		fmt.Println("Usage: michel serve [OPTIONS...]")
		fmt.Println(description)
		fmt.Println()
		flagSet.PrintDefaults()
	}

	return command{
		flagSet:     flagSet,
		description: description,
		run: func(args []string) {
			// Build before running server
			err := build.Build(*outdir)
			if err != nil {
				build.PrintBuildError(err)
				os.Exit(1)
			}

			// Run server
			err = server.Run("./public", 8080, *outdir)
			fmt.Fprintf(os.Stderr, "Server exited: %v\n", err)
		},
	}
}

func configCmd() command {
	flagSet := flag.NewFlagSet("michel config", flag.ExitOnError)

	description := "Print parsed config"

	flagSet.Usage = func() {
		fmt.Println("Usage: michel config")
		fmt.Println(description)
	}

	return command{
		flagSet:     flagSet,
		description: description,
		run: func(args []string) {
			c, err := config.Load()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}

			s := c.Dump()
			fmt.Print(s)
		},
	}
}

func versionCmd() command {
	flagSet := flag.NewFlagSet("michel version", flag.ExitOnError)

	description := "Print version"

	flagSet.Usage = func() {
		fmt.Println("Usage: michel version")
		fmt.Println(description)
	}

	return command{
		flagSet:     flagSet,
		description: description,
		run: func(args []string) {
			fmt.Println(info.GetVersionString())
			fmt.Println(info.GetAtrusVersionString())
		},
	}
}
