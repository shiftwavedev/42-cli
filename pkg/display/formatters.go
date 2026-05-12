package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// RelativeTime converts a timestamp to a human-readable relative time
// e.g., "2 hours ago", "3 days ago", "just now"
func RelativeTime(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp
	}

	return RelativeTimeFromTime(t)
}

// RelativeTimeFromTime converts a time.Time to a human-readable relative time
func RelativeTimeFromTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 0 {
		// Future time
		diff = -diff
		return formatDurationFuture(diff)
	}

	return formatDurationPast(diff)
}

func formatDurationPast(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case d < 365*24*time.Hour:
		months := int(d.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(d.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

func formatDurationFuture(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "in a moment"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "in 1 minute"
		}
		return fmt.Sprintf("in %d minutes", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "in 1 hour"
		}
		return fmt.Sprintf("in %d hours", hours)
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "in 1 day"
		}
		return fmt.Sprintf("in %d days", days)
	default:
		months := int(d.Hours() / 24 / 30)
		if months == 1 {
			return "in 1 month"
		}
		return fmt.Sprintf("in %d months", months)
	}
}

// FormatDuration formats a duration for display
// e.g., "3 hours", "45 minutes", "2 days"
func FormatDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "less than a minute"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
}

// ExtractEmailDomain extracts the domain from an email address
func ExtractEmailDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return email
}

// TruncateString truncates a string to a maximum length with ellipsis (ANSI-aware)
func TruncateString(s string, maxLen int) string {
	visibleWidth := lipgloss.Width(s)
	if visibleWidth <= maxLen {
		return s
	}
	if maxLen <= 3 {
		// For very short limits, use lipgloss to handle ANSI codes properly
		return lipgloss.NewStyle().MaxWidth(maxLen).Render(s)
	}
	// Truncate with ellipsis, preserving ANSI codes
	return lipgloss.NewStyle().MaxWidth(maxLen-3).Render(s) + "..."
}

// PadRight pads a string to a minimum width (using visible width, ignoring ANSI codes)
func PadRight(s string, width int) string {
	visibleWidth := lipgloss.Width(s)
	if visibleWidth >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visibleWidth)
}

// PadLeft pads a string to a minimum width on the left (using visible width, ignoring ANSI codes)
func PadLeft(s string, width int) string {
	visibleWidth := lipgloss.Width(s)
	if visibleWidth >= width {
		return s
	}
	return strings.Repeat(" ", width-visibleWidth) + s
}

// MaskToken masks a token showing only first and last 4 characters
// e.g., "abcd...wxyz"
func MaskToken(token string) string {
	if len(token) <= 12 {
		return token
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// JoinWithComma joins strings with commas and "and" for the last item
func JoinWithComma(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " and " + items[1]
	default:
		return strings.Join(items[:len(items)-1], ", ") + ", and " + items[len(items)-1]
	}
}

// Pluralize returns singular or plural form based on count
func Pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// FormatCount formats a count with its label
// e.g., "1 project", "5 projects"
func FormatCount(count int, singular, plural string) string {
	return fmt.Sprintf("%d %s", count, Pluralize(count, singular, plural))
}

// CapitalizeFirst capitalizes the first letter of a string
func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FormatLevel formats a float level with 2 decimal places
func FormatLevel(level float64) string {
	return fmt.Sprintf("%.2f", level)
}

// AutoPadWidth returns responsive padding width based on terminal size
// Use 25% of terminal width, bounded by min/max for readability
func AutoPadWidth() int {
	width := GlobalTerminalDimensions.Width / 4
	if width < 15 {
		return 15
	}
	if width > 30 {
		return 30
	}
	return width
}

// CalculateMaxProjectNameLength calculates safe truncation point based on terminal size
// Allocate ~30% of terminal width to project name, bounded for readability
func CalculateMaxProjectNameLength() int {
	available := GlobalTerminalDimensions.Width / 3
	if available < 15 {
		return 15
	}
	if available > 40 {
		return 40
	}
	return available
}

// WrapTextResponsive wraps text with responsive width
// reservedWidth is the space needed for other columns/formatting
func WrapTextResponsive(text string, reservedWidth int) string {
	maxWidth := GlobalTerminalDimensions.Width - reservedWidth
	if maxWidth < 20 {
		maxWidth = 20
	}
	if len(text) > maxWidth {
		return text[:maxWidth-3] + "..."
	}
	return text
}

// DividerForContent creates a section divider that matches the width of rendered content
func DividerForContent(title string, renderedContent string) string {
	contentWidth := lipgloss.Width(renderedContent)
	if contentWidth <= 0 {
		return SectionDivider(title)
	}
	return SectionDividerWithWidth(title, contentWidth)
}

// SectionDividerWithWidth creates a section divider with explicit width
func SectionDividerWithWidth(title string, width int) string {
	if width <= 0 {
		return SectionDivider(title)
	}
	sideWidth := max((width-len(title)-2)/2, 1)

	left := strings.Repeat(DividerChar, sideWidth)
	right := strings.Repeat(DividerChar, width-len(title)-2-sideWidth)

	return lipgloss.NewStyle().
		Foreground(Border).
		Render(left + " " + title + " " + right)
}
