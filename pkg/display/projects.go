package display

import (
	"fmt"
	"strings"
)

// ProjectInfo represents a project for rendering
type ProjectInfo struct {
	Name      string
	Slug      string
	Status    string
	Validated bool
	FinalMark int
	MarkedAt  string
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

// RenderProjectsList renders a styled projects list
func RenderProjectsList(login string, projects []ProjectInfo) string {
	var b strings.Builder

	// Header with count
	subtitle := fmt.Sprintf("%d total", len(projects))
	b.WriteString(Header(login+"'s projects", subtitle))
	b.WriteString("\n")
	b.WriteString(Divider(65))
	b.WriteString("\n\n")

	if len(projects) == 0 {
		b.WriteString(Indent + RenderIf(Subtle, "No projects found.") + "\n")
		return b.String()
	}

	// Categorize projects (matching original logic from cmd/projects.go)
	var finished, inProgress, failed []ProjectInfo
	for _, p := range projects {
		status := strings.ToLower(p.Status)
		// Finished: status == "finished" AND validated
		if status == "finished" && p.Validated {
			finished = append(finished, p)
		}
		// In Progress: (status == "waiting_for_correction" OR status == "in_progress") AND NOT validated
		if (status == "waiting_for_correction" || status == "in_progress") && !p.Validated {
			inProgress = append(inProgress, p)
		}
		// Failed: status == "failed" only
		if status == "failed" {
			failed = append(failed, p)
		}
	}

	// Finished projects
	if len(finished) > 0 {
		b.WriteString(SectionHeader(fmt.Sprintf("FINISHED (%d)", len(finished))))
		b.WriteString("\n")
		for _, p := range finished {
			b.WriteString(renderProjectLine(p, "finished") + "\n")
		}
	}

	// In progress projects
	if len(inProgress) > 0 {
		if len(finished) > 0 {
			b.WriteString("\n")
		}
		b.WriteString(SectionHeader(fmt.Sprintf("IN PROGRESS (%d)", len(inProgress))))
		b.WriteString("\n")
		for _, p := range inProgress {
			b.WriteString(renderProjectLine(p, "in_progress") + "\n")
		}
	}

	// Failed projects
	if len(failed) > 0 {
		if len(finished) > 0 || len(inProgress) > 0 {
			b.WriteString("\n")
		}
		b.WriteString(SectionHeader(fmt.Sprintf("FAILED (%d)", len(failed))))
		b.WriteString("\n")
		for _, p := range failed {
			b.WriteString(renderProjectLine(p, "failed") + "\n")
		}
	}

	// Summary
	summary := calculateSummary(finished, inProgress, failed)
	b.WriteString("\n")
	b.WriteString(SectionHeader("SUMMARY"))
	b.WriteString("\n")
	summaryLine := fmt.Sprintf("%s %d finished · %s %d in progress · %s %d failed",
		Badge("✓", Success), summary.Finished,
		Badge("◦", Warning), summary.InProgress,
		Badge("✗", Error), summary.Failed,
	)
	b.WriteString(Indent + summaryLine + "\n")

	return b.String()
}

func renderProjectLine(p ProjectInfo, category string) string {
	// Project name (left-aligned, 30 chars)
	name := PadRight(p.Name, 30)
	if len(p.Name) > 30 {
		name = TruncateString(p.Name, 27)
	}

	var statusBadge, score, meta string

	switch category {
	case "finished":
		statusBadge = Badge("Finished", Success)
		if p.FinalMark > 0 {
			score = Badge(fmt.Sprintf("%d", p.FinalMark), Success)
		}
		if p.MarkedAt != "" {
			meta = RenderIf(Subtle, RelativeTime(p.MarkedAt))
		}

	case "in_progress":
		if strings.ToLower(p.Status) == "waiting_for_correction" {
			statusBadge = Badge("Waiting", Warning)
		} else {
			statusBadge = Badge("In Progress", Primary)
		}
		if p.TeamCount > 1 {
			meta = RenderIf(Subtle, fmt.Sprintf("Team: %d members", p.TeamCount))
		}

	case "failed":
		statusBadge = Badge("Failed", Error)
		if p.FinalMark >= 0 {
			score = Badge(fmt.Sprintf("%d", p.FinalMark), Error)
		}
		if p.Retriable {
			meta = Badge("retriable", Warning)
		}
	}

	// Build line
	parts := []string{name, PadRight(statusBadge, 15)}
	if score != "" {
		parts = append(parts, PadRight(score, 5))
	} else {
		parts = append(parts, PadRight("", 5))
	}
	if meta != "" {
		parts = append(parts, meta)
	}

	return ListItem(strings.Join(parts, " "))
}

func calculateSummary(finished, inProgress, failed []ProjectInfo) ProjectSummary {
	return ProjectSummary{
		Finished:   len(finished),
		InProgress: len(inProgress),
		Failed:     len(failed),
	}
}
