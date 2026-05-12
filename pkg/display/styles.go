package display

import (
	"os"
	"strings"

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

// Color palette - adaptive for light and dark terminal backgrounds
// TODO: Change that's not really adaptive :/ I don't found the best colors for light/dark terminals
var (
	Primary   = lipgloss.AdaptiveColor{Light: "27", Dark: "39"}   // Bright blue on light, soft blue on dark
	Secondary = lipgloss.AdaptiveColor{Light: "243", Dark: "246"} // Dark gray on light, light gray on dark
	Success   = lipgloss.AdaptiveColor{Light: "28", Dark: "2"}    // Dark green on light, bright green on dark
	Warning   = lipgloss.AdaptiveColor{Light: "130", Dark: "3"}   // Brown/orange on light, bright yellow on dark
	Error     = lipgloss.AdaptiveColor{Light: "160", Dark: "1"}   // Dark red on light, bright red on dark
	Muted     = lipgloss.AdaptiveColor{Light: "245", Dark: "240"} // Mid-gray for both
	Text      = lipgloss.AdaptiveColor{Light: "235", Dark: "252"} // Near-black on light, near-white on dark
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

// Icon returns a semantic icon symbol (✓, ✗, ⚠, ●, ○, ⚙)
// In NO_COLOR mode, lipgloss automatically strips the color codes
func Icon(state string) string {
	state = strings.ToLower(state)

	icons := map[string]struct {
		symbol string
		color  lipgloss.TerminalColor
	}{
		"success":     {symbol: "✓", color: Success},
		"error":       {symbol: "✗", color: Error},
		"warning":     {symbol: "⚠", color: Warning},
		"online":      {symbol: "●", color: Success},
		"offline":     {symbol: "○", color: Muted},
		"in_progress": {symbol: "⚙", color: Primary},
	}

	ic, exists := icons[state]
	if !exists {
		// Fallback for unknown states
		return state
	}

	// Render symbol with color - in NO_COLOR mode, lipgloss will strip the ANSI codes
	return lipgloss.NewStyle().Foreground(ic.color).Render(ic.symbol)
}

// ListPrefix returns the list bullet character, styled
func ListPrefix() string {
	return lipgloss.NewStyle().
		Foreground(Primary).
		Render("•")
}

// Badge renders a colored badge/pill with text
// TerminalColor interface accepts both Color and AdaptiveColor
func Badge(text string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().
		Background(color).
		Foreground(Text).
		Bold(true).
		Padding(0, 1).
		Render(text)
}

// StatusBadge renders a status indicator with appropriate color and icon
func StatusBadge(status string) string {
	status = strings.ToLower(status)

	statusMap := map[string]struct {
		icon  string
		label string
	}{
		"finished":               {icon: "success", label: "Finished"},
		"completed":              {icon: "success", label: "Completed"},
		"valid":                  {icon: "success", label: "Valid"},
		"active":                 {icon: "online", label: "Active"},
		"waiting_for_correction": {icon: "warning", label: "Waiting"},
		"in_progress":            {icon: "in_progress", label: "In Progress"},
		"failed":                 {icon: "error", label: "Failed"},
		"expired":                {icon: "error", label: "Expired"},
		"error":                  {icon: "error", label: "Error"},
	}

	sm, exists := statusMap[status]
	if !exists {
		return Badge(status, Secondary)
	}

	icon := Icon(sm.icon)
	label := sm.label

	// Render icon + styled label - in NO_COLOR mode, lipgloss strips ANSI
	var labelStyle lipgloss.Style
	switch sm.icon {
	case "success":
		labelStyle = lipgloss.NewStyle().Foreground(Success)
	case "error":
		labelStyle = lipgloss.NewStyle().Foreground(Error)
	case "warning":
		labelStyle = lipgloss.NewStyle().Foreground(Warning)
	case "online":
		labelStyle = lipgloss.NewStyle().Foreground(Success)
	case "in_progress":
		labelStyle = lipgloss.NewStyle().Foreground(Primary)
	default:
		labelStyle = lipgloss.NewStyle().Foreground(Secondary)
	}

	return icon + " " + labelStyle.Render(label)
}

// Panel renders a rounded border panel with title and body
func Panel(title, body string) string {
	// COLOR mode: lipgloss rounded border
	// In NO_COLOR mode, lipgloss strips ANSI but keeps Unicode box-drawing chars
	titleStyle := lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Muted).
		Padding(1, 2).
		Render(titleStyle.Render(title) + "\n" + body)
}

// DividerWidth returns responsive divider width with margins
func DividerWidth() int {
	width := GlobalTerminalDimensions.Width
	if width > 24 {
		return width - 4
	}
	return 20
}

// SectionDividerWidth returns width for section headers (slightly narrower)
func SectionDividerWidth() int {
	width := DividerWidth()
	if width > 8 {
		return width - 4
	}
	return width
}

// SectionDivider renders a section header line with title
func SectionDivider(title string) string {
	width := SectionDividerWidth()
	sideWidth := max((width-len(title)-2)/2, 1)

	left := strings.Repeat("─", sideWidth)
	right := strings.Repeat("─", width-len(title)-2-sideWidth)

	return lipgloss.NewStyle().
		Foreground(Muted).
		Render(left + " " + title + " " + right)
}

// Divider returns a horizontal divider line with responsive width
// Pass width <= 0 to auto-detect from terminal width
func Divider(width int) string {
	if width <= 0 {
		width = DividerWidth()
	}

	dividerLine := strings.Repeat(DividerChar, width)

	return lipgloss.NewStyle().
		Foreground(Muted).
		Render(dividerLine)
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

// SectionHeader renders a secondary section header (deprecated: use SectionDivider instead)
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
