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

	var suggestErr merrors.SuggestError
	if errors.As(err, &suggestErr) {
		fmt.Fprintf(os.Stderr, "Error during build: %v\n", suggestErr)
		fmt.Fprintf(os.Stderr, "  %v\n", err)
		fmt.Println()
		fmt.Println(suggestErr.Suggestion())
	} else {
		fmt.Fprintf(os.Stderr, "Error during build: %v\n", err)
	}
}
