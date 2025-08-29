package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `{
		"commands": [
			{
				"name": "test",
				"aliases": ["t"],
				"description": "Test command",
				"url": "https://example.com/{{.Query}}",
				"requiresQuery": true
			}
		]
	}`

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_commands.json")

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(config.Commands))
	}

	cmd := config.Commands[0]
	if cmd.Name != "test" {
		t.Errorf("Expected command name 'test', got %q", cmd.Name)
	}

	if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "t" {
		t.Errorf("Expected alias 't', got %v", cmd.Aliases)
	}

	if !cmd.RequiresQuery {
		t.Error("Expected RequiresQuery to be true")
	}
}

func TestNewCommandRegistry(t *testing.T) {
	config := &CommandConfig{
		Commands: []Command{
			{
				Name:          "google",
				Aliases:       []string{"g", "search"},
				Description:   "Search Google",
				URL:           "https://google.com/search?q={{.Query}}",
				RequiresQuery: true,
				Default:       true,
			},
			{
				Name:          "github",
				Aliases:       []string{"gh"},
				Description:   "GitHub",
				URL:           "https://github.com",
				RequiresQuery: false,
				Subcommands: []Subcommand{
					{
						Name:        "pr",
						Aliases:     []string{"pull"},
						Description: "Pull requests",
						URL:         "https://github.com/pulls?q={{.Query}}",
					},
				},
			},
		},
	}

	registry := NewCommandRegistry(config)

	// Test finding main commands
	if cmd := registry.FindCommand("google"); cmd == nil {
		t.Error("Failed to find 'google' command")
	}

	if cmd := registry.FindCommand("g"); cmd == nil || cmd.Name != "google" {
		t.Error("Failed to find 'google' command by alias 'g'")
	}

	if cmd := registry.FindCommand("search"); cmd == nil || cmd.Name != "google" {
		t.Error("Failed to find 'google' command by alias 'search'")
	}

	// Test finding subcommands
	if sub := registry.FindSubcommand("github", "pr"); sub == nil {
		t.Error("Failed to find 'pr' subcommand for 'github'")
	}

	if sub := registry.FindSubcommand("gh", "pull"); sub == nil || sub.Name != "pr" {
		t.Error("Failed to find 'pr' subcommand by alias 'pull' for alias 'gh'")
	}

	// Test not finding non-existent commands
	if cmd := registry.FindCommand("nonexistent"); cmd != nil {
		t.Error("Found non-existent command")
	}

	if sub := registry.FindSubcommand("github", "nonexistent"); sub != nil {
		t.Error("Found non-existent subcommand")
	}

	// Test default command
	defaultCmd := registry.GetDefaultCommand()
	if defaultCmd == nil {
		t.Error("Expected to find default command")
	} else if defaultCmd.Name != "google" {
		t.Errorf("Expected default command to be 'google', got %q", defaultCmd.Name)
	}
}

func TestExecuteURL(t *testing.T) {
	testCases := []struct {
		template string
		query    string
		expected string
	}{
		{
			template: "https://example.com/search?q={{.Query}}",
			query:    "test query",
			expected: "https://example.com/search?q=test query",
		},
		{
			template: "https://example.com/",
			query:    "",
			expected: "https://example.com/",
		},
		{
			template: "https://example.com/{{.Query}}/page",
			query:    "results",
			expected: "https://example.com/results/page",
		},
	}

	for _, tc := range testCases {
		result, err := ExecuteURL(tc.template, tc.query)
		if err != nil {
			t.Errorf("ExecuteURL failed for template %q with query %q: %v", tc.template, tc.query, err)
			continue
		}

		if result != tc.expected {
			t.Errorf("ExecuteURL(%q, %q) = %q, expected %q", tc.template, tc.query, result, tc.expected)
		}
	}
}

func TestExecuteURL_InvalidTemplate(t *testing.T) {
	_, err := ExecuteURL("{{.InvalidField}}", "test")
	if err == nil {
		t.Error("Expected error for invalid template, got nil")
	}
}

func TestListCommands(t *testing.T) {
	config := &CommandConfig{
		Commands: []Command{
			{Name: "cmd1", Description: "First command"},
			{Name: "cmd2", Description: "Second command"},
		},
	}

	registry := NewCommandRegistry(config)
	commands := registry.ListCommands()

	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	// Check that both commands are present
	found := make(map[string]bool)
	for _, cmd := range commands {
		found[cmd.Name] = true
	}

	if !found["cmd1"] || !found["cmd2"] {
		t.Error("Not all commands found in list")
	}
}
