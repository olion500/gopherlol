package commands

import (
	"net/url"
	"strings"
	"testing"
)

func TestCommands_G(t *testing.T) {
	c := &Commands{}
	result := c.G("hello world")

	expectedBase := "https://www.google.com/#q="
	if !strings.HasPrefix(result, expectedBase) {
		t.Errorf("Expected result to start with %q, got %q", expectedBase, result)
	}

	// Check URL encoding (Go's QueryEscape uses + for spaces)
	if !strings.Contains(result, "hello+world") {
		t.Errorf("Expected result to contain URL encoded query, got %q", result)
	}
}

func TestCommands_G_EmptyArg(t *testing.T) {
	c := &Commands{}
	result := c.G("")

	expected := "https://www.google.com/#q="
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCommands_G_SpecialCharacters(t *testing.T) {
	c := &Commands{}
	result := c.G("test & special chars!")

	expectedEncoded := url.QueryEscape("test & special chars!")
	expectedURL := "https://www.google.com/#q=" + expectedEncoded

	if result != expectedURL {
		t.Errorf("Expected %q, got %q", expectedURL, result)
	}
}

func TestCommands_Author(t *testing.T) {
	c := &Commands{}
	result := c.Author()

	expected := "https://www.markusdosch.com"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCommands_So(t *testing.T) {
	c := &Commands{}
	result := c.So("golang testing")

	expectedBase := "https://stackoverflow.com/search?q="
	if !strings.HasPrefix(result, expectedBase) {
		t.Errorf("Expected result to start with %q, got %q", expectedBase, result)
	}

	// Check URL encoding (Go's QueryEscape uses + for spaces)
	if !strings.Contains(result, "golang+testing") {
		t.Errorf("Expected result to contain URL encoded query, got %q", result)
	}
}

func TestCommands_So_EmptyArg(t *testing.T) {
	c := &Commands{}
	result := c.So("")

	expected := "https://stackoverflow.com/search?q="
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCommands_So_SpecialCharacters(t *testing.T) {
	c := &Commands{}
	result := c.So("how to handle & in golang?")

	expectedEncoded := url.QueryEscape("how to handle & in golang?")
	expectedURL := "https://stackoverflow.com/search?q=" + expectedEncoded

	if result != expectedURL {
		t.Errorf("Expected %q, got %q", expectedURL, result)
	}
}