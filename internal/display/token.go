package display

import (
	"fmt"
	"strconv"
	"time"
)

// TokenExpiryInfo displays token expiry information
func TokenExpiryInfo(expiryStr string) {
	if expiryStr == "" {
		return
	}

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return
	}

	expiryTime := time.Unix(expiry, 0)
	fmt.Printf("expires_at: %s\n", expiryTime.Format("2006-01-02 15:04:05"))

	timeLeft := time.Until(expiryTime)
	if timeLeft > 0 {
		fmt.Printf("time_left: %v\n", timeLeft.Round(time.Minute))
	} else {
		fmt.Printf("status: expired\n")
	}
}
