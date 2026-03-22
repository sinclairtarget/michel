package site

import (
	"net/url"
	"strings"
)

// Returns an origin-relative URL incorporating any leading path part present
// in the given base URL.
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

	if strings.HasPrefix(u.Path, "/") {
		return u.Path
	} else {
		return "/" + u.Path
	}
}

// Returns an absolute URL incorporating the base URL.
//
// If no base URL is configured, returns an origin-relative URL.
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
		return u.String()
	} else if strings.HasPrefix(u.Path, "/") {
		return u.Path
	} else {
		return "/" + u.Path
	}
}
