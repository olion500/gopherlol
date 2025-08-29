package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/markusdosch/gopherlol/internal/analytics"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	analytics *analytics.Analytics
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(analytics *analytics.Analytics) *DashboardHandler {
	return &DashboardHandler{
		analytics: analytics,
	}
}

// HandleDashboard serves the main dashboard page
func (d *DashboardHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get date parameter, default to today
	dateParam := r.URL.Query().Get("date")
	if dateParam == "" {
		dateParam = time.Now().Format("2006-01-02")
	}

	// Get date range if specified
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `<!DOCTYPE html>
<html>
<head>
    <title>gopherlol Analytics Dashboard</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #333; margin: 0; }
        .header p { color: #666; margin: 10px 0; }
        .controls { margin-bottom: 20px; padding: 15px; background: #f8f9fa; border-radius: 5px; }
        .controls input, .controls button { margin: 5px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
        .controls button { background: #007bff; color: white; cursor: pointer; }
        .controls button:hover { background: #0056b3; }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-bottom: 20px; }
        .stat-card { background: white; padding: 20px; border-radius: 8px; border-left: 4px solid #007bff; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        .stat-card h3 { margin: 0 0 10px 0; color: #333; }
        .stat-value { font-size: 2em; font-weight: bold; color: #007bff; margin: 10px 0; }
        .chart-container { margin-top: 20px; }
        .command-list { list-style: none; padding: 0; }
        .command-item { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #eee; }
        .command-name { font-weight: bold; }
        .command-count { color: #666; }
        .duration { color: #28a745; font-size: 0.9em; }
        .refresh-btn { float: right; }
        .no-data { text-align: center; color: #666; padding: 40px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 gopherlol Analytics Dashboard</h1>
            <p>Command usage statistics and insights</p>
        </div>
        
        <div class="controls">
            <form method="GET" style="display: inline;">
                <label>Date:</label>
                <input type="date" name="date" value="` + dateParam + `">
                <button type="submit">View Day</button>
            </form>
            
            <form method="GET" style="display: inline; margin-left: 20px;">
                <label>Range:</label>
                <input type="date" name="start" placeholder="Start date" value="` + startDate + `">
                <input type="date" name="end" placeholder="End date" value="` + endDate + `">
                <button type="submit">View Range</button>
            </form>
            
            <button class="refresh-btn" onclick="window.location.reload()">🔄 Refresh</button>
        </div>

        <div id="stats-content">
            <div class="no-data">Loading statistics...</div>
        </div>
    </div>

    <script>
        // Load statistics based on current URL parameters
        window.onload = function() {
            const urlParams = new URLSearchParams(window.location.search);
            const date = urlParams.get('date') || '` + time.Now().Format("2006-01-02") + `';
            const start = urlParams.get('start');
            const end = urlParams.get('end');
            
            let apiUrl = '/api/stats';
            if (start && end) {
                apiUrl += '?start=' + start + '&end=' + end;
            } else {
                apiUrl += '?date=' + date;
            }
            
            fetch(apiUrl)
                .then(response => response.json())
                .then(data => displayStats(data))
                .catch(error => {
                    document.getElementById('stats-content').innerHTML = 
                        '<div class="no-data">Error loading statistics: ' + error.message + '</div>';
                });
        };

        function displayStats(data) {
            const container = document.getElementById('stats-content');
            
            if (Array.isArray(data)) {
                // Range data
                let totalUsage = 0;
                let totalTime = 0;
                let allCommands = {};
                
                data.forEach(day => {
                    totalUsage += day.total_usage;
                    totalTime += day.total_time_ms;
                    Object.keys(day.commands).forEach(cmd => {
                        allCommands[cmd] = (allCommands[cmd] || 0) + day.commands[cmd];
                    });
                });
                
                container.innerHTML = generateStatsHTML({
                    date: data.length > 0 ? data[0].date + ' to ' + data[data.length-1].date : 'No data',
                    total_usage: totalUsage,
                    total_time_ms: totalTime,
                    commands: allCommands,
                    unique_users: Math.max(...data.map(d => d.unique_users || 0)),
                    top_commands: Object.keys(allCommands).map(cmd => ({command: cmd, count: allCommands[cmd]}))
                        .sort((a,b) => b.count - a.count).slice(0, 10)
                });
            } else {
                // Single day data
                container.innerHTML = generateStatsHTML(data);
            }
        }

        function generateStatsHTML(stats) {
            const topCommands = stats.top_commands || [];
            const totalTimeHours = (stats.total_time_ms / (1000 * 60 * 60)).toFixed(1);
            
            return ` + "`" + `
                <div class="stats-grid">
                    <div class="stat-card">
                        <h3>📊 Total Usage</h3>
                        <div class="stat-value">${stats.total_usage || 0}</div>
                        <div>commands executed</div>
                    </div>
                    <div class="stat-card">
                        <h3>👥 Unique Users</h3>
                        <div class="stat-value">${stats.unique_users || 0}</div>
                        <div>different users</div>
                    </div>
                    <div class="stat-card">
                        <h3>⏱️ Total Time</h3>
                        <div class="stat-value">${totalTimeHours}h</div>
                        <div>time spent</div>
                    </div>
                    <div class="stat-card">
                        <h3>🔝 Top Command</h3>
                        <div class="stat-value">${topCommands[0]?.command || 'N/A'}</div>
                        <div>${topCommands[0]?.count || 0} uses</div>
                    </div>
                </div>
                
                <div class="stat-card">
                    <h3>📈 Command Usage (${stats.date || 'Unknown'})</h3>
                    <ul class="command-list">
                        ${topCommands.map(cmd => ` + "`" + `
                            <li class="command-item">
                                <span class="command-name">${cmd.command}</span>
                                <span>
                                    <span class="command-count">${cmd.count} uses</span>
                                    ${stats.avg_duration?.[cmd.command] ? 
                                        ` + "`" + ` • <span class="duration">${(stats.avg_duration[cmd.command] / 1000).toFixed(1)}s avg</span>` + "`" + ` 
                                        : ''}
                                </span>
                            </li>
                        ` + "`" + `).join('')}
                    </ul>
                </div>
            ` + "`" + `;
        }
    </script>
</body>
</html>`

	_, _ = w.Write([]byte(html))
}

// HandleStatsAPI serves JSON statistics data
func (d *DashboardHandler) HandleStatsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Check for date range query
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")

	if startDate != "" && endDate != "" {
		// Return range data
		stats, err := d.analytics.GetDateRange(startDate, endDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting date range stats: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return
	}

	// Get single date parameter
	dateParam := r.URL.Query().Get("date")
	if dateParam == "" {
		dateParam = time.Now().Format("2006-01-02")
	}

	stats, err := d.analytics.GetDayStats(dateParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting day stats: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// HandleOverallStatsAPI serves overall statistics
func (d *DashboardHandler) HandleOverallStatsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	stats, err := d.analytics.GetOverallStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting overall stats: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}