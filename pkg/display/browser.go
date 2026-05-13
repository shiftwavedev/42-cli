package display

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

type IntraURL struct{}

// User returns the profile URL for a given login
func (IntraURL) User(login string) string {
	return fmt.Sprintf("https://profile.intra.42.fr/users/%s", login)
}

// Projects returns the projects URL
// For the current user: https://projects.intra.42.fr/
// For other users: https://projects.intra.42.fr/projects/graph?login=<user>
// NOTE: --web is currently focused on the 42 Central/Paris system and does not yet support other etablissements, which may require a different format :'(
// TODO: maybe ask user if they use version v2 or v3 of intranet?
func (IntraURL) Projects(login, currentUserLogin string) string {
	if login == currentUserLogin {
		return "https://projects.intra.42.fr/"
	}
	return fmt.Sprintf("https://projects.intra.42.fr/projects/graph?login=%s", login)
}
