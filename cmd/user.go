package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io"
	"log"
	"net/http"
	"time"
)

type User struct {
	Login       string `json:"login"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DisplayName string `json:"displayname"`
	Staff       bool   `json:"staff?"`
	CorrectionPoint int `json:"correction_point"`
	PoolMonth   string `json:"pool_month"`
	PoolYear    string `json:"pool_year"`
	Location    string `json:"location"`
	Wallet      int    `json:"wallet"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Alumni      bool   `json:"alumni?"`
	Active      bool   `json:"active?"`
	Campus      []struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"campus"`
	CursusUsers []struct {
		Grade  string  `json:"grade"`
		Level  float64 `json:"level"`
		Cursus struct {
			Name string `json:"name"`
		} `json:"cursus"`
	} `json:"cursus_users"`
}

type Location struct {
	ID       int    `json:"id"`
	BeginAt  string `json:"begin_at"`
	EndAt    *string `json:"end_at"`
	Primary  bool   `json:"primary"`
	Host     string `json:"host"`
	CampusID int    `json:"campus_id"`
	User     struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		URL   string `json:"url"`
	} `json:"user"`
}

var userCmd = &cobra.Command{
	Use:   "user [login]",
	Short: "Display user profile information",
	Long:  "Display user profile information. If no login is provided, shows your own profile.",
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
		
		locations, err := getLocationInfo(login)
		if err != nil {
			log.Fatal("Error fetching location information:", err)
		}
		displayLocationInfo(login, locations)
	},
}

func getUserProfile(login string) (*User, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.intra.42.fr/v2/users/%s", login), nil)
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

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func displayUserProfile(user *User) {
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

	fmt.Printf("\nCreated: %s\n", formatDate(user.CreatedAt))
	fmt.Printf("Updated: %s\n", formatDate(user.UpdatedAt))
}

func formatDate(dateStr string) string {
	if dateStr == "" {
		return "N/A"
	}
	
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	
	return t.Format("2006-01-02 15:04")
}

func formatDateTime(dateStr string) string {
	if dateStr == "" {
		return "N/A"
	}
	
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	
	return t.Format("2006-01-02 15:04")
}

func getLocationInfo(login string) ([]Location, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.intra.42.fr/v2/users/%s/locations", login), nil)
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

	var locations []Location
	if err := json.Unmarshal(body, &locations); err != nil {
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
			fmt.Printf("│ 🕰️  Since: %-34s │\n", formatDateTime(loc.BeginAt))
			fmt.Printf("│ ⏱️  Duration: %-30s │\n", calculateDuration(loc.BeginAt, nil))
			fmt.Println("├──────────────────────────────────────────────────┤")
		}
		fmt.Println("└──────────────────────────────────────────────────┘")
	} else {
		fmt.Println("🚫 Currently offline")
	}
}

func calculateDuration(beginAt string, endAt *string) string {
	beginTime, err := time.Parse(time.RFC3339, beginAt)
	if err != nil {
		return "N/A"
	}

	var endTime time.Time
	if endAt == nil {
		endTime = time.Now()
	} else {
		endTime, err = time.Parse(time.RFC3339, *endAt)
		if err != nil {
			return "N/A"
		}
	}

	duration := endTime.Sub(beginTime)
	return formatDuration(duration)
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	
	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

func init() {
	userCmd.AddCommand(locationCmd)
}