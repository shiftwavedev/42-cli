package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shiftwavedev/42-cli/internal/api"
	"github.com/shiftwavedev/42-cli/internal/helpers"
	"github.com/shiftwavedev/42-cli/pkg/display"
)

var projectsCmd = &cobra.Command{
	Use:   "projects [login]",
	Short: "Display user projects information",
	Long:  "Display user projects information including completed and ongoing projects. If no login is provided, shows your own projects.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		login, err := helpers.GetLoginOrDefault(args, credManager)
		if err != nil {
			log.Fatal("Error:", err)
		}

		webFlag, _ := cmd.Flags().GetBool("web")
		if webFlag {
			currentUserLogin, err := credManager.GetLogin()
			if err != nil {
				log.Fatal("Error getting current user login:", err)
			}
			url := display.IntraURL{}.Projects(login, currentUserLogin)
			fmt.Printf("Opening %s in browser...\n", url)
			if err := display.OpenInBrowser(url); err != nil {
				log.Fatal("Error opening browser:", err)
			}
			return
		}

		spinner := display.NewSimpleSpinner(fmt.Sprintf("Fetching %s's projects...", login))
		spinner.Start()

		projects, err := getUserProjects(login)
		if err != nil {
			spinner.StopWithError("Failed to fetch projects")
			log.Fatal("Error fetching user projects:", err)
		}
		spinner.Stop()

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

	var projects []ProjectUser
	err = api.DefaultClient.Get(fmt.Sprintf("/v2/users/%s/projects_users", login), token, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func displayUserProjects(login string, projects []ProjectUser) {
	// Convert to display types
	var displayProjects []display.ProjectInfo
	for _, p := range projects {
		teamCount := 0
		if len(p.Teams) > 0 {
			teamCount = len(p.Teams[0].Users)
		}
		displayProjects = append(displayProjects, display.ProjectInfo{
			Name:      p.Project.Name,
			Slug:      p.Project.Slug,
			Status:    p.Status,
			Validated: p.Validated,
			FinalMark: p.FinalMark,
			MarkedAt:  p.MarkedAt,
			Retriable: p.Retriable,
			TeamCount: teamCount,
		})
	}

	fmt.Print(display.RenderProjectsList(login, displayProjects))
}

func cloneProject(projectName, targetDir string) error {
	login, err := credManager.GetLogin()
	if err != nil {
		return fmt.Errorf("you need to login first. Use '42-cli auth login'")
	}

	fmt.Println(display.Header("Clone Project", projectName))
	fmt.Println()

	// Find repository with spinner
	spinner := display.NewSimpleSpinner("Finding repository...")
	spinner.Start()

	repoURL, err := getProjectRepoURL(login, projectName)
	if err != nil {
		spinner.StopWithError("Failed to find repository")
		return fmt.Errorf("failed to find repository: %v", err)
	}

	if repoURL == "" {
		spinner.StopWithError("Repository not found")
		return fmt.Errorf("no repository found for project '%s'. Make sure you have access to this project", projectName)
	}
	spinner.StopWithSuccess(fmt.Sprintf("Found: %s", repoURL))

	fmt.Printf("%sTarget: %s\n\n", display.Indent, targetDir)

	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", targetDir)
	}

	// Clone with spinner
	cloneSpinner := display.NewSimpleSpinner("Cloning repository...")
	cloneSpinner.Start()

	cmd := exec.Command("git", "clone", repoURL, targetDir)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		cloneSpinner.StopWithError("Clone failed")
		return fmt.Errorf("git clone failed: %v", err)
	}

	cloneSpinner.StopWithSuccess(fmt.Sprintf("Successfully cloned %s", projectName))
	fmt.Printf("\n%sRun: cd %s && ls -la\n", display.Indent, filepath.Base(targetDir))

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

	var projectUser ProjectUser
	err = api.DefaultClient.Get(fmt.Sprintf("/v2/projects_users/%d", projectUserID), token, &projectUser)
	if err != nil {
		return nil, err
	}

	return projectUser.Teams, nil
}

func init() {
	projectsCmd.AddCommand(cloneCmd)

	projectsCmd.Flags().BoolP("web", "w", false, "Open projects page in browser")
}
