package display

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// NoColor indicates whether colors should be disabled
var NoColor bool

func init() {
	// Detect NO_COLOR env var or CI environment
	if os.Getenv("NO_COLOR") != "" || os.Getenv("CI") != "" {
		NoColor = true
	}
}

// Color palette - minimalist & accessible (GitHub CLI inspired)
var (
	// Primary colors - subtle, professional
	Primary   = lipgloss.Color("39")  // Soft blue (similar to gh)
	Secondary = lipgloss.Color("246") // Cool gray
	Success   = lipgloss.Color("2")   // Green (finished, valid)
	Warning   = lipgloss.Color("3")   // Yellow (expiring, pending)
	Error     = lipgloss.Color("1")   // Red (failed, expired)
	Muted     = lipgloss.Color("240") // Dim gray (metadata)
	Text      = lipgloss.Color("252") // Light gray (primary text)
)

// Typography styles - clean, readable
var (
	// H1 - Main headers
	H1 = lipgloss.NewStyle().
		Bold(true).
		Foreground(Text)

	// H2 - Section headers
	H2 = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		MarginTop(1)

	// Body - Regular text
	Body = lipgloss.NewStyle().
		Foreground(Text)

	// Subtle - Dimmed text for metadata
	Subtle = lipgloss.NewStyle().
		Foreground(Muted).
		Italic(true)

	// Code - Monospace styled text
	Code = lipgloss.NewStyle().
		Foreground(Secondary)

	// ErrorText - Error messages
	ErrorText = lipgloss.NewStyle().
		Foreground(Error).
		Bold(true)

	// SuccessText - Success messages
	SuccessText = lipgloss.NewStyle().
		Foreground(Success)

	// WarningText - Warning messages
	WarningText = lipgloss.NewStyle().
		Foreground(Warning)
)

// Layout elements
var (
	// Divider character for separating sections
	DividerChar = "─"

	// Indent for nested content
	Indent = "  "
)

// ListPrefix returns the list bullet character, styled if colors are enabled
func ListPrefix() string {
	if NoColor {
		return "•"
	}
	return lipgloss.NewStyle().
		Foreground(Primary).
		Render("•")
}

// Badge renders a colored badge/tag
func Badge(text string, color lipgloss.Color) string {
	if NoColor {
		return text
	}
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(text)
}

// StatusBadge renders a status indicator with appropriate color
func StatusBadge(status string) string {
	switch status {
	case "finished", "completed", "valid", "active":
		return Badge(status, Success)
	case "in_progress", "pending", "expiring":
		return Badge(status, Warning)
	case "failed", "expired", "error":
		return Badge(status, Error)
	default:
		return Badge(status, Secondary)
	}
}

// Divider returns a horizontal divider line
func Divider(width int) string {
	if NoColor {
		return repeat(DividerChar, width)
	}
	return lipgloss.NewStyle().
		Foreground(Muted).
		Render(repeat(DividerChar, width))
}

// repeat returns a string repeated n times
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// RenderIf returns styled text if colors are enabled, otherwise plain text
func RenderIf(style lipgloss.Style, text string) string {
	if NoColor {
		return text
	}
	return style.Render(text)
}

// Header renders a section header with optional subtitle
func Header(title string, subtitle string) string {
	result := RenderIf(H1, title)
	if subtitle != "" {
		result += "  " + RenderIf(Subtle, subtitle)
	}
	return result
}

// SectionHeader renders a secondary section header
func SectionHeader(title string) string {
	return RenderIf(H2, title)
}

// ListItem renders a bulleted list item
func ListItem(text string) string {
	return Indent + ListPrefix() + " " + text
}

// KeyValue renders a key-value pair with aligned formatting
func KeyValue(key string, value string, keyWidth int) string {
	paddedKey := key
	for len(paddedKey) < keyWidth {
		paddedKey += " "
	}
	return Indent + RenderIf(Subtle, paddedKey) + "  " + value
}
