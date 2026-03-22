package util

import (
	"fmt"
	"iter"
	"regexp"
	"strings"
)

type Keyed interface {
	Key() string
}

const globPattern string = `.+`

// Consumes the given sequence, yielding each element with a key matching the
// given glob pattern.
//
// The glob is just a wildcard matching one or more arbitrary characters.
func Select[T Keyed](
	seq iter.Seq[T],
	field string,
	pattern string,
) iter.Seq[T] {
	field = strings.ToLower(field)
	if field != "key" {
		msg := fmt.Sprintf("unsupported field \"%s\" in select", field)
		panic(msg)
	}

	re := compileGlobRegex(pattern)
	return func(yield func(T) bool) {
		for elem := range seq {
			if re.MatchString(elem.Key()) {
				proceed := yield(elem)
				if !proceed {
					return
				}
			}
		}
	}
}

// Consumes the given sequence, yielding each element with a key NOT matching
// the given glob pattern.
//
// The glob is just a wildcard matching one or more arbitrary characters.
func Reject[T Keyed](
	seq iter.Seq[T],
	field string,
	pattern string,
) iter.Seq[T] {
	field = strings.ToLower(field)
	if field != "key" {
		msg := fmt.Sprintf("unsupported field \"%s\" in reject", field)
		panic(msg)
	}

	re := compileGlobRegex(pattern)
	return func(yield func(T) bool) {
		for elem := range seq {
			if !re.MatchString(elem.Key()) {
				proceed := yield(elem)
				if !proceed {
					return
				}
			}
		}
	}
}

func compileGlobRegex(pattern string) *regexp.Regexp {
	parts := strings.Split(pattern, "*")

	escaped := []string{}
	for _, p := range parts {
		// Handle adjacent "*"
		if p == "" {
			continue
		}

		escaped = append(escaped, regexp.QuoteMeta(p))
	}

	expr := strings.Join(escaped, globPattern)
	return regexp.MustCompile(expr)
}

// Converts iter.Seq[T] to iter.Seq[U].
//
// Performs a runtime type assertion.
func CoerceSeq[T any, U any](seq iter.Seq[T]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for elem := range seq {
			if !yield(any(elem).(U)) {
				return
			}
		}
	}
}
