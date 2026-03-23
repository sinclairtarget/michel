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
	return fmt.Sprintf("%s with key \"%s\" not found hi", e.Type, e.Key)
}
