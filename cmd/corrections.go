package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var correctionsCmd = &cobra.Command{
	Use:   "corrections",
	Short: "Display your corrections to give and evaluations to receive",
	Long:  "Display your upcoming corrections to give and evaluations to receive.",
	Run: func(cmd *cobra.Command, args []string) {
		login42, err := keyring.Get(service, "login42")
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

		displayCorrections(login42, toCorrect, toReceive)
	},
}

func getCorrectionsAsCorrector() ([]ScaleTeam, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", "https://api.intra.42.fr/v2/me/scale_teams/as_corrector", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var scaleTeams []ScaleTeam
	if err := json.Unmarshal(body, &scaleTeams); err != nil {
		return nil, err
	}

	return scaleTeams, nil
}

func getCorrectionsAsCorrected() ([]ScaleTeam, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", "https://api.intra.42.fr/v2/me/scale_teams/as_corrected", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var scaleTeams []ScaleTeam
	if err := json.Unmarshal(body, &scaleTeams); err != nil {
		return nil, err
	}

	return scaleTeams, nil
}

func displayCorrections(_ string, toCorrect []ScaleTeam, toReceive []ScaleTeam) {
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
