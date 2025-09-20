package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var projectsCmd = &cobra.Command{
	Use:   "projects [login]",
	Short: "Display user projects information",
	Long:  "Display user projects information including completed and ongoing projects. If no login is provided, shows your own projects.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var login string
		if len(args) == 0 {
			login42, err := keyring.Get(service, "login42")
			if err != nil {
				log.Fatal("Error: You need to login first. Use '42-cli auth login'")
			}
			login = login42
		} else {
			login = args[0]
		}

		projects, err := getUserProjects(login)
		if err != nil {
			log.Fatal("Error fetching user projects:", err)
		}
		displayUserProjects(login, projects)
	},
}

var cloneCmd = &cobra.Command{
	Use:   "clone [project_name] [directory]",
	Short: "Clone a project repository",
	Long:  "Clone a project repository to the specified directory or current directory if not specified.",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		var targetDir string

		if len(args) == 2 {
			targetDir = args[1]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal("Error getting current directory:", err)
			}
			targetDir = filepath.Join(cwd, projectName)
		}

		err := cloneProject(projectName, targetDir)
		if err != nil {
			log.Fatal("Error cloning project:", err)
		}
	},
}

func getUserProjects(login string) ([]ProjectUser, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.intra.42.fr/v2/users/%s/projects_users", login), nil)
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

	var projects []ProjectUser
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func displayUserProjects(login string, projects []ProjectUser) {
	var finished []ProjectUser
	var inProgress []ProjectUser
	var failed []ProjectUser

	fmt.Printf("=== Projects: %s ===\n", login)
	if len(projects) == 0 {
		fmt.Println("No projects found.")
		return
	}

	for _, project := range projects {
		var status string = strings.ToLower(project.Status)
		if status == "finished" && project.Validated {
			finished = append(finished, project)
		}
		if (status == "waiting_for_correction" || status == "in_progress") && !project.Validated {
			inProgress = append(inProgress, project)
		}
		if status == "failed" {
			failed = append(failed, project)
		}
	}

	if len(finished) > 0 {
		fmt.Printf("\n=== Finished Projects (%d) ===\n", len(finished))
		for _, project := range finished {
			fmt.Printf("✅ %s", project.Project.Name)
			if project.FinalMark > 0 {
				fmt.Printf(" - Score: %d", project.FinalMark)
			}
			fmt.Println()
		}
	}

	if len(inProgress) > 0 {
		fmt.Printf("\n=== In Progress Projects (%d) ===\n", len(inProgress))
		for _, project := range inProgress {
			fmt.Printf("🔄 %s", project.Project.Name)
			if project.MarkedAt != "" {
				fmt.Printf(" (Started: %s)", FormatDate(project.MarkedAt))
			}
			fmt.Println()
		}
	}

	if len(failed) > 0 {
		fmt.Printf("\n=== Failed Projects (%d) ===\n", len(failed))
		for _, project := range failed {
			fmt.Printf("❌ %s", project.Project.Name)
			if project.FinalMark >= 0 {
				fmt.Printf(" - Score: %d", project.FinalMark)
			}
			if project.Retriable {
				fmt.Print(" (Retriable)")
			}
			if project.MarkedAt != "" {
				fmt.Printf(" (Failed: %s)", FormatDate(project.MarkedAt))
			}
			fmt.Println()
		}
	}

	total := len(projects)
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total Projects: %d\n", total)
	fmt.Printf("Finished: %d | In Progress: %d | Failed: %d\n",
		len(finished), len(inProgress), len(failed))
}

func cloneProject(projectName, targetDir string) error {
	login42, err := keyring.Get(service, "login42")
	if err != nil {
		return fmt.Errorf("you need to login first. Use '42-cli auth login'")
	}

	fmt.Printf("=== Cloning Project: %s ===\n\n", projectName)
	fmt.Print("🔍 Finding repository...\n")

	repoURL, err := getProjectRepoURL(login42, projectName)
	if err != nil {
		return fmt.Errorf("failed to find repository: %v", err)
	}

	if repoURL == "" {
		return fmt.Errorf("no repository found for project '%s'. Make sure you have access to this project", projectName)
	}

	fmt.Printf("Found: %s\n", repoURL)
	fmt.Printf("📁 Target directory: %s\n\n", targetDir)

	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", targetDir)
	}

	fmt.Print("🚀 Cloning...\n")

	cmd := exec.Command("git", "clone", repoURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %v", err)
	}

	fmt.Printf("\n✅ Successfully cloned %s\n", projectName)
	fmt.Printf("📂 cd %s && ls -la to get started!\n", filepath.Base(targetDir))

	return nil
}

func getProjectRepoURL(login, projectName string) (string, error) {
	projects, err := getUserProjects(login)
	if err != nil {
		return "", err
	}

	for _, project := range projects {
		if strings.EqualFold(project.Project.Name, projectName) || strings.EqualFold(project.Project.Slug, projectName) {
			teams, err := getProjectTeams(project.ID)
			if err != nil {
				return "", err
			}

			// Find the team with the highest ID (most recent)
			var latestTeam *Team
			for _, team := range teams {
				for _, user := range team.Users {
					if user.Login == login {
						if latestTeam == nil || team.ID > latestTeam.ID {
							latestTeam = &team
						}
					}
				}
			}

			if latestTeam != nil {
				return latestTeam.RepoURL, nil
			}
		}
	}

	return "", fmt.Errorf("project '%s' not found in your projects", projectName)
}

func getProjectTeams(projectUserID int) ([]Team, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.intra.42.fr/v2/projects_users/%d", projectUserID), nil)
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

	var projectUser ProjectUser
	if err := json.Unmarshal(body, &projectUser); err != nil {
		return nil, err
	}

	return projectUser.Teams, nil
}

func init() {
	projectsCmd.AddCommand(cloneCmd)
}
