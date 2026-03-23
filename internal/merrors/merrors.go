// Custom Michel errors.
package merrors

import (
	"fmt"
)

type KeyNotFoundError struct {
	Key  string
	Type string // e.g. content, page, asset
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("%s with key \"%s\" not found", e.Type, e.Key)
}

func (e KeyNotFoundError) Suggestion() string {
	return fmt.Sprintf(
		"Is there a file matching \"%s\"? Is the key correct?",
		e.Key,
	)
}
