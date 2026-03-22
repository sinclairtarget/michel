package site

import (
	"log/slog"
	"path/filepath"
)

// An asset is any file we want to be part of the built site that is not an
// HTML file.
//
// Examples of assets: images, stylesheets, javascript.
//
// Unlike pages and content, the key for an asset includes the file extension.
// This makes it possible to have assets named e.g. main.js and main.css.
type AssetMetadata struct {
	key      string
	Filepath string
	relURL   string
	absURL   string
}

func (m AssetMetadata) Key() string { return m.key }

func (m AssetMetadata) RelURL() string { return m.relURL }

func (m AssetMetadata) AbsURL() string {
	if m.absURL == "" {
		slog.Warn(
			"no AbsURL for asset; did you configure baseURL?",
			"key",
			m.Key(),
		)
		return m.relURL
	}

	return m.absURL
}

func NewAsset(dir string, path string, baseURL string) AssetMetadata {
	key, err := filepath.Rel(dir, path)
	if err != nil {
		panic("asset path could not be made relative to site directory")
	}

	m := AssetMetadata{
		key:      key,
		Filepath: path,
		relURL:   RelURL(key, baseURL),
	}

	if baseURL != "" {
		m.absURL = AbsURL(key, baseURL)
	}

	return m
}
