package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/dominikoh/gopherlol/internal/analytics"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

func main() {
	var (
		date      = flag.String("date", "", "Show stats for specific date (YYYY-MM-DD)")
		startDate = flag.String("start", "", "Start date for range (YYYY-MM-DD)")
		endDate   = flag.String("end", "", "End date for range (YYYY-MM-DD)")
		logFile   = flag.String("log", "usage.log", "Path to usage log file")
		top       = flag.Int("top", 10, "Number of top commands to show")
		overall   = flag.Bool("overall", false, "Show overall statistics")
		help      = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Initialize analytics
	analyticsSystem := analytics.NewAnalytics(*logFile)

	// Check if log file exists
	if _, err := os.Stat(*logFile); os.IsNotExist(err) {
		fmt.Printf("%s❌ No analytics data found%s\n", ColorRed, ColorReset)
		fmt.Printf("   Log file '%s' does not exist.\n", *logFile)
		fmt.Printf("   Start using gopherlol to generate analytics data!\n\n")
		return
	}

	printHeader()

	if *overall {
		showOverallStats(analyticsSystem, *top)
	} else if *startDate != "" && *endDate != "" {
		showDateRangeStats(analyticsSystem, *startDate, *endDate, *top)
	} else {
		targetDate := *date
		if targetDate == "" {
			targetDate = time.Now().Format("2006-01-02")
		}
		showDayStats(analyticsSystem, targetDate, *top)
	}
}

func showHelp() {
	fmt.Printf("%s📊 gopherlol Analytics CLI%s\n\n", ColorBold+ColorCyan, ColorReset)
	fmt.Println("Display command usage analytics in the terminal.")
	fmt.Println()
	fmt.Printf("%sUsage:%s\n", ColorBold, ColorReset)
	fmt.Println("  analytics [options]")
	fmt.Println()
	fmt.Printf("%sOptions:%s\n", ColorBold, ColorReset)
	fmt.Println("  -date YYYY-MM-DD    Show stats for specific date (default: today)")
	fmt.Println("  -start YYYY-MM-DD   Start date for range analysis")
	fmt.Println("  -end YYYY-MM-DD     End date for range analysis")
	fmt.Println("  -log FILE          Path to usage log file (default: usage.log)")
	fmt.Println("  -top N             Number of top commands to show (default: 10)")
	fmt.Println("  -overall           Show overall/all-time statistics")
	fmt.Println("  -help              Show this help message")
	fmt.Println()
	fmt.Printf("%sExamples:%s\n", ColorBold, ColorReset)
	fmt.Println("  analytics                    # Today's stats")
	fmt.Println("  analytics -date 2024-01-15   # Specific date")
	fmt.Println("  analytics -overall           # All-time stats")
	fmt.Println("  analytics -start 2024-01-01 -end 2024-01-31  # Range")
	fmt.Println("  analytics -top 5             # Show top 5 commands only")
}

func printHeader() {
	fmt.Printf("%s", ColorBold+ColorCyan)
	fmt.Println("┌─────────────────────────────────────────────────┐")
	fmt.Println("│           🔍 gopherlol Analytics                │")
	fmt.Println("└─────────────────────────────────────────────────┘")
	fmt.Printf("%s\n", ColorReset)
}

func showDayStats(analyticsSystem *analytics.Analytics, date string, topCount int) {
	stats, err := analyticsSystem.GetDayStats(date)
	if err != nil {
		fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("%s📅 Statistics for %s%s\n\n", ColorBold+ColorBlue, date, ColorReset)

	if stats.TotalUsage == 0 {
		fmt.Printf("%s📭 No command usage data for this date%s\n", ColorYellow, ColorReset)
		return
	}

	showStatsOverview(stats)
	showTopCommands(stats.TopCommands, stats.AvgDuration, topCount)
}

func showDateRangeStats(analyticsSystem *analytics.Analytics, startDate, endDate string, topCount int) {
	rangeStats, err := analyticsSystem.GetDateRange(startDate, endDate)
	if err != nil {
		fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	if len(rangeStats) == 0 {
		fmt.Printf("%s📭 No data found for the specified date range%s\n", ColorYellow, ColorReset)
		return
	}

	// Aggregate range data
	totalUsage := 0
	totalTime := int64(0)
	allCommands := make(map[string]int)
	allDurations := make(map[string][]int64)
	maxUsers := 0

	for _, dayStats := range rangeStats {
		totalUsage += dayStats.TotalUsage
		totalTime += dayStats.TotalTime
		if dayStats.UniqueUsers > maxUsers {
			maxUsers = dayStats.UniqueUsers
		}

		for cmd, count := range dayStats.Commands {
			allCommands[cmd] += count
		}

		// Collect duration data (simplified aggregation)
		for cmd, avgDur := range dayStats.AvgDuration {
			if avgDur > 0 {
				allDurations[cmd] = append(allDurations[cmd], int64(avgDur))
			}
		}
	}

	// Calculate average durations
	avgDurations := make(map[string]float64)
	for cmd, durations := range allDurations {
		if len(durations) > 0 {
			var total int64
			for _, d := range durations {
				total += d
			}
			avgDurations[cmd] = float64(total) / float64(len(durations))
		}
	}

	// Create aggregated stats
	aggregatedStats := &analytics.DayStats{
		Date:        fmt.Sprintf("%s to %s", startDate, endDate),
		TotalUsage:  totalUsage,
		Commands:    allCommands,
		AvgDuration: avgDurations,
		TotalTime:   totalTime,
		UniqueUsers: maxUsers,
		TopCommands: getTopCommandsFromMap(allCommands, topCount),
	}

	fmt.Printf("%s📊 Statistics for %s%s\n\n", ColorBold+ColorBlue, aggregatedStats.Date, ColorReset)
	showStatsOverview(aggregatedStats)
	showTopCommands(aggregatedStats.TopCommands, aggregatedStats.AvgDuration, topCount)
}

func showOverallStats(analyticsSystem *analytics.Analytics, topCount int) {
	stats, err := analyticsSystem.GetOverallStats()
	if err != nil {
		fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("%s🌟 Overall Statistics (All Time)%s\n\n", ColorBold+ColorPurple, ColorReset)

	if stats.TotalUsage == 0 {
		fmt.Printf("%s📭 No command usage data found%s\n", ColorYellow, ColorReset)
		return
	}

	showStatsOverview(stats)
	showTopCommands(stats.TopCommands, stats.AvgDuration, topCount)
}

func showStatsOverview(stats *analytics.DayStats) {
	totalTimeHours := float64(stats.TotalTime) / (1000 * 60 * 60)

	fmt.Printf("%s┌─ Overview ─────────────────────────────────────┐%s\n", ColorDim, ColorReset)
	fmt.Printf("│ %s📊 Total Usage:%s     %-25d │\n", ColorGreen, ColorReset, stats.TotalUsage)
	fmt.Printf("│ %s👥 Unique Users:%s    %-25d │\n", ColorCyan, ColorReset, stats.UniqueUsers)
	fmt.Printf("│ %s⏱️  Total Time:%s     %-22.1fh │\n", ColorYellow, ColorReset, totalTimeHours)

	if len(stats.TopCommands) > 0 {
		topCmd := stats.TopCommands[0]
		fmt.Printf("│ %s🔝 Top Command:%s     %-15s (%d uses) │\n", ColorPurple, ColorReset, topCmd.Command, topCmd.Count)
	}

	fmt.Printf("%s└───────────────────────────────────────────────┘%s\n\n", ColorDim, ColorReset)
}

func showTopCommands(commands []analytics.CommandCount, avgDurations map[string]float64, limit int) {
	if len(commands) == 0 {
		return
	}

	if len(commands) > limit {
		commands = commands[:limit]
	}

	fmt.Printf("%s📈 Top Commands%s\n", ColorBold+ColorBlue, ColorReset)
	fmt.Printf("%s┌─────────────────────────────────────────────────────────────┐%s\n", ColorDim, ColorReset)
	fmt.Printf("│ %s%-15s %8s %12s %12s%s │\n", ColorBold, "Command", "Count", "Percentage", "Avg Duration", ColorReset)
	fmt.Printf("├─────────────────────────────────────────────────────────────┤\n")

	// Calculate total for percentages
	total := 0
	for _, cmd := range commands {
		total += cmd.Count
	}

	for i, cmd := range commands {
		percentage := float64(cmd.Count) / float64(total) * 100
		avgDur := avgDurations[cmd.Command]

		// Color coding for ranking
		rankColor := ColorWhite
		switch i {
		case 0:
			rankColor = ColorGreen + ColorBold
		case 1:
			rankColor = ColorYellow + ColorBold
		case 2:
			rankColor = ColorRed + ColorBold
		}

		durationStr := "-"
		if avgDur > 0 {
			if avgDur < 1000 {
				durationStr = fmt.Sprintf("%.0fms", avgDur)
			} else {
				durationStr = fmt.Sprintf("%.1fs", avgDur/1000)
			}
		}

		fmt.Printf("│ %s%-15s%s %7d %11.1f%% %11s │\n",
			rankColor, cmd.Command, ColorReset, cmd.Count, percentage, durationStr)
	}

	fmt.Printf("%s└─────────────────────────────────────────────────────────────┘%s\n", ColorDim, ColorReset)
}

func getTopCommandsFromMap(commands map[string]int, limit int) []analytics.CommandCount {
	var commandList []analytics.CommandCount
	for cmd, count := range commands {
		commandList = append(commandList, analytics.CommandCount{Command: cmd, Count: count})
	}

	// Sort by count (descending)
	sort.Slice(commandList, func(i, j int) bool {
		return commandList[i].Count > commandList[j].Count
	})

	if len(commandList) > limit {
		commandList = commandList[:limit]
	}

	return commandList
}
