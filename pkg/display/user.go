package display

import (
	"fmt"
	"strings"
	"time"
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
	b.WriteString(Divider(65))
	b.WriteString("\n\n")

	// Status line with inline badges
	status := []string{}
	if user.Active {
		status = append(status, Badge("Active", Success))
	} else {
		status = append(status, Badge("Inactive", Muted))
	}
	if user.Staff {
		status = append(status, Badge("Staff", Primary))
	}
	if user.Alumni {
		status = append(status, Badge("Alumni", Secondary))
	}
	status = append(status, fmt.Sprintf("%d correction points", user.CorrectionPoint))
	status = append(status, fmt.Sprintf("%d wallet", user.Wallet))

	b.WriteString(Indent + strings.Join(status, " · ") + "\n\n")

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
		b.WriteString(SectionHeader("CAMPUS"))
		b.WriteString("\n")
		for _, campus := range user.Campus {
			b.WriteString(ListItem(fmt.Sprintf("%s, %s", campus.Name, campus.Country)))
			b.WriteString("\n")
		}
	}

	// Cursus list with aligned levels
	if len(user.CursusUsers) > 0 {
		if len(user.Campus) == 0 {
			b.WriteString("\n")
		}
		b.WriteString(SectionHeader("CURSUS"))
		b.WriteString("\n")
		for _, cursus := range user.CursusUsers {
			name := PadRight(cursus.Cursus.Name, 20)
			level := fmt.Sprintf("Level %s", FormatLevel(cursus.Level))
			if cursus.Grade != "" {
				level += RenderIf(Subtle, fmt.Sprintf(" (%s)", cursus.Grade))
			}
			b.WriteString(ListItem(fmt.Sprintf("%s %s", name, level)))
			b.WriteString("\n")
		}
	}

	// Footer - relative timestamps
	created := RelativeTime(user.CreatedAt)
	updated := RelativeTime(user.UpdatedAt)
	footer := fmt.Sprintf("Created %s · Updated %s", created, updated)
	b.WriteString("\n" + Indent + RenderIf(Subtle, footer) + "\n")

	return b.String()
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
	b.WriteString(Divider(65))
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
		b.WriteString(SectionHeader("CURRENT"))
		b.WriteString("\n")
		for _, loc := range active {
			// Parse begin time for duration calculation
			beginTime, _ := time.Parse(time.RFC3339, loc.BeginAt)
			duration := FormatDuration(time.Since(beginTime))

			b.WriteString(ListItem(fmt.Sprintf("%s %s", Badge(loc.Host, Success), RenderIf(Subtle, "since "+RelativeTime(loc.BeginAt)))))
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf("%s%s%s %s\n", Indent, Indent, RenderIf(Subtle, "Duration:"), duration))
		}
	} else {
		b.WriteString(Indent + Badge("Offline", Muted) + "\n")
	}

	return b.String()
}
