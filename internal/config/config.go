package config

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

// Fixed filename for config file
const Filename string = "michel.yaml"

// Configuration for the Michel site.
type Config struct {
	Title       string
	Description string
	BaseURL     string
}

// Returns the default config.
func DefaultConfig() Config {
	return Config{
		Title: "My Michel Site",
	}
}

// Returns a YAML string representing the config.
func (c Config) Dump() string {
	d, err := yaml.Marshal(&c)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal config: %v", err))
	}
	return string(d)
}

// Loads the config from disk.
//
// We first instantiate the default config, then update it with any non-empty
// fields loaded from disk.
func Load() (Config, error) {
	c := DefaultConfig()

	f, err := os.Open(Filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return c, nil
		}

		return c, err
	}
	defer f.Close()

	d, err := io.ReadAll(f)
	if err != nil {
		return c, err
	}

	loaded := Config{}
	err = yaml.Unmarshal(d, &loaded)
	if err != nil {
		return c, err
	}

	if loaded.Title != "" {
		c.Title = loaded.Title
	}
	if loaded.Description != "" {
		c.Description = loaded.Description
	}
	if loaded.BaseURL != "" {
		c.BaseURL = loaded.BaseURL
	}

	return c, nil
}
