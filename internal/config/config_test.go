package config_test

import (
	"testing"

	"github.com/sinclairtarget/michel/internal/config"
)

func TestDefault(t *testing.T) {
	defaults := config.DefaultConfig()
	c, err := config.Load()
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}

	// Just test a few fields for equality
	if c.Title != defaults.Title {
		t.Errorf(
			"title not equal; wanted \"%s\", got: \"%s\"",
			defaults.Title,
			c.Title,
		)
	}
	if c.Description != defaults.Description {
		t.Errorf(
			"description not equal; wanted \"%s\", got: \"%s\"",
			defaults.Description,
			c.Description,
		)
	}
}

func TestDump(t *testing.T) {
	defaults := config.DefaultConfig()
	s := defaults.Dump()
	if s == "" {
		t.Error("dumped config was empty")
	}
}
