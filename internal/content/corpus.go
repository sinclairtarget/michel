package content

import (
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/sinclairtarget/michel/internal/content/myst"
	"github.com/sinclairtarget/michel/internal/util"
)

type Entry struct {
	Metadata
	corpus *Corpus
}

func (e Entry) Root() (*myst.Node, error) {
	content, err := e.corpus.Get(e.Key())
	if err != nil {
		return nil, err
	}

	return content.Root, nil
}

// Collection of content loaded from disk.
//
// Metadata is kept in memory for every content file. The parsed MyST ASTs are
// loaded lazily.
//
// TODO: Cache the parsed MyST nodes so we don't have to re-read files, if this
// reveals itself to be more performant.
type Corpus struct {
	entries map[string]Entry
}

// Loads all content metadata into memory.
func LoadCorpus(dir string) (Corpus, error) {
	corpus := Corpus{
		entries: map[string]Entry{},
	}

	seq, finish := util.WalkFiles(dir)
	for path := range seq {
		m, err := LoadMetadata(dir, path)
		if err != nil {
			return corpus, err
		}

		corpus.entries[m.key] = Entry{
			Metadata: m,
			corpus:   &corpus,
		}
	}

	err := finish()
	if err != nil {
		return corpus, err
	}

	return corpus, nil
}

func (c Corpus) Get(key string) (Content, error) {
	entry, ok := c.entries[key]
	if !ok {
		return Content{}, fmt.Errorf(
			"content with key \"%s\" not found",
			key,
		)
	}

	content, err := LoadContent(entry.Metadata)
	if err != nil {
		return content, err
	}

	return content, nil
}

// Returns iterator over all content.
//
// This will load and parse markdown content for each content file.
func (c Corpus) All() iter.Seq[Entry] {
	return maps.Values(c.entries)
}

func (c Corpus) ByDate() iter.Seq[Entry] {
	values := slices.Collect(maps.Values(c.entries))
	slices.SortFunc(values, func(a, b Entry) int {
		if a.Date.Before(b.Date) {
			return -1
		} else if a.Date.Equal(b.Date) {
			return 0
		} else {
			return 1
		}
	})

	return slices.Values(values)
}
