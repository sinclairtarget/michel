package cli

import (
	"io"
	"os"
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

func copyTestdataFile(name string, destDir string) error {
	fin, err := os.Open("testdata/" + name)
	if err != nil {
		return err
	}
	defer fin.Close()

	fout, err := os.OpenFile(
		destDir+"/"+name,
		os.O_WRONLY|os.O_CREATE,
		0644,
	)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	return err
}

// You can run `michel -h` or `michel --help` to get usage information.
func TestHelp(t *testing.T) {
	for _, flag := range []string{"-h", "--help"} {
		result, err := michel.Run(t.TempDir(), flag)
		if err != nil {
			t.Fatal(err)
		}

		checkContains(t, result, "stdout", "Usage")
	}
}

func TestConfig(t *testing.T) {
	dir := t.TempDir()
	err := copyTestdataFile("michel.yaml", dir)
	if err != nil {
		t.Fatalf("failed to copy testdata file: %v", err)
	}

	result, err := michel.Run(dir, "config")
	if err != nil {
		t.Fatal(err)
	}

	expected := `title: CLI Test
description: This is a test site.`
	actual := strings.TrimSpace(result.Stdout)
	if actual != expected {
		t.Errorf(
			"output did not match. Wanted...\n%s\nGot...\n%s",
			expected,
			actual,
		)
	}
}

// You can run `michel version` to get the current version.
func TestVersion(t *testing.T) {
	result, err := michel.Run(t.TempDir(), "version")
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	if len(lines) != 2 {
		t.Fatalf("output of version subcommand should have two lines")
	}

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
	version := lines[0]
	if version == "" {
		t.Fatal("version string was empty")
	}

	parts := strings.Split(version, " ")
	if len(parts) < 2 {
		t.Fatalf("version string \"%s\" has no build tag", version)
	}

	buildTag := parts[len(parts)-1]
	if buildTag != "test" {
		t.Errorf("build tag should be \"test\", but is \"%s\"", buildTag)
	}

	libatrusVersion := lines[1]
	if !strings.HasPrefix(libatrusVersion, "libatrus: ") {
		t.Errorf("libatrus version line should start with \"libatrus: \"")
	}
}

// Specifying a non-existent subcommand produces an error.
func TestNonExistentSubcommand(t *testing.T) {
	result, err := michel.Run(t.TempDir(), "foo")
	if err != nil {
		t.Fatal(err)
	}

	if result.ExitCode == 0 {
		t.Errorf("exit code should be non-zero but was %d", result.ExitCode)
	}

	checkContains(t, result, "stderr", "Unrecognized subcommand")
}
