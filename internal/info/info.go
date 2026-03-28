package info

import (
	"fmt"

	atrus "github.com/sinclairtarget/libatrus-go"
)

var Version string  // Semantic version
var BuildTag string // Optional tag

func GetVersionString() string {
	if Version == "" {
		return "unknown"
	}

	if BuildTag != "" {
		return fmt.Sprintf("%s %s", Version, BuildTag)
	}

	return Version
}

func GetAtrusVersionString() string {
	return fmt.Sprintf("libatrus: %s", atrus.Version())
}
