package main

import (
	"github.com/dominikoh/gopherlol/internal/analytics"
	"github.com/dominikoh/gopherlol/internal/config"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setupTestRegistry() {
	// Initialize analytics for testing
	analyticsSystem = analytics.NewAnalytics("test_usage.log")
	
	testConfig := &config.CommandConfig{
		Commands: []config.Command{
			{
				Name:          "google",
				Aliases:       []string{"g", "search"},
				Description:   "Search Google",
				URL:           "https://www.google.com/?q={{.Query}}",
				RequiresQuery: true,
				Default:       true,
			},
			{
				Name:          "stackoverflow",
				Aliases:       []string{"so", "stack"},
				Description:   "Search Stack Overflow",
				URL:           "https://stackoverflow.com/search?q={{.Query}}",
				RequiresQuery: true,
			},
			{
				Name:          "author",
				Aliases:       []string{},
				Description:   "Go to author's website",
				URL:           "https://www.markusdosch.com",
				RequiresQuery: false,
			},
			{
				Name:          "github",
				Aliases:       []string{"gh"},
				Description:   "GitHub operations",
				URL:           "https://github.com/search?q={{.Query}}",
				RequiresQuery: false,
				Subcommands: []config.Subcommand{
					{
						Name:        "pr",
						Aliases:     []string{"pull"},
						Description: "GitHub pull requests",
						URL:         "https://github.com/search?type=pullrequests&q={{.Query}}",
					},
				},
			},
		},
	}
	commandRegistry = config.NewCommandRegistry(testConfig)
}

func TestHandler_Help(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/?q=help", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<h1>gopherlol command list</h1>") {
		t.Error("Expected help page to contain title")
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected Content-Type to contain text/html, got %q", contentType)
	}
}

func TestHandler_List(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/?q=list", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<h1>gopherlol command list</h1>") {
		t.Error("Expected list page to contain title")
	}

	// Check that known commands appear in the list
	if !strings.Contains(body, "author") {
		t.Error("Expected list to contain 'author' command")
	}
	if !strings.Contains(body, "g") {
		t.Error("Expected list to contain 'g' command")
	}
	if !strings.Contains(body, "so") {
		t.Error("Expected list to contain 'so' command")
	}
}

func TestHandler_GoogleCommand(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/?q=g%20test%20query", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	expectedURL := "https://www.google.com/?q=test+query"
	if location != expectedURL {
		t.Errorf("Expected Location header %q, got %q", expectedURL, location)
	}
}

func TestHandler_StackOverflowCommand(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/?q=so%20golang%20testing", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	expectedURL := "https://stackoverflow.com/search?q=golang+testing"
	if location != expectedURL {
		t.Errorf("Expected Location header %q, got %q", expectedURL, location)
	}
}

func TestHandler_AuthorCommand(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/?q=author", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	expectedURL := "https://www.markusdosch.com"
	if location != expectedURL {
		t.Errorf("Expected Location header %q, got %q", expectedURL, location)
	}
}

func TestHandler_FallbackToGoogle(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/?q=nonexistent%20command%20test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	expectedURL := "https://www.google.com/?q=nonexistent+command+test"
	if location != expectedURL {
		t.Errorf("Expected Location header %q, got %q", expectedURL, location)
	}
}

func TestHandler_EmptyQuery(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/?q=", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	expectedURL := "https://www.google.com/?q="
	if location != expectedURL {
		t.Errorf("Expected Location header %q, got %q", expectedURL, location)
	}
}

func TestHandler_NoQuery(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	expectedURL := "https://www.google.com/?q="
	if location != expectedURL {
		t.Errorf("Expected Location header %q, got %q", expectedURL, location)
	}
}

func TestHandler_CaseInsensitive(t *testing.T) {
	setupTestRegistry()

	testCases := []struct {
		query    string
		expected string
	}{
		{"g%20test", "https://www.google.com/?q=test"},
		{"G%20test", "https://www.google.com/?q=test"},
		{"author", "https://www.markusdosch.com"},
		{"Author", "https://www.markusdosch.com"},
		{"so%20test", "https://stackoverflow.com/search?q=test"},
		{"So%20test", "https://stackoverflow.com/search?q=test"},
	}

	for _, tc := range testCases {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/?q="+tc.query, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusSeeOther {
				t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
			}

			location := w.Header().Get("Location")
			if location != tc.expected {
				t.Errorf("Expected Location header %q, got %q", tc.expected, location)
			}
		})
	}
}

func TestHandler_URLEncoding(t *testing.T) {
	setupTestRegistry()

	req := httptest.NewRequest("GET", "/?q=g%20hello%20world%20%26%20special%20chars%21", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	if !strings.Contains(location, "hello+world+%26+special+chars%21") {
		t.Errorf("Expected URL to be properly encoded, got %q", location)
	}
}

func TestHandler_FallbackToDefault(t *testing.T) {
	setupTestRegistry()
	defer os.Remove("test_usage.log") // Clean up test log file

	req := httptest.NewRequest("GET", "/?q=nonexistent%20query", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	expectedURL := "https://www.google.com/?q=nonexistent+query"
	if location != expectedURL {
		t.Errorf("Expected Location header %q, got %q", expectedURL, location)
	}
}
