package cli

import (
	"strings"
	"testing"

	"github.com/sinclairtarget/michel/test/michel"
)

func checkContains(
	t *testing.T,
	result michel.Result,
	stream string,
	needle string,
) {
	haystack := result.Stdout
	if stream == "stderr" {
		haystack = result.Stderr
	}

	if !strings.Contains(strings.ToLower(haystack), strings.ToLower(needle)) {
		t.Errorf("expected %s to contain \"%s\"", stream, needle)
	}
}

// You can run `michel -h` or `michel --help` to get usage information.
func TestHelp(t *testing.T) {
	for _, flag := range []string{ "-h", "--help" } {
		result, err := michel.Run(flag)
		if err != nil {
			t.Fatal(err)
		}

		checkContains(t, result, "stdout", "Usage")
	}
}

// You can run `michel -version` or `michel --version` to get the current
// version.
func TestVersion(t *testing.T) {
	for _, flag := range []string{ "--version", "-version" } {
		result, err := michel.Run(flag)
		if err != nil {
			t.Fatal(err)
		}

		version := strings.TrimSpace(result.Stdout)

		// The version consists of the latest Git tag (possibly including a
		// commit hash if the HEAD coommit at build time had no tag, and
		// possibly including the string `dirty` if there are uncommitted
		// changes in the working directory at build time) and an optional build
		// tag.
		//
		// Examples:
		// v0.1
		// v0.1 dev
		// v0.12.1 test
		// v0.1-1-abc123 test
		if version == "" {
			t.Fatal("version string was empty")
		}

		parts := strings.Split(version, " ")
		if len(parts) < 2 {
			t.Fatalf("version string \"%s\" has no build tag", version)
		}

		buildTag := parts[len(parts) - 1]
		if buildTag != "test" {
			t.Fatalf("build tag should be \"test\", but is \"%s\"", buildTag)
		}
	}
}

// Specifying a non-existent subcommand produces an error.
func TestNonExistentSubcommand(t *testing.T) {
	result, err := michel.Run("foo")
	if err != nil {
		t.Fatal(err)
	}

	if result.ExitCode == 0 {
		t.Errorf("exit code should be non-zero but was %d", result.ExitCode)
	}

	checkContains(t, result, "stderr", "Unrecognized subcommand")
}
