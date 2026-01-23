package cmd

import (
	"fmt"
	"strings"
	"time"
)

// Common display utilities for consistent formatting across the CLI

// PrintHeader prints a formatted header with title
func PrintHeader(title string) {
	fmt.Printf("=== %s ===\n", title)
}

// PrintHeaderWithUser prints a formatted header with user context
func PrintHeaderWithUser(title string, user string) {
	fmt.Printf("=== %s: %s ===\n", title, user)
}

// PrintSection prints a section header with count
func PrintSection(title string, count int) {
	fmt.Printf("\n=== %s (%d) ===\n", title, count)
}

// PrintSectionSimple prints a simple section header
func PrintSectionSimple(title string) {
	fmt.Printf("\n=== %s ===\n", title)
}

// PrintListItem prints a formatted list item with icon
func PrintListItem(icon string, text string) {
	fmt.Printf("   • %s %s\n", icon, text)
}

// PrintKeyValue prints a key-value pair
func PrintKeyValue(key string, value string) {
	fmt.Printf("%s: %s\n", key, value)
}

// PrintKeyValueConditional prints a key-value pair only if value is not empty
func PrintKeyValueConditional(key string, value string) {
	if value != "" {
		fmt.Printf("%s: %s\n", key, value)
	}
}

// PrintKeyValueInt prints a key-value pair for integers
func PrintKeyValueInt(key string, value int) {
	fmt.Printf("%s: %d\n", key, value)
}

// PrintKeyValueBool prints a key-value pair for booleans
func PrintKeyValueBool(key string, value bool) {
	fmt.Printf("%s: %t\n", key, value)
}

// Date and time formatting functions

// FormatDate formats a RFC3339 date string to readable format
func FormatDate(dateStr string) string {
	if dateStr == "" {
		return "N/A"
	}

	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}

	return t.Format("2006-01-02 15:04")
}

// FormatDateTime formats a RFC3339 datetime string (alias for FormatDate for consistency)
func FormatDateTime(dateStr string) string {
	return FormatDate(dateStr)
}

// FormatEvaluationTime formats evaluation time with local timezone conversion
func FormatEvaluationTime(dateStr string) string {
	if dateStr == "" {
		return "N/A"
	}

	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}

	// Convert to local time
	localTime := t.Local()
	return localTime.Format("2006-01-02 15:04")
}

// FormatDuration formats a duration to human-readable format
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

// CalculateDuration calculates duration between two time strings
func CalculateDuration(beginAt string, endAt *string) string {
	beginTime, err := time.Parse(time.RFC3339, beginAt)
	if err != nil {
		return "N/A"
	}

	var endTime time.Time
	if endAt == nil {
		endTime = time.Now()
	} else {
		endTime, err = time.Parse(time.RFC3339, *endAt)
		if err != nil {
			return "N/A"
		}
	}

	duration := endTime.Sub(beginTime)
	return FormatDuration(duration)
}

// String utilities

// JoinUserLogins joins a slice of TeamUser logins into a comma-separated string
func JoinUserLogins(users []TeamUser) string {
	var names []string
	for _, user := range users {
		names = append(names, user.Login)
	}
	return strings.Join(names, ", ")
}

// GetProjectDisplayName returns the best display name for a project
func GetProjectDisplayName(team TeamInfo) string {
	if team.Project.Name != "" {
		return team.Project.Name
	}
	if team.Project.Slug != "" {
		return team.Project.Slug
	}
	return "Unknown Project"
}

// GetTeamDisplayName returns the best display name for a team
func GetTeamDisplayName(team TeamInfo) string {
	if team.Name != "" {
		return team.Name
	}
	return fmt.Sprintf("team-%d", team.ID)
}

// Status and result formatting

// PrintSummary prints a summary line with counts
func PrintSummary(total int, categories map[string]int) {
	PrintSectionSimple("Summary")
	fmt.Printf("Total: %d\n", total)

	var parts []string
	for name, count := range categories {
		parts = append(parts, fmt.Sprintf("%s: %d", name, count))
	}
	fmt.Println(strings.Join(parts, " | "))
}

// PrintNoResults prints a message when no results are found
func PrintNoResults(message string) {
	fmt.Println(message)
}

// PrintNoResultsDefault prints a default "no results" message
func PrintNoResultsDefault() {
	fmt.Println("No results found.")
}
