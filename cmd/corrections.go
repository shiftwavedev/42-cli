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

type ScaleTeam struct {
	ID                   int      `json:"id"`
	Scale                Scale    `json:"scale"`
	Comment              string   `json:"comment"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
	Feedback             string   `json:"feedback"`
	FinalMark            *int     `json:"final_mark"`
	Flag                 Flag     `json:"flag"`
	BeginAt              string   `json:"begin_at"`
	Correcteds           any      `json:"correcteds"`
	Correctors           any      `json:"correctors"`
	Truant               struct{} `json:"truant"`
	FilledAt             *string  `json:"filled_at"`
	QuestionsWithAnswers []any    `json:"questions_with_answers"`
	Team                 TeamInfo `json:"team"`
}

type Scale struct {
	ID               int `json:"id"`
	CorrectionNumber int `json:"correction_number"`
	Duration         int `json:"duration"`
}

type Flag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Positive  bool   `json:"positive"`
	Icon      string `json:"icon"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TeamInfo struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	URL           string  `json:"url"`
	FinalMark     *int    `json:"final_mark"`
	ProjectID     int     `json:"project_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Status        string  `json:"status"`
	TerminatingAt *string `json:"terminating_at"`
	Users         []struct {
		ID             int    `json:"id"`
		Login          string `json:"login"`
		URL            string `json:"url"`
		Leader         bool   `json:"leader"`
		Occurrence     int    `json:"occurrence"`
		Validated      bool   `json:"validated"`
		ProjectsUserID int    `json:"projects_user_id"`
	} `json:"users"`
	ProjectSessionID int `json:"project_session_id"`
	Project          struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"project"`
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

		scaleTeams, err := getCorrections(login42)
		if err != nil {
			log.Fatal("Error fetching corrections information:", err)
		}
		displayCorrections(login42, scaleTeams)
	},
}

func getCorrections(login string) ([]ScaleTeam, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.intra.42.fr/v2/users/%s/scale_teams", login), nil)
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
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
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

func displayCorrections(login string, scaleTeams []ScaleTeam) {
	fmt.Printf("=== Corrections: %s ===\n\n", login)

	if len(scaleTeams) == 0 {
		fmt.Println("No corrections or evaluations found.")
		return
	}

	toCorrect := []ScaleTeam{}
	toReceive := []ScaleTeam{}

	for _, st := range scaleTeams {
		if st.FilledAt == nil {
			if isCorrectorForUser(st.Correctors, login) {
				toCorrect = append(toCorrect, st)
			} else {
				toReceive = append(toReceive, st)
			}
		}
	}

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
			fmt.Printf("     👤 Corrector: %s\n", formatCorrectors(st.Correctors))
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

func getStudentNames(users []struct {
	ID             int    `json:"id"`
	Login          string `json:"login"`
	URL            string `json:"url"`
	Leader         bool   `json:"leader"`
	Occurrence     int    `json:"occurrence"`
	Validated      bool   `json:"validated"`
	ProjectsUserID int    `json:"projects_user_id"`
}) string {
	var names []string
	for _, user := range users {
		names = append(names, user.Login)
	}
	return strings.Join(names, ", ")
}

func isCorrectorForUser(correctors any, login string) bool {
	switch v := correctors.(type) {
	case string:
		return strings.Contains(strings.ToLower(v), strings.ToLower(login))
	case []any:
		for _, corrector := range v {
			if correctorStr, ok := corrector.(string); ok {
				if strings.Contains(strings.ToLower(correctorStr), strings.ToLower(login)) {
					return true
				}
			}
		}
	}
	return false
}

func formatCorrectors(correctors any) string {
	switch v := correctors.(type) {
	case string:
		return v
	case []any:
		var names []string
		for _, corrector := range v {
			if correctorStr, ok := corrector.(string); ok {
				names = append(names, correctorStr)
			}
		}
		return strings.Join(names, ", ")
	default:
		return "Unknown"
	}
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
