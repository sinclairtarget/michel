package site_test

import (
	"testing"

	"github.com/sinclairtarget/michel/internal/site"
)

func TestRelURL(t *testing.T) {
	tests := []struct {
		name     string
		suffix   string
		baseURL  string
		expected string
	}{
		{
			name:     "basic",
			suffix:   "foo/bar",
			baseURL:  "https://bim.com",
			expected: "/foo/bar",
		},
		{
			name:     "leading_slash",
			suffix:   "/foo/bar",
			baseURL:  "https://bim.com",
			expected: "/foo/bar",
		},
		{
			name:     "empty_baseURL",
			suffix:   "foo/bar",
			baseURL:  "",
			expected: "/foo/bar",
		},
		{
			name:     "baseURL_directory",
			suffix:   "foo/bar",
			baseURL:  "https://bim.com/bat",
			expected: "/bat/foo/bar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := site.RelURL(test.suffix, test.baseURL)
			if result != test.expected {
				t.Errorf(
					"RelURL was wrong; wanted \"%s\" but got \"%s\"",
					test.expected,
					result,
				)
			}
		})
	}
}

func TestAbsURL(t *testing.T) {
	tests := []struct {
		name     string
		suffix   string
		baseURL  string
		expected string
	}{
		{
			name:     "basic",
			suffix:   "foo/bar",
			baseURL:  "https://bim.com",
			expected: "https://bim.com/foo/bar",
		},
		{
			name:     "leading_slash",
			suffix:   "/foo/bar",
			baseURL:  "https://bim.com",
			expected: "https://bim.com/foo/bar",
		},
		{
			name:     "baseURL_directory",
			suffix:   "foo/bar",
			baseURL:  "https://bim.com/bat",
			expected: "https://bim.com/bat/foo/bar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := site.AbsURL(test.suffix, test.baseURL)
			if result != test.expected {
				t.Errorf(
					"AbsURL was wrong; wanted \"%s\" but got \"%s\"",
					test.expected,
					result,
				)
			}
		})
	}
}
