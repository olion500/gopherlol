package analytics

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// CommandUsage represents a single command usage event
type CommandUsage struct {
	Command      string    `json:"command"`
	Query        string    `json:"query,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	UserAgent    string    `json:"user_agent,omitempty"`
	RemoteAddr   string    `json:"remote_addr,omitempty"`
	Duration     int64     `json:"duration_ms,omitempty"` // Time until next command
	IsDefault    bool      `json:"is_default"`            // Whether this was a fallback to default
	IsSubcommand bool      `json:"is_subcommand"`
	Subcommand   string    `json:"subcommand,omitempty"`
}

// DayStats represents aggregated stats for a single day
type DayStats struct {
	Date        string             `json:"date"`
	TotalUsage  int                `json:"total_usage"`
	Commands    map[string]int     `json:"commands"`
	AvgDuration map[string]float64 `json:"avg_duration"`
	TotalTime   int64              `json:"total_time_ms"`
	UniqueUsers int                `json:"unique_users"`
	TopCommands []CommandCount     `json:"top_commands"`
}

// CommandCount represents command usage count for ranking
type CommandCount struct {
	Command string `json:"command"`
	Count   int    `json:"count"`
}

// Analytics manages command usage tracking and statistics
type Analytics struct {
	logFile      string
	mutex        sync.RWMutex
	lastUsage    map[string]time.Time // Track last usage time per user for duration calculation
	sessionStart time.Time
}

// NewAnalytics creates a new analytics instance
func NewAnalytics(logFile string) *Analytics {
	return &Analytics{
		logFile:      logFile,
		lastUsage:    make(map[string]time.Time),
		sessionStart: time.Now(),
	}
}

// LogCommandUsage logs a command usage event
func (a *Analytics) LogCommandUsage(command, query, userAgent, remoteAddr string, isDefault, isSubcommand bool, subcommand string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	now := time.Now()
	userKey := remoteAddr + "|" + userAgent // Simple user identification

	usage := CommandUsage{
		Command:      command,
		Query:        query,
		Timestamp:    now,
		UserAgent:    userAgent,
		RemoteAddr:   remoteAddr,
		IsDefault:    isDefault,
		IsSubcommand: isSubcommand,
		Subcommand:   subcommand,
	}

	// Calculate duration since last command for this user
	if lastTime, exists := a.lastUsage[userKey]; exists {
		duration := now.Sub(lastTime)
		// Only log duration if it's reasonable (less than 1 hour)
		if duration < time.Hour {
			usage.Duration = duration.Milliseconds()
		}
	}
	a.lastUsage[userKey] = now

	// Log to file
	if err := a.writeLogEntry(usage); err != nil {
		fmt.Printf("Error logging command usage: %v\n", err)
	}
}

// writeLogEntry writes a single log entry to the log file
func (a *Analytics) writeLogEntry(usage CommandUsage) error {
	file, err := os.OpenFile(a.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(usage)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(data) + "\n")
	return err
}

// GetDayStats returns aggregated statistics for a specific date
func (a *Analytics) GetDayStats(date string) (*DayStats, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	entries, err := a.readLogEntries()
	if err != nil {
		return nil, err
	}

	stats := &DayStats{
		Date:        date,
		Commands:    make(map[string]int),
		AvgDuration: make(map[string]float64),
	}

	var filteredEntries []CommandUsage
	uniqueUsers := make(map[string]bool)
	commandDurations := make(map[string][]int64)

	// Filter entries for the specific date
	for _, entry := range entries {
		entryDate := entry.Timestamp.Format("2006-01-02")
		if entryDate == date {
			filteredEntries = append(filteredEntries, entry)
			userKey := entry.RemoteAddr + "|" + entry.UserAgent
			uniqueUsers[userKey] = true

			// Count command usage
			stats.Commands[entry.Command]++
			stats.TotalUsage++

			// Track durations for average calculation
			if entry.Duration > 0 {
				commandDurations[entry.Command] = append(commandDurations[entry.Command], entry.Duration)
				stats.TotalTime += entry.Duration
			}
		}
	}

	stats.UniqueUsers = len(uniqueUsers)

	// Calculate average durations
	for command, durations := range commandDurations {
		if len(durations) > 0 {
			var total int64
			for _, d := range durations {
				total += d
			}
			stats.AvgDuration[command] = float64(total) / float64(len(durations))
		}
	}

	// Generate top commands
	stats.TopCommands = a.getTopCommands(stats.Commands, 10)

	return stats, nil
}

// GetDateRange returns statistics for a range of dates
func (a *Analytics) GetDateRange(startDate, endDate string) ([]*DayStats, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	var results []*DayStats
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		stats, err := a.GetDayStats(dateStr)
		if err != nil {
			return nil, err
		}
		results = append(results, stats)
	}

	return results, nil
}

// GetOverallStats returns aggregated statistics across all time
func (a *Analytics) GetOverallStats() (*DayStats, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	entries, err := a.readLogEntries()
	if err != nil {
		return nil, err
	}

	stats := &DayStats{
		Date:        "all-time",
		Commands:    make(map[string]int),
		AvgDuration: make(map[string]float64),
	}

	uniqueUsers := make(map[string]bool)
	commandDurations := make(map[string][]int64)

	for _, entry := range entries {
		userKey := entry.RemoteAddr + "|" + entry.UserAgent
		uniqueUsers[userKey] = true

		stats.Commands[entry.Command]++
		stats.TotalUsage++

		if entry.Duration > 0 {
			commandDurations[entry.Command] = append(commandDurations[entry.Command], entry.Duration)
			stats.TotalTime += entry.Duration
		}
	}

	stats.UniqueUsers = len(uniqueUsers)

	// Calculate average durations
	for command, durations := range commandDurations {
		if len(durations) > 0 {
			var total int64
			for _, d := range durations {
				total += d
			}
			stats.AvgDuration[command] = float64(total) / float64(len(durations))
		}
	}

	stats.TopCommands = a.getTopCommands(stats.Commands, 10)

	return stats, nil
}

// readLogEntries reads all log entries from the log file
func (a *Analytics) readLogEntries() ([]CommandUsage, error) {
	file, err := os.Open(a.logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []CommandUsage{}, nil // Return empty slice if file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	var entries []CommandUsage
	decoder := json.NewDecoder(file)

	for {
		var usage CommandUsage
		if err := decoder.Decode(&usage); err != nil {
			if err.Error() == "EOF" {
				break
			}
			// Skip invalid JSON lines
			continue
		}
		entries = append(entries, usage)
	}

	return entries, nil
}

// getTopCommands returns the top N commands by usage count
func (a *Analytics) getTopCommands(commandCounts map[string]int, limit int) []CommandCount {
	var commands []CommandCount
	for cmd, count := range commandCounts {
		commands = append(commands, CommandCount{Command: cmd, Count: count})
	}

	// Simple bubble sort for top commands
	for i := 0; i < len(commands)-1; i++ {
		for j := 0; j < len(commands)-i-1; j++ {
			if commands[j].Count < commands[j+1].Count {
				commands[j], commands[j+1] = commands[j+1], commands[j]
			}
		}
	}

	if len(commands) > limit {
		commands = commands[:limit]
	}

	return commands
}
