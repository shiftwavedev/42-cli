package display

import (
	"fmt"
	"time"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// ProjectInfo represents a project for rendering
type ProjectInfo struct {
	Name      string
	Slug      string
	Status    string
	Validated bool
	FinalMark int
	MarkedAt  string
	CursusIds int
	Retriable bool
	TeamCount int
}

// ProjectSummary holds counts for project categories
type ProjectSummary struct {
	Finished   int
	InProgress int
	Failed     int
	Waiting    int
}

// RenderProjectsList renders a styled projects list using a table
func RenderProjectsList(login string, projects []ProjectInfo) string {
	var b strings.Builder

	// Header with count
	subtitle := fmt.Sprintf("%d total", len(projects))
	b.WriteString(Header(login+"'s projects", subtitle))
	b.WriteString("\n")
	b.WriteString(Divider(0))
	b.WriteString("\n\n")

	if len(projects) == 0 {
		b.WriteString(Indent + RenderIf(Subtle, "No projects found.") + "\n")
		return b.String()
	}

	var finished, inProgress, failed []ProjectInfo
	for _, p := range projects {
		// NOTE: hardcoded to cursus_id 21 (common core)
		if p.CursusIds == 21 {
			if (p.FinalMark >= 80 && p.Validated == true) || (p.FinalMark >= 80 && p.Validated == true && p.Status == "finished") {
				finished = append(finished, p)
			}

			if p.FinalMark < 80 && !p.Validated && p.Status == "finished" {
				failed = append(failed, p)
			}

			if p.FinalMark == 0 && !p.Validated && (p.Status == "in_progress" || p.Status == "waiting_for_correction") {
				inProgress = append(inProgress, p)
			}
		}

	}

	sortProjectsByMarkedAt(finished)
	sortProjectsByMarkedAt(inProgress)
	sortProjectsByMarkedAt(failed)

	// Finished projects
	if len(finished) > 0 {
		finishedTable := renderProjectsTable(finished, "finished")
		b.WriteString(DividerForContent(fmt.Sprintf("FINISHED (%d)", len(finished)), finishedTable))
		b.WriteString("\n\n")
		b.WriteString(finishedTable)
		b.WriteString("\n")
	}

	// In progress projects
	if len(inProgress) > 0 {
		if len(finished) > 0 {
			b.WriteString("\n")
		}
		inProgressTable := renderProjectsTable(inProgress, "in_progress")
		b.WriteString(DividerForContent(fmt.Sprintf("IN PROGRESS (%d)", len(inProgress)), inProgressTable))
		b.WriteString("\n\n")
		b.WriteString(inProgressTable)
		b.WriteString("\n")
	}

	// Failed projects
	if len(failed) > 0 {
		if len(finished) > 0 || len(inProgress) > 0 {
			b.WriteString("\n")
		}
		failedTable := renderProjectsTable(failed, "failed")
		b.WriteString(DividerForContent(fmt.Sprintf("FAILED (%d)", len(failed)), failedTable))
		b.WriteString("\n\n")
		b.WriteString(failedTable)
		b.WriteString("\n")
	}

	// Summary
	summary := calculateSummary(finished, inProgress, failed)
	b.WriteString("\n")
	summaryLine := fmt.Sprintf("%s %d finished · %s %d in progress · %s %d failed",
		Badge("✓", Success), summary.Finished,
		Badge("◦", Warning), summary.InProgress,
		Badge("✗", Error), summary.Failed,
	)
	b.WriteString(DividerForContent("SUMMARY", Indent+summaryLine))
	b.WriteString("\n\n")
	b.WriteString(Indent + summaryLine + "\n")

	return b.String()
}

// renderProjectsTable builds a lipgloss table for a list of projects
func renderProjectsTable(projects []ProjectInfo, category string) string {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(Border)).
		StyleFunc(func(row, col int) lipgloss.Style {
			// Header row styling
			if row == 0 {
				return lipgloss.NewStyle().
					Bold(true).
					Foreground(Primary).
					Padding(0, 1)
			}
			// Regular rows
			return lipgloss.NewStyle().Padding(0, 1)
		}).
		Headers("PROJECT", "STATUS", "SCORE", "DATE").
		Rows()

	for _, p := range projects {
		row := buildProjectRow(p, category)
		t = t.Row(row...)
	}

	return Indent + strings.ReplaceAll(t.Render(), "\n", "\n"+Indent)
}

// buildProjectRow constructs a row ([]string) for a project in the table
func buildProjectRow(p ProjectInfo, category string) []string {
	// Column 1: Project name (truncate if too long)
	name := p.Name
	maxLen := CalculateMaxProjectNameLength()
	if len(name) > maxLen {
		name = TruncateString(name, maxLen-3)
	}

	// Column 2: Status badge
	var statusBadge string
	switch category {
	case "finished":
		statusBadge = Badge("Finished", Success)

	case "in_progress":
		if strings.ToLower(p.Status) == "waiting_for_correction" {
			statusBadge = Badge("Waiting", Warning)
		} else {
			statusBadge = Badge("In Progress", Primary)
		}

	case "failed":
		statusBadge = Badge("Failed", Error)
	}

	// Column 3: Score
	var score string
	switch category {
	case "finished":
		if p.FinalMark > 0 {
			score = Badge(fmt.Sprintf("%d", p.FinalMark), Success)
		} else {
			score = "—"
		}

	case "in_progress":
		score = "—"

	case "failed":
		if p.FinalMark >= 0 {
			score = Badge(fmt.Sprintf("%d", p.FinalMark), Error)
		} else {
			score = "—"
		}
	}

	// Column 4: Date (marked_at)
	notes := "—"
	if p.MarkedAt != "" {
		notes = RenderIf(Subtle, formatMarkedAt(p.MarkedAt))
	}

	return []string{name, statusBadge, score, notes}
}

func calculateSummary(finished, inProgress, failed []ProjectInfo) ProjectSummary {
	return ProjectSummary{
		Finished:   len(finished),
		InProgress: len(inProgress),
		Failed:     len(failed),
	}
}

func sortProjectsByMarkedAt(projects []ProjectInfo) {
	sort.SliceStable(projects, func(i, j int) bool {
		first := parseMarkedAt(projects[i].MarkedAt)
		second := parseMarkedAt(projects[j].MarkedAt)

		if first.IsZero() && second.IsZero() {
			return false
		}
		if first.IsZero() {
			return false
		}
		if second.IsZero() {
			return true
		}

		return first.After(second)
	})
}

func parseMarkedAt(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}
	}

	return parsed
}

func formatMarkedAt(value string) string {
	parsed := parseMarkedAt(value)
	if parsed.IsZero() {
		return value
	}

	return RelativeTimeFromTime(parsed)
}
