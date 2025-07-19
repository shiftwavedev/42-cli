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

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User-related commands",
	Long:  "Commands to interact with user data from 42 intranet",
}

var profileCmd = &cobra.Command{
	Use:   "profile [login]",
	Short: "Display user profile information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		login := args[0]
		user, err := getUserProfile(login)
		if err != nil {
			log.Fatal("Error fetching user profile:", err)
		}
		displayUserProfile(user)
	},
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Display your own profile information",
	Run: func(cmd *cobra.Command, args []string) {
		login42, err := keyring.Get(service, "login42")
		if err != nil {
			log.Fatal("Error: You need to login first. Use '42-cli auth login'")
		}
		
		user, err := getUserProfile(login42)
		if err != nil {
			log.Fatal("Error fetching your profile:", err)
		}
		displayUserProfile(user)
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

func init() {
	userCmd.AddCommand(profileCmd)
	userCmd.AddCommand(meCmd)
}