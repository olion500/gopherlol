package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// CommandConfig represents the JSON configuration structure
type CommandConfig struct {
	Commands []Command `json:"commands"`
}

// Command represents a single command configuration
type Command struct {
	Name          string       `json:"name"`
	Aliases       []string     `json:"aliases"`
	Description   string       `json:"description"`
	URL           string       `json:"url"`
	RequiresQuery bool         `json:"requiresQuery"`
	Default       bool         `json:"default,omitempty"`
	Subcommands   []Subcommand `json:"subcommands,omitempty"`
}

// Subcommand represents a subcommand configuration
type Subcommand struct {
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
}

// TemplateData holds data for URL template processing
type TemplateData struct {
	Query string
}

// CommandRegistry manages command lookup and execution
type CommandRegistry struct {
	commands       map[string]*Command
	aliases        map[string]*Command
	subcommands    map[string]map[string]*Subcommand
	defaultCommand *Command
}

// LoadConfig loads command configuration from a JSON file
func LoadConfig(filename string) (*CommandConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config CommandConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &config, nil
}

// NewCommandRegistry creates a new command registry from configuration
func NewCommandRegistry(config *CommandConfig) *CommandRegistry {
	registry := &CommandRegistry{
		commands:       make(map[string]*Command),
		aliases:        make(map[string]*Command),
		subcommands:    make(map[string]map[string]*Subcommand),
		defaultCommand: nil,
	}

	// Register commands and aliases
	for i := range config.Commands {
		cmd := &config.Commands[i]

		// Register main command name
		registry.commands[strings.ToLower(cmd.Name)] = cmd

		// Register aliases
		for _, alias := range cmd.Aliases {
			registry.aliases[strings.ToLower(alias)] = cmd
		}

		// Set default command if specified
		if cmd.Default {
			registry.defaultCommand = cmd
		}

		// Register subcommands if any
		if len(cmd.Subcommands) > 0 {
			subMap := make(map[string]*Subcommand)
			for j := range cmd.Subcommands {
				sub := &cmd.Subcommands[j]
				subMap[strings.ToLower(sub.Name)] = sub

				// Register subcommand aliases
				for _, alias := range sub.Aliases {
					subMap[strings.ToLower(alias)] = sub
				}
			}
			registry.subcommands[strings.ToLower(cmd.Name)] = subMap

			// Also register subcommands under aliases
			for _, alias := range cmd.Aliases {
				registry.subcommands[strings.ToLower(alias)] = subMap
			}
		}
	}

	return registry
}

// FindCommand looks up a command by name or alias
func (r *CommandRegistry) FindCommand(name string) *Command {
	name = strings.ToLower(name)

	if cmd, exists := r.commands[name]; exists {
		return cmd
	}

	if cmd, exists := r.aliases[name]; exists {
		return cmd
	}

	return nil
}

// FindSubcommand looks up a subcommand for a given command
func (r *CommandRegistry) FindSubcommand(cmdName, subName string) *Subcommand {
	cmdName = strings.ToLower(cmdName)
	subName = strings.ToLower(subName)

	if subMap, exists := r.subcommands[cmdName]; exists {
		if sub, exists := subMap[subName]; exists {
			return sub
		}
	}

	return nil
}

// GetDefaultCommand returns the configured default command
func (r *CommandRegistry) GetDefaultCommand() *Command {
	return r.defaultCommand
}

// ExecuteURL processes the URL template with the given query
func ExecuteURL(urlTemplate, query string) (string, error) {
	tmpl, err := template.New("url").Parse(urlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL template: %w", err)
	}

	var buf bytes.Buffer
	data := TemplateData{Query: query}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute URL template: %w", err)
	}

	return buf.String(), nil
}

// ListCommands returns all available commands for help display
func (r *CommandRegistry) ListCommands() []Command {
	var commands []Command
	seen := make(map[string]bool)

	for _, cmd := range r.commands {
		if !seen[cmd.Name] {
			commands = append(commands, *cmd)
			seen[cmd.Name] = true
		}
	}

	return commands
}
