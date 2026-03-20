package content

import (
	"fmt"

	"github.com/sinclairtarget/michel/internal/util"
)

// Collection of content loaded from disk.
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

	content, err := metadata.LoadContent()
	if err != nil {
		return content, err
	}

	return content, nil
}
