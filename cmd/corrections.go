package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/shiftwavedev/42-cli/internal/api"
	"github.com/shiftwavedev/42-cli/pkg/auth"
	"github.com/shiftwavedev/42-cli/pkg/display"
)

var correctionsCmd = &cobra.Command{
	Use:   "corrections",
	Short: "Display your corrections to give and evaluations to receive",
	Long:  "Display your upcoming corrections to give and evaluations to receive.",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := credManager.GetLogin()
		if err != nil {
			log.Fatal("Error: You need to login first. Use '42-cli auth login'")
		}

		spinner := display.NewSpinner("Fetching corrections...")

		toCorrect, err := getCorrectionsAsCorrector()
		if err != nil {
			spinner.Fail("Failed to fetch corrections")
			log.Fatal("Error fetching corrections to give:", err)
		}

		toReceive, err := getCorrectionsAsCorrected()
		if err != nil {
			spinner.Fail("Failed to fetch corrections")
			log.Fatal("Error fetching corrections to receive:", err)
		}
		spinner.Stop()

		displayCorrections(toCorrect, toReceive)
	},
}

func getCorrectionsAsCorrector() ([]ScaleTeam, error) {
	token, err := tokenManager.GetValidToken(auth.ScopeUser)
	if err != nil {
		return nil, err
	}

	var scaleTeams []ScaleTeam
	err = api.DefaultClient.Get("/v2/me/scale_teams/as_corrector", token, &scaleTeams)
	if err != nil {
		return nil, err
	}

	return scaleTeams, nil
}

func getCorrectionsAsCorrected() ([]ScaleTeam, error) {
	token, err := tokenManager.GetValidToken(auth.ScopeUser)
	if err != nil {
		return nil, err
	}

	var scaleTeams []ScaleTeam
	err = api.DefaultClient.Get("/v2/me/scale_teams/as_corrected", token, &scaleTeams)
	if err != nil {
		return nil, err
	}

	return scaleTeams, nil
}

func displayCorrections(toCorrect []ScaleTeam, toReceive []ScaleTeam) {
	// Convert to display.CorrectionCard for rendering
	var toGiveCards []display.CorrectionCard
	for _, st := range toCorrect {
		projectName := st.Team.Project.Name
		if projectName == "" {
			projectName = st.Team.Project.Slug
		}
		if projectName == "" {
			projectName = "Unknown"
		}

		teamName := st.Team.Name
		if teamName == "" {
			teamName = fmt.Sprintf("team-%d", st.Team.ID)
		}

		// Map users to logins
		var participants []string
		for _, user := range st.Team.Users {
			participants = append(participants, user.Login)
		}

		toGiveCards = append(toGiveCards, display.CorrectionCard{
			Project:      projectName,
			Team:         teamName,
			Date:         display.RelativeTime(st.BeginAt),
			Participants: participants,
			IsGiver:      true,
		})
	}

	var toReceiveCards []display.CorrectionCard
	for _, st := range toReceive {
		projectName := st.Team.Project.Name
		if projectName == "" {
			projectName = st.Team.Project.Slug
		}
		if projectName == "" {
			projectName = "Unknown"
		}

		teamName := st.Team.Name
		if teamName == "" {
			teamName = fmt.Sprintf("team-%d", st.Team.ID)
		}

		toReceiveCards = append(toReceiveCards, display.CorrectionCard{
			Project:   projectName,
			Team:      teamName,
			Date:      display.RelativeTime(st.BeginAt),
			Corrector: st.Corrector.Login,
			IsGiver:   false,
		})
	}

	fmt.Print(display.RenderCorrections(toGiveCards, toReceiveCards))
}
