package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CorrectionCard represents a correction session for rendering
type CorrectionCard struct {
	Project      string   // Project name or slug
	Team         string   // Team name or ID
	Date         string   // Pre-formatted via RelativeTime (e.g., "Tomorrow at 14:00")
	Participants []string // Logins (for giver cards) or empty
	Corrector    string   // Login (for receiver cards) or empty
	IsGiver      bool     // true = "Corrections to Give", false = "Evaluations to Receive"
}

// RenderCorrections renders corrections and evaluations as bordered cards
func RenderCorrections(toGive, toReceive []CorrectionCard) string {
	var b strings.Builder

	// Header
	b.WriteString(Header("Corrections & Evaluations", ""))
	b.WriteString("\n")
	b.WriteString(Divider(0))
	b.WriteString("\n\n")

	// Corrections to give
	if len(toGive) > 0 {
		b.WriteString(SectionDivider(fmt.Sprintf("TO GIVE (%d)", len(toGive))))
		b.WriteString("\n\n")
		for _, card := range toGive {
			b.WriteString(renderCorrectionCard(card))
			b.WriteString("\n")
		}
	}

	// Evaluations to receive
	if len(toReceive) > 0 {
		if len(toGive) > 0 {
			b.WriteString("\n")
		}
		b.WriteString(SectionDivider(fmt.Sprintf("TO RECEIVE (%d)", len(toReceive))))
		b.WriteString("\n\n")
		for _, card := range toReceive {
			b.WriteString(renderCorrectionCard(card))
			b.WriteString("\n")
		}
	}

	// Empty state
	if len(toGive) == 0 && len(toReceive) == 0 {
		emptyMsg := Icon("success") + "  No pending corrections or evaluations."
		b.WriteString(Indent + RenderIf(Subtle, emptyMsg) + "\n")
	}

	return b.String()
}

// renderCorrectionCard renders a single correction card with rounded border
func renderCorrectionCard(card CorrectionCard) string {
	// Build card content lines with accessible text labels
	var lines []string

	// Project and team line
	projectTeam := fmt.Sprintf("%s · %s", card.Project, card.Team)
	lines = append(lines, projectTeam)

	// Date line - use text label for accessibility
	lines = append(lines, "When:  "+card.Date)

	// Participants or corrector line - use text labels for accessibility
	if len(card.Participants) > 0 {
		participantsStr := strings.Join(card.Participants, ", ")
		lines = append(lines, "Team:  "+participantsStr)
	} else if card.Corrector != "" {
		lines = append(lines, "By:  "+card.Corrector)
	}

	// Render as compact panel with rounded border
	bodyStr := strings.Join(lines, "\n")
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Padding(0, 2)

	return Indent + strings.ReplaceAll(cardStyle.Render(bodyStr), "\n", "\n"+Indent)
}
