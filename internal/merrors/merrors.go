// Custom Michel errors.
package merrors

import (
	"fmt"
)

type SuggestError interface {
	error
	Suggestion() string
}

// Raised when there is no contnet, page, or asset found with the given key.
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

// Raised when .Page.Content is used in a template but the page frontmatter
// does not include a content key.
type NoAssociatedContentError struct {
	PageKey      string
	PageFilepath string
}

func (e NoAssociatedContentError) Error() string {
	return fmt.Sprintf("page \"%s\" has no associated content", e.PageKey)
}

func (e NoAssociatedContentError) Suggestion() string {
	return fmt.Sprintf(
		"Add a \"content\" field to the frontmatter in \"%s\" to create the "+
			"association.",
		e.PageFilepath,
	)
}
