package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

type CorrectionUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	URL   string `json:"url"`
}

type ScaleTeam struct {
	ID                   int              `json:"id"`
	ScaleID              int              `json:"scale_id"`
	Comment              *string          `json:"comment"`
	CreatedAt            string           `json:"created_at"`
	UpdatedAt            string           `json:"updated_at"`
	Feedback             *string          `json:"feedback"`
	FinalMark            *int             `json:"final_mark"`
	Flag                 Flag             `json:"flag"`
	BeginAt              string           `json:"begin_at"`
	Correcteds           []CorrectionUser `json:"correcteds"`
	Corrector            CorrectionUser   `json:"corrector"`
	Truant               struct{}         `json:"truant"`
	FilledAt             *string          `json:"filled_at"`
	QuestionsWithAnswers []any            `json:"questions_with_answers"`
	Scale                Scale            `json:"scale"`
	Team                 TeamInfo         `json:"team"`
	Feedbacks            []any            `json:"feedbacks"`
}

type Language struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type Scale struct {
	ID                 int        `json:"id"`
	EvaluationID       int        `json:"evaluation_id"`
	Name               string     `json:"name"`
	IsPrimary          bool       `json:"is_primary"`
	Comment            string     `json:"comment"`
	IntroductionMd     string     `json:"introduction_md"`
	DisclaimerMd       string     `json:"disclaimer_md"`
	GuidelinesMd       string     `json:"guidelines_md"`
	CreatedAt          string     `json:"created_at"`
	CorrectionNumber   int        `json:"correction_number"`
	Duration           int        `json:"duration"`
	ManualSubscription bool       `json:"manual_subscription"`
	Languages          []Language `json:"languages"`
	Flags              []Flag     `json:"flags"`
	Free               bool       `json:"free"`
}

type Flag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Positive  bool   `json:"positive"`
	Icon      string `json:"icon"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TeamUser struct {
	ID             int    `json:"id"`
	Login          string `json:"login"`
	URL            string `json:"url"`
	Leader         bool   `json:"leader"`
	Occurrence     int    `json:"occurrence"`
	Validated      bool   `json:"validated"`
	ProjectsUserID int    `json:"projects_user_id"`
}

type CorrectionProject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type TeamInfo struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	URL               string            `json:"url"`
	FinalMark         *int              `json:"final_mark"`
	ProjectID         int               `json:"project_id"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
	Status            string            `json:"status"`
	TerminatingAt     *string           `json:"terminating_at"`
	Users             []TeamUser        `json:"users"`
	Locked            bool              `json:"locked?"`
	Validated         *bool             `json:"validated?"`
	Closed            bool              `json:"closed?"`
	RepoURL           string            `json:"repo_url"`
	RepoUUID          string            `json:"repo_uuid"`
	LockedAt          string            `json:"locked_at"`
	ClosedAt          string            `json:"closed_at"`
	ProjectSessionID  int               `json:"project_session_id"`
	ProjectGitlabPath string            `json:"project_gitlab_path"`
	Project           CorrectionProject `json:"project"`
}

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

func displayCorrections(login string, toCorrect []ScaleTeam, toReceive []ScaleTeam) {
	fmt.Print("=== Corrections ===\n\n")

	if len(toCorrect) > 0 {
		fmt.Printf("🔍 Corrections to Give (%d)\n", len(toCorrect))
		for _, st := range toCorrect {
			projectName := getProjectName(st)
			teamName := st.Team.Name
			if teamName == "" {
				teamName = fmt.Sprintf("team-%d", st.Team.ID)
			}

			fmt.Printf("   • %s - %s\n", projectName, teamName)
			fmt.Printf("     📅 %s\n", formatEvaluationTime(st.BeginAt))
			fmt.Printf("     👥 Correct: %s\n", getStudentNames(st.Team.Users))
			fmt.Println()
		}
	}

	if len(toReceive) > 0 {
		fmt.Printf("📝 Evaluations to Receive (%d)\n", len(toReceive))
		for _, st := range toReceive {
			projectName := getProjectName(st)
			teamName := st.Team.Name
			if teamName == "" {
				teamName = fmt.Sprintf("team-%d", st.Team.ID)
			}

			fmt.Printf("   • %s - %s\n", projectName, teamName)
			fmt.Printf("     📅 %s\n", formatEvaluationTime(st.BeginAt))
			fmt.Printf("     👤 Corrector: %s\n", st.Corrector.Login)
			fmt.Println()
		}
	}

	if len(toCorrect) == 0 && len(toReceive) == 0 {
		fmt.Println("✅ No pending corrections or evaluations.")
	}
}

func getProjectName(st ScaleTeam) string {
	if st.Team.Project.Name != "" {
		return st.Team.Project.Name
	}
	if st.Team.Project.Slug != "" {
		return st.Team.Project.Slug
	}
	return "Unknown Project"
}

func getStudentNames(users []TeamUser) string {
	var names []string
	for _, user := range users {
		names = append(names, user.Login)
	}
	return strings.Join(names, ", ")
}

func formatEvaluationTime(dateStr string) string {
	if dateStr == "" {
		return "N/A"
	}

	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}

	return t.Format("2006-01-02 15:04")
}
