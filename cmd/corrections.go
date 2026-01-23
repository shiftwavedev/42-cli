package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/shiftwavedev/42-cli/internal/api"
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

		toCorrect, err := getCorrectionsAsCorrector()
		if err != nil {
			log.Fatal("Error fetching corrections to give:", err)
		}

		toReceive, err := getCorrectionsAsCorrected()
		if err != nil {
			log.Fatal("Error fetching corrections to receive:", err)
		}
		displayCorrections(toCorrect, toReceive)
	},
}

func getCorrectionsAsCorrector() ([]ScaleTeam, error) {
	token, err := GetAccessToken()
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
	token, err := GetAccessToken()
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
	PrintHeader("Corrections")
	fmt.Print("\n\n")

	if len(toCorrect) > 0 {
		fmt.Printf("🔍 Corrections to Give (%d)\n", len(toCorrect))
		for _, st := range toCorrect {
			projectName := GetProjectDisplayName(st.Team)
			teamName := GetTeamDisplayName(st.Team)

			fmt.Printf("   • %s - %s\n", projectName, teamName)
			fmt.Printf("     📅 %s\n", FormatEvaluationTime(st.BeginAt))
			fmt.Printf("     👥 Correct: %s\n\n", JoinUserLogins(st.Team.Users))
		}
	}

	if len(toReceive) > 0 {
		fmt.Printf("📝 Evaluations to Receive (%d)\n", len(toReceive))
		for _, st := range toReceive {
			projectName := GetProjectDisplayName(st.Team)
			teamName := GetTeamDisplayName(st.Team)

			fmt.Printf("   • %s - %s\n", projectName, teamName)
			fmt.Printf("     📅 %s\n", FormatEvaluationTime(st.BeginAt))
			fmt.Printf("     👤 Corrector: %s\n\n", st.Corrector.Login)
		}
	}

	if len(toCorrect) == 0 && len(toReceive) == 0 {
		fmt.Println("✅ No pending corrections or evaluations.")
	}
}
