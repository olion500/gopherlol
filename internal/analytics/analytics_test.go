package analytics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewAnalytics(t *testing.T) {
	logFile := "test_analytics.log"
	analytics := NewAnalytics(logFile)

	if analytics.logFile != logFile {
		t.Errorf("Expected logFile to be %q, got %q", logFile, analytics.logFile)
	}

	if analytics.lastUsage == nil {
		t.Error("Expected lastUsage to be initialized")
	}

	if analytics.sessionStart.IsZero() {
		t.Error("Expected sessionStart to be set")
	}
}

func TestLogCommandUsage(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	analytics := NewAnalytics(logFile)

	// Log a command usage
	analytics.LogCommandUsage("google", "test query", "test-agent", "127.0.0.1", false, false, "")

	// Read the log file
	file, err := os.Open(logFile)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	var usage CommandUsage
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&usage); err != nil {
		t.Fatalf("Failed to decode log entry: %v", err)
	}

	if usage.Command != "google" {
		t.Errorf("Expected command 'google', got %q", usage.Command)
	}

	if usage.Query != "test query" {
		t.Errorf("Expected query 'test query', got %q", usage.Query)
	}

	if usage.UserAgent != "test-agent" {
		t.Errorf("Expected user agent 'test-agent', got %q", usage.UserAgent)
	}

	if usage.RemoteAddr != "127.0.0.1" {
		t.Errorf("Expected remote addr '127.0.0.1', got %q", usage.RemoteAddr)
	}

	if usage.IsDefault != false {
		t.Errorf("Expected IsDefault false, got %v", usage.IsDefault)
	}

	if usage.IsSubcommand != false {
		t.Errorf("Expected IsSubcommand false, got %v", usage.IsSubcommand)
	}
}

func TestLogCommandUsageWithSubcommand(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	analytics := NewAnalytics(logFile)

	// Log a subcommand usage
	analytics.LogCommandUsage("github", "test query", "test-agent", "127.0.0.1", false, true, "pr")

	// Read the log file
	file, err := os.Open(logFile)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	var usage CommandUsage
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&usage); err != nil {
		t.Fatalf("Failed to decode log entry: %v", err)
	}

	if usage.Command != "github" {
		t.Errorf("Expected command 'github', got %q", usage.Command)
	}

	if usage.IsSubcommand != true {
		t.Errorf("Expected IsSubcommand true, got %v", usage.IsSubcommand)
	}

	if usage.Subcommand != "pr" {
		t.Errorf("Expected subcommand 'pr', got %q", usage.Subcommand)
	}
}

func TestGetDayStats(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	analytics := NewAnalytics(logFile)

	// Create some test data
	today := time.Now().Format("2006-01-02")
	analytics.LogCommandUsage("google", "test1", "agent1", "127.0.0.1", false, false, "")
	analytics.LogCommandUsage("stackoverflow", "test2", "agent2", "127.0.0.2", false, false, "")
	analytics.LogCommandUsage("google", "test3", "agent1", "127.0.0.1", false, false, "")

	// Wait a bit to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	stats, err := analytics.GetDayStats(today)
	if err != nil {
		t.Fatalf("Failed to get day stats: %v", err)
	}

	if stats.Date != today {
		t.Errorf("Expected date %q, got %q", today, stats.Date)
	}

	if stats.TotalUsage != 3 {
		t.Errorf("Expected total usage 3, got %d", stats.TotalUsage)
	}

	if stats.Commands["google"] != 2 {
		t.Errorf("Expected google command count 2, got %d", stats.Commands["google"])
	}

	if stats.Commands["stackoverflow"] != 1 {
		t.Errorf("Expected stackoverflow command count 1, got %d", stats.Commands["stackoverflow"])
	}

	if stats.UniqueUsers != 2 {
		t.Errorf("Expected unique users 2, got %d", stats.UniqueUsers)
	}
}

func TestGetDayStats_NonExistentDate(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	analytics := NewAnalytics(logFile)

	stats, err := analytics.GetDayStats("2020-01-01")
	if err != nil {
		t.Fatalf("Failed to get day stats: %v", err)
	}

	if stats.TotalUsage != 0 {
		t.Errorf("Expected total usage 0, got %d", stats.TotalUsage)
	}

	if stats.UniqueUsers != 0 {
		t.Errorf("Expected unique users 0, got %d", stats.UniqueUsers)
	}
}

func TestGetOverallStats(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	analytics := NewAnalytics(logFile)

	// Create some test data across multiple days
	analytics.LogCommandUsage("google", "test1", "agent1", "127.0.0.1", false, false, "")
	analytics.LogCommandUsage("stackoverflow", "test2", "agent2", "127.0.0.2", false, false, "")
	analytics.LogCommandUsage("google", "test3", "agent3", "127.0.0.3", false, false, "")

	stats, err := analytics.GetOverallStats()
	if err != nil {
		t.Fatalf("Failed to get overall stats: %v", err)
	}

	if stats.Date != "all-time" {
		t.Errorf("Expected date 'all-time', got %q", stats.Date)
	}

	if stats.TotalUsage != 3 {
		t.Errorf("Expected total usage 3, got %d", stats.TotalUsage)
	}

	if stats.Commands["google"] != 2 {
		t.Errorf("Expected google command count 2, got %d", stats.Commands["google"])
	}

	if stats.Commands["stackoverflow"] != 1 {
		t.Errorf("Expected stackoverflow command count 1, got %d", stats.Commands["stackoverflow"])
	}

	if stats.UniqueUsers != 3 {
		t.Errorf("Expected unique users 3, got %d", stats.UniqueUsers)
	}

	// Check top commands
	if len(stats.TopCommands) == 0 {
		t.Error("Expected top commands to be populated")
	} else {
		if stats.TopCommands[0].Command != "google" {
			t.Errorf("Expected top command to be 'google', got %q", stats.TopCommands[0].Command)
		}
		if stats.TopCommands[0].Count != 2 {
			t.Errorf("Expected top command count 2, got %d", stats.TopCommands[0].Count)
		}
	}
}

func TestGetTopCommands(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	analytics := NewAnalytics(logFile)

	commandCounts := map[string]int{
		"google":        5,
		"stackoverflow": 3,
		"github":        7,
		"author":        1,
	}

	topCommands := analytics.getTopCommands(commandCounts, 2)

	if len(topCommands) != 2 {
		t.Errorf("Expected 2 top commands, got %d", len(topCommands))
	}

	// Should be sorted by count descending
	if topCommands[0].Command != "github" || topCommands[0].Count != 7 {
		t.Errorf("Expected first command to be 'github' with count 7, got %q with count %d",
			topCommands[0].Command, topCommands[0].Count)
	}

	if topCommands[1].Command != "google" || topCommands[1].Count != 5 {
		t.Errorf("Expected second command to be 'google' with count 5, got %q with count %d",
			topCommands[1].Command, topCommands[1].Count)
	}
}

func TestReadLogEntries_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "empty.log")
	analytics := NewAnalytics(logFile)

	// Create empty file
	file, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	file.Close()

	entries, err := analytics.readLogEntries()
	if err != nil {
		t.Fatalf("Failed to read empty log entries: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}

func TestReadLogEntries_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "nonexistent.log")
	analytics := NewAnalytics(logFile)

	entries, err := analytics.readLogEntries()
	if err != nil {
		t.Fatalf("Failed to read non-existent log entries: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for non-existent file, got %d", len(entries))
	}
}
