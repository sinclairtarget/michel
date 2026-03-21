package util_test

import (
	"slices"
	"testing"

	"github.com/sinclairtarget/michel/internal/util"
)

type Dummy struct {
	key string
}

func (d Dummy) Key() string { return d.key }

func TestSelect(t *testing.T) {
	elements := []Dummy{
		Dummy{key: "blog/foo.md"},
		Dummy{key: "blog/bar.md"},
		Dummy{key: "blog/intro.md"},
		Dummy{key: "intro.md"},
		Dummy{key: "roll/foo.md"},
	}

	want := []string{
		"blog/foo.md",
		"blog/bar.md",
		"blog/intro.md",
	}

	got := []string{}
	for elem := range util.Select(slices.Values(elements), "blog/*") {
		got = append(got, elem.Key())
	}

	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestReject(t *testing.T) {
	elements := []Dummy{
		Dummy{key: "blog/foo.md"},
		Dummy{key: "blog/bar.md"},
		Dummy{key: "blog/intro.md"},
		Dummy{key: "blog.md"},
		Dummy{key: "roll/foo.md"},
	}

	want := []string{
		"blog/bar.md",
		"blog/intro.md",
		"blog.md",
	}

	got := []string{}
	for elem := range util.Reject(slices.Values(elements), "*/foo.md") {
		got = append(got, elem.Key())
	}

	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSelectReject(t *testing.T) {
	elements := []Dummy{
		Dummy{key: "blog/foo.md"},
		Dummy{key: "blog/bar.md"},
		Dummy{key: "blog/intro.md"},
		Dummy{key: "blog.md"},
		Dummy{key: "roll/foo.md"},
	}

	want := []string{
		"blog/foo.md",
		"blog/bar.md",
	}

	got := []string{}
	for elem := range util.Reject(
		util.Select(slices.Values(elements), "blog/*"),
		"blog/intro.md",
	) {
		got = append(got, elem.Key())
	}

	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
