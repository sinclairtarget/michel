package content

import (
	"errors"
	"iter"
	"log/slog"
	"maps"
	"slices"

	"github.com/sinclairtarget/michel/internal/content/myst"
	"github.com/sinclairtarget/michel/internal/merrors"
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
	used    map[string]bool // content that has been fully loaded via Get()
}

// Loads all content metadata into memory.
func LoadCorpus(dir string) (Corpus, error) {
	corpus := Corpus{
		entries: map[string]Entry{},
		used:    map[string]bool{},
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
		return Content{}, &merrors.KeyNotFoundError{
			Key:  key,
			Type: "content",
		}
	}

	// Record that we used this content file
	c.used[key] = true

	content, err := LoadContent(entry.Metadata)
	if err != nil {
		return content, err
	}

	return content, nil
}

func (c Corpus) GetMaybe(key string) (*Content, error) {
	content, err := c.Get(key)
	if err != nil {
		var keyerr *merrors.KeyNotFoundError
		if errors.As(err, &keyerr) {
			return nil, nil
		}

		return nil, err
	}

	return &content, nil
}

// Returns iterator over all content.
//
// This will load and parse markdown content for each content file.
func (c Corpus) All() iter.Seq[Entry] {
	return maps.Values(c.entries)
}

func (c Corpus) ByDate() iter.Seq[Entry] {
	return c.sorted(func(a, b Entry) int {
		if a.Date.Before(b.Date) {
			return -1
		} else if a.Date.Equal(b.Date) {
			return 0
		} else {
			return 1
		}
	})
}

func (c Corpus) ByTitle() iter.Seq[Entry] {
	return c.sorted(func(a, b Entry) int {
		if a.Title < b.Title {
			return -1
		} else if a.Title == b.Title {
			return 0
		} else {
			return 1
		}
	})
}

func (c Corpus) sorted(sortFunc func(Entry, Entry) int) iter.Seq[Entry] {
	values := slices.Collect(maps.Values(c.entries))
	slices.SortFunc(values, sortFunc)
	return slices.Values(values)
}

// This is a function rather than a method so it can't be called by users
// within templates.
func ReportUnused(c Corpus) {
	for entry := range c.All() {
		_, ok := c.used[entry.Key()]
		if !ok {
			slog.Warn("unused content", "path", entry.Filepath)
		}
	}
}
