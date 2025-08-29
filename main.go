package main

import (
	"fmt"
	"github.com/markusdosch/gopherlol/internal/config"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var commandRegistry *config.CommandRegistry

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	// Parse command and arguments
	parts := strings.SplitN(q, " ", 3)
	cmdName := strings.ToLower(parts[0])
	
	// Handle help/list commands
	if cmdName == "list" || cmdName == "help" {
		generateHelpPage(w)
		return
	}

	// Try to find the command
	cmd := commandRegistry.FindCommand(cmdName)
	if cmd == nil {
		// Command not found => fall back to google
		fallbackURL := fmt.Sprintf("https://www.google.com/#q=%s", url.QueryEscape(q))
		http.Redirect(w, r, fallbackURL, http.StatusSeeOther)
		return
	}

	// Check for subcommands
	var targetURL string
	var query string

	if len(parts) >= 2 {
		// Check if second part is a subcommand
		subCmd := commandRegistry.FindSubcommand(cmdName, parts[1])
		if subCmd != nil {
			// Found subcommand, use remaining parts as query
			if len(parts) >= 3 {
				query = url.QueryEscape(parts[2])
			}
			var err error
			targetURL, err = config.ExecuteURL(subCmd.URL, query)
			if err != nil {
				log.Printf("Error executing subcommand URL template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			// No subcommand found, treat everything after command as query
			query = url.QueryEscape(strings.Join(parts[1:], " "))
			var err error
			targetURL, err = config.ExecuteURL(cmd.URL, query)
			if err != nil {
				log.Printf("Error executing command URL template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	} else {
		// No arguments, just the command
		if cmd.RequiresQuery {
			// Command requires query but none provided, fallback to google
			fallbackURL := fmt.Sprintf("https://www.google.com/#q=%s", url.QueryEscape(q))
			http.Redirect(w, r, fallbackURL, http.StatusSeeOther)
			return
		}
		var err error
		targetURL, err = config.ExecuteURL(cmd.URL, "")
		if err != nil {
			log.Printf("Error executing command URL template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, targetURL, http.StatusSeeOther)
}

func generateHelpPage(w http.ResponseWriter) {
	commands := commandRegistry.ListCommands()

	var html strings.Builder
	html.WriteString("<h1>gopherlol command list</h1>")
	html.WriteString("<ul>")
	
	for _, cmd := range commands {
		aliases := ""
		if len(cmd.Aliases) > 0 {
			aliases = fmt.Sprintf(" (aliases: %s)", strings.Join(cmd.Aliases, ", "))
		}

		requiresQuery := ""
		if cmd.RequiresQuery {
			requiresQuery = ", requires query"
		}

		html.WriteString(fmt.Sprintf(
			"<li><strong>%s</strong>%s%s - %s</li>",
			cmd.Name,
			aliases,
			requiresQuery,
			cmd.Description,
		))

		// Show subcommands if any
		if len(cmd.Subcommands) > 0 {
			html.WriteString("<ul>")
			for _, sub := range cmd.Subcommands {
				subAliases := ""
				if len(sub.Aliases) > 0 {
					subAliases = fmt.Sprintf(" (aliases: %s)", strings.Join(sub.Aliases, ", "))
				}
				html.WriteString(fmt.Sprintf(
					"<li><strong>%s %s</strong>%s - %s</li>",
					cmd.Name,
					sub.Name,
					subAliases,
					sub.Description,
				))
			}
			html.WriteString("</ul>")
		}
	}
	html.WriteString("</ul>")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, html.String())
}

func main() {
	// Load configuration
	configFile := "commands.json"
	commandConfig, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load command configuration: %v", err)
	}

	// Initialize command registry
	commandRegistry = config.NewCommandRegistry(commandConfig)

	log.Printf("Loaded %d commands from %s", len(commandConfig.Commands), configFile)
	
	http.HandleFunc("/", handler)
	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
