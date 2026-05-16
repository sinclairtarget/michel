package site

import (
	"net/url"
	"strings"
)

// Returns an origin-relative URL incorporating any leading path part present
// in the given base URL.
//
// Any trailing index.html is stripped.
//
// e.g.
// foo/bar  https://bim.com     -> /foo/bar
// /foo/bar https://bim.com     -> /foo/bar
// foo/bar  ""                  -> /foo/bar
// foo/bar  https://bim.com/bat -> /bat/foo/bar
func RelURL(suffix string, baseURL string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	elems := strings.Split(suffix, "/")
	for _, elem := range elems {
		u = u.JoinPath(url.PathEscape(elem))
	}

	p := stripIndex(u.Path)
	if strings.HasPrefix(p, "/") {
		return p
	} else {
		return "/" + p
	}
}

// Returns an absolute URL incorporating the base URL.
//
// If no base URL is configured, panics.
//
// Any trailing index.html is stripped.
//
// e.g.
// foo/bar  https://bim.com     -> https://bim.com/foo/bar
// /foo/bar https://bim.com     -> https://bim.com/foo/bar
// foo/bar  ""                  -> panic!
// foo/bar  https://bim.com/bat -> https://bim.com/bat/foo/bar
func AbsURL(suffix string, baseURL string) string {
	if baseURL == "" {
		panic("can't compute absolute URL without base URL")
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	elems := strings.Split(suffix, "/")
	for _, elem := range elems {
		u = u.JoinPath(url.PathEscape(elem))
	}

	if u.IsAbs() {
		return stripIndex(u.String())
	} else if strings.HasPrefix(u.Path, "/") {
		return stripIndex(u.Path)
	} else {
		return "/" + stripIndex(u.Path)
	}
}

// TODO: Should this be configurable?
func stripIndex(s string) string {
	if strings.HasSuffix(s, "/index.html") {
		return strings.TrimSuffix(s, "index.html")
	}

	return s
}
