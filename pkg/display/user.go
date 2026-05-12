package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
)

// UserProfile represents a 42 user profile for rendering
type UserProfile struct {
	Login           string
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	DisplayName     string
	Staff           bool
	CorrectionPoint int
	PoolMonth       string
	PoolYear        string
	Location        string
	Wallet          int
	CreatedAt       string
	UpdatedAt       string
	Alumni          bool
	Active          bool
	Campus          []CampusInfo
	CursusUsers     []CursusUser
}

// CampusInfo represents campus information
type CampusInfo struct {
	Name    string
	Country string
}

// CursusUser represents a user's cursus information
type CursusUser struct {
	Grade  string
	Level  float64
	Cursus CursusInfo
}

// CursusInfo represents cursus details
type CursusInfo struct {
	Name string
}

// RenderUserProfile renders a user profile with clean, styled output
func RenderUserProfile(user *UserProfile) string {
	var b strings.Builder

	// Header: login + email domain
	emailDomain := ""
	if user.Email != "" {
		emailDomain = "@" + ExtractEmailDomain(user.Email)
	}
	b.WriteString(Header(user.Login, emailDomain))
	b.WriteString("\n")
	b.WriteString(Divider(0))
	b.WriteString("\n\n")

	// Status line with semantic indicators and badges
	status := []string{}

	// Active/Inactive status with icon
	if user.Active {
		status = append(status, Icon("online")+" Active")
	} else {
		status = append(status, Icon("offline")+" Offline")
	}

	// Staff/Alumni pills
	if user.Staff {
		status = append(status, Badge("Staff", Primary))
	}
	if user.Alumni {
		status = append(status, Badge("Alumni", Secondary))
	}

	// Correction points and wallet
	status = append(status, fmt.Sprintf("%d correction pts", user.CorrectionPoint))
	status = append(status, fmt.Sprintf("%d₳", user.Wallet))

	b.WriteString(Indent + strings.Join(status, "   ") + "\n\n")

	// Location (if present)
	if user.Location != "" {
		b.WriteString(fmt.Sprintf("%s%s %s\n", Indent, RenderIf(Subtle, "Location"), user.Location))
	}

	// Pool info
	if user.PoolMonth != "" && user.PoolYear != "" {
		poolInfo := CapitalizeFirst(user.PoolMonth) + " " + user.PoolYear
		b.WriteString(fmt.Sprintf("%s%s %s\n", Indent, RenderIf(Subtle, "Piscine"), poolInfo))
	}

	// Campus list
	if len(user.Campus) > 0 {
		b.WriteString("\n")
		campusList := ""
		for _, campus := range user.Campus {
			campusList += ListItem(fmt.Sprintf("%s, %s", campus.Name, campus.Country)) + "\n"
		}
		b.WriteString(DividerForContent("CAMPUS", campusList))
		b.WriteString("\n\n")
		b.WriteString(campusList)
	}

	// Cursus list with progress bars
	if len(user.CursusUsers) > 0 {
		if len(user.Campus) == 0 {
			b.WriteString("\n")
		}
		cursusList := renderCursusList(user.CursusUsers)
		b.WriteString(DividerForContent("CURSUS", cursusList))
		b.WriteString("\n\n")
		b.WriteString(cursusList)
	}

	created := RelativeTime(user.CreatedAt)
	updated := RelativeTime(user.UpdatedAt)
	footer := fmt.Sprintf("Created %s · Updated %s", created, updated)
	b.WriteString(Indent + RenderIf(Subtle, footer) + "\n")

	return b.String()
}

// renderProgressBar creates a visual progress bar using bubbles/progress
// In NO_COLOR mode, renders a plain ASCII bar without ANSI codes
func renderProgressBar(level float64, width int) string {
	// Level is out of 21
	percent := level / 21.0
	if percent > 1.0 {
		percent = 1.0
	}

	if NoColor {
		// NO_COLOR mode: render plain ASCII bar
		filledWidth := int(float64(width) * percent)
		emptyWidth := width - filledWidth

		// Use plain ASCII characters for the bar
		bar := strings.Repeat("#", filledWidth) + strings.Repeat("-", emptyWidth)

		// Add percentage
		percentStr := fmt.Sprintf("  %d%%  %.2f / 21", int(percent*100), level)
		return "[" + bar + "]" + percentStr
	}

	// COLOR mode: Use bubbles progress bar with gradient
	prog := progress.New(
		progress.WithScaledGradient("#00d7ff", "#0087ff"),
		progress.WithWidth(width),
	)

	// Render the progress bar
	barStr := prog.ViewAs(percent)

	// Add the level info at the end
	levelStr := fmt.Sprintf("  %.2f / 21", level)
	return barStr + levelStr
}

// LocationInfo represents user location for rendering
type LocationInfo struct {
	Host    string
	BeginAt string
	EndAt   *string
	Primary bool
}

// RenderLocationInfo renders user location information
func RenderLocationInfo(login string, locations []LocationInfo) string {
	var b strings.Builder

	b.WriteString(Header("Location", login))
	b.WriteString("\n")
	b.WriteString(Divider(0))
	b.WriteString("\n\n")

	if len(locations) == 0 {
		b.WriteString(Indent + RenderIf(Subtle, "No location information found.") + "\n")
		return b.String()
	}

	// Find active locations
	var active []LocationInfo
	for _, loc := range locations {
		if loc.EndAt == nil {
			active = append(active, loc)
		}
	}

	if len(active) > 0 {
		sessionContent := ""
		for _, loc := range active {
			// Parse begin time for duration calculation
			beginTime, _ := time.Parse(time.RFC3339, loc.BeginAt)
			duration := FormatDuration(time.Since(beginTime))

			// Build location panel body
			locationBody := fmt.Sprintf(
				"%s  %s\n%s  %s\n%s  %s",
				RenderIf(Subtle, "Host"),
				loc.Host,
				RenderIf(Subtle, "Since"),
				RelativeTime(loc.BeginAt),
				RenderIf(Subtle, "Duration"),
				duration,
			)
			sessionContent += Panel("Session", locationBody) + "\n"
		}
		b.WriteString(DividerForContent("CURRENT SESSION", sessionContent))
		b.WriteString("\n\n")
		b.WriteString(sessionContent)
	} else {
		// Offline status
		b.WriteString(Icon("offline") + " Offline\n")
	}

	return b.String()
}

// renderCursusList renders the cursus list as a single string for divider width calculation
func renderCursusList(cursusUsers []CursusUser) string {
	var result strings.Builder

	for _, cursus := range cursusUsers {
		name := PadRight(cursus.Cursus.Name, AutoPadWidth())
		levelStr := FormatLevel(cursus.Level)
		if cursus.Grade != "" {
			levelStr += RenderIf(Subtle, fmt.Sprintf(" (%s)", cursus.Grade))
		}
		result.WriteString(name + "  " + levelStr + "\n")

		bar := renderProgressBar(cursus.Level, 30)
		if NoColor {
			percentage := (cursus.Level / 21.0) * 100.0
			result.WriteString(fmt.Sprintf("%s%s  (%.1f%%)\n", Indent, bar, percentage))
		} else {
			result.WriteString(Indent + bar + "\n")
		}
		result.WriteString("\n")
	}

	return result.String()
}
