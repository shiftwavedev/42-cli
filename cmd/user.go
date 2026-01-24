package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/shiftwavedev/42-cli/internal/api"
	"github.com/shiftwavedev/42-cli/internal/helpers"
	"github.com/shiftwavedev/42-cli/pkg/display"
)

var userCmd = &cobra.Command{
	Use:   "user [login]",
	Short: "Display user profile information",
	Long:  "Display user profile information. If no login is provided, shows your own profile.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		login, err := helpers.GetLoginOrDefault(args, credManager)
		if err != nil {
			log.Fatal("Error:", err)
		}

		webFlag, _ := cmd.Flags().GetBool("web")
		if webFlag {
			url := display.IntraURL{}.User(login)
			fmt.Printf("Opening %s in browser...\n", url)
			if err := display.OpenInBrowser(url); err != nil {
				log.Fatal("Error opening browser:", err)
			}
			return
		}

		spinner := display.NewSimpleSpinner(fmt.Sprintf("Fetching %s's profile...", login))
		spinner.Start()

		user, err := getUserProfile(login)
		if err != nil {
			spinner.StopWithError("Failed to fetch profile")
			log.Fatal("Error fetching user profile:", err)
		}
		spinner.Stop()

		displayUserProfile(user)
	},
}

var locationCmd = &cobra.Command{
	Use:   "location [login]",
	Short: "Display user location information",
	Long:  "Display current and recent location information for a user. If no login is provided, shows your own location.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		login, err := helpers.GetLoginOrDefault(args, credManager)
		if err != nil {
			log.Fatal("Error:", err)
		}

		spinner := display.NewSimpleSpinner(fmt.Sprintf("Fetching %s's location...", login))
		spinner.Start()

		locations, err := getLocationInfo(login)
		if err != nil {
			spinner.StopWithError("Failed to fetch location")
			log.Fatal("Error fetching location information:", err)
		}
		spinner.Stop()

		displayLocationInfo(login, locations)
	},
}

func getUserProfile(login string) (*UserProfile, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	var user UserProfile
	err = api.DefaultClient.Get(fmt.Sprintf("/v2/users/%s", login), token, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func displayUserProfile(user *UserProfile) {
	// Convert to display types
	profile := &display.UserProfile{
		Login:           user.Login,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Phone:           user.Phone,
		DisplayName:     user.DisplayName,
		Staff:           user.Staff,
		CorrectionPoint: user.CorrectionPoint,
		PoolMonth:       user.PoolMonth,
		PoolYear:        user.PoolYear,
		Location:        user.Location,
		Wallet:          user.Wallet,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Alumni:          user.Alumni,
		Active:          user.Active,
	}

	// Convert campus
	for _, c := range user.Campus {
		profile.Campus = append(profile.Campus, display.CampusInfo{
			Name:    c.Name,
			Country: c.Country,
		})
	}

	// Convert cursus
	for _, c := range user.CursusUsers {
		profile.CursusUsers = append(profile.CursusUsers, display.CursusUser{
			Grade: c.Grade,
			Level: c.Level,
			Cursus: display.CursusInfo{
				Name: c.Cursus.Name,
			},
		})
	}

	fmt.Print(display.RenderUserProfile(profile))
}

func getLocationInfo(login string) ([]Location, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	var locations []Location
	err = api.DefaultClient.Get(fmt.Sprintf("/v2/users/%s/locations", login), token, &locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func displayLocationInfo(login string, locations []Location) {
	// Convert to display types
	var displayLocs []display.LocationInfo
	for _, loc := range locations {
		displayLocs = append(displayLocs, display.LocationInfo{
			Host:    loc.Host,
			BeginAt: loc.BeginAt,
			EndAt:   loc.EndAt,
			Primary: loc.Primary,
		})
	}

	fmt.Print(display.RenderLocationInfo(login, displayLocs))
}

func init() {
	userCmd.AddCommand(locationCmd)

	userCmd.Flags().BoolP("web", "w", false, "Open user profile in browser")
}
