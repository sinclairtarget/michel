package util

import (
	"fmt"
	"iter"
	"maps"
)

type Collectable interface {
	Key() string
}

type Collection[T Collectable] struct {
	loaded map[string]T
}

func NewCollection[T Collectable]() Collection[T] {
	return Collection[T]{
		loaded: map[string]T{},
	}
}

func (c *Collection[T]) Add(key string, member T) {
	c.loaded[key] = member
}

func (c Collection[T]) Get(key string) (T, error) {
	member, ok := c.loaded[key]
	if !ok {
		return member, fmt.Errorf(
			"collection member with key \"%s\" not found",
			key,
		)
	}

	return member, nil
}

func (c Collection[T]) All() iter.Seq[T] {
	return maps.Values(c.loaded)
}
