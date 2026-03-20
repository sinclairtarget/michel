package content

import (
	"fmt"
	"iter"
	"maps"

	"github.com/sinclairtarget/michel/internal/util"
)

// Collection of content loaded from disk.
//
// Metadata is kept in memory for every content file. The parsed MyST ASTs are
// loaded lazily.
//
// TODO: Cache the parsed MyST nodes so we don't have to re-read files, if this
// reveals itself to be more performant.
type Corpus struct {
	metadata map[string]Metadata
}

// Loads all content metadata into memory.
func LoadCorpus(dir string) (Corpus, error) {
	metadata := map[string]Metadata{}

	seq, finish := util.WalkFiles(dir)
	for path := range seq {
		m, err := LoadMetadata(dir, path)
		if err != nil {
			return Corpus{metadata}, err
		}

		metadata[m.key] = m
	}

	err := finish()
	if err != nil {
		return Corpus{metadata}, err
	}

	return Corpus{metadata}, nil
}

func (c Corpus) Get(key string) (Content, error) {
	metadata, ok := c.metadata[key]
	if !ok {
		return Content{}, fmt.Errorf(
			"content with key \"%s\" not found",
			key,
		)
	}

	content, err := LoadContent(metadata)
	if err != nil {
		return content, err
	}

	return content, nil
}

// Returns iterator over all content.
//
// This will load and parse markdown content for each content file.
func (c Corpus) All() iter.Seq[Content] {
	return func(yield func(Content) bool) {
		for k, _ := range c.metadata {
			content, err := c.Get(k)
			if err != nil {
				// This method is meant to be called within templates. A panic
				// during template execution will return an error from
				// tmpl.Execute().
				panic(err)
			}

			if !yield(content) {
				return
			}
		}
	}
}

// Returns iterator over all content metadata.
func (c Corpus) Meta() iter.Seq[Metadata] {
	return maps.Values(c.metadata)
}
