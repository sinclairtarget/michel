package build

import (
	"errors"
	"fmt"
	"os"

	"github.com/sinclairtarget/michel/internal/merrors"
)

func PrintBuildError(err error) {
	if err == nil {
		return
	}

	var keyerr *merrors.KeyNotFoundError
	if errors.As(err, &keyerr) {
		fmt.Fprintf(os.Stderr, "Error during build: %v\n", keyerr)
		fmt.Fprintf(os.Stderr, "  %v\n", err)
		fmt.Println()
		fmt.Println(keyerr.Suggestion())
	} else {
		fmt.Fprintf(os.Stderr, "Error during build: %v\n", err)
	}
}
