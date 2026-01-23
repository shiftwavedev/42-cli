package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/shiftwavedev/42-cli/internal/api"
	"github.com/shiftwavedev/42-cli/internal/helpers"
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

		user, err := getUserProfile(login)
		if err != nil {
			log.Fatal("Error fetching user profile:", err)
		}
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

		locations, err := getLocationInfo(login)
		if err != nil {
			log.Fatal("Error fetching location information:", err)
		}
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
	fmt.Printf("=== Profile: %s ===\n", user.Login)
	fmt.Printf("Name: %s %s\n", user.FirstName, user.LastName)
	fmt.Printf("Display Name: %s\n", user.DisplayName)
	fmt.Printf("Email: %s\n", user.Email)

	if user.Phone != "" {
		fmt.Printf("Phone: %s\n", user.Phone)
	}

	fmt.Printf("Staff: %t\n", user.Staff)
	fmt.Printf("Alumni: %t\n", user.Alumni)
	fmt.Printf("Active: %t\n", user.Active)
	fmt.Printf("Correction Points: %d\n", user.CorrectionPoint)
	fmt.Printf("Wallet: %d\n", user.Wallet)

	if user.Location != "" {
		fmt.Printf("Location: %s\n", user.Location)
	}

	if user.PoolMonth != "" && user.PoolYear != "" {
		fmt.Printf("Pool: %s %s\n", user.PoolMonth, user.PoolYear)
	}

	if len(user.Campus) > 0 {
		fmt.Printf("\n=== Campus ===\n")
		for _, campus := range user.Campus {
			fmt.Printf("- %s (%s)\n", campus.Name, campus.Country)
		}
	}

	if len(user.CursusUsers) > 0 {
		fmt.Printf("\n=== Cursus ===\n")
		for _, cursus := range user.CursusUsers {
			fmt.Printf("- %s: Level %.2f", cursus.Cursus.Name, cursus.Level)
			if cursus.Grade != "" {
				fmt.Printf(" (Grade: %s)", cursus.Grade)
			}
			fmt.Println()
		}
	}

	PrintKeyValue("Created", FormatDate(user.CreatedAt))
	PrintKeyValue("Updated", FormatDate(user.UpdatedAt))
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
	active := []Location{}

	fmt.Printf("=== Location Information: %s ===\n\n", login)
	if len(locations) == 0 {
		fmt.Println("No location information found.")
		return
	}

	for _, loc := range locations {
		if loc.EndAt == nil {
			active = append(active, loc)
		}
	}

	if len(active) > 0 {
		fmt.Printf("📍 Current Location\n")
		fmt.Println("┌──────────────────────────────────────────────────┐")
		for _, loc := range active {
			fmt.Printf("│ 🖥️  Host: %-35s │\n", loc.Host)
			fmt.Printf("│ 🕰️  Since: %-34s │\n", FormatDateTime(loc.BeginAt))
			fmt.Printf("│ ⏱️  Duration: %-30s │\n", CalculateDuration(loc.BeginAt, nil))
			fmt.Println("├──────────────────────────────────────────────────┤")
		}
		fmt.Println("└──────────────────────────────────────────────────┘")
	} else {
		fmt.Println("🚫 Currently offline")
	}
}

func init() {
	userCmd.AddCommand(locationCmd)
}
