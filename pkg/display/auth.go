package display

import (
	"fmt"
	"strings"
	"time"
)

// TokenInfo represents token information for rendering
type TokenInfo struct {
	Token        string
	ExpiryUnix   int64
	RefreshToken string
}

// CredentialInfo represents stored credentials for rendering
type CredentialInfo struct {
	Login        string
	ClientID     string
	ClientSecret string
}

// RenderAuthStatus renders a complete authentication status view
func RenderAuthStatus(creds *CredentialInfo, appToken, userToken *TokenInfo) string {
	var b strings.Builder

	b.WriteString(Header("Authentication Status", ""))
	b.WriteString("\n")
	b.WriteString(Divider(0))
	b.WriteString("\n\n")

	// Credentials section
	if creds != nil {
		credentialsBody := fmt.Sprintf(
			"%s  %s\n%s  %s\n%s  %s",
			RenderIf(Subtle, "Login"),
			creds.Login,
			RenderIf(Subtle, "Client ID"),
			creds.ClientID,
			RenderIf(Subtle, "Secret"),
			MaskToken(creds.ClientSecret),
		)
		b.WriteString(Panel("Credentials", credentialsBody))
		b.WriteString("\n\n")
	} else {
		noCredsBody := "No credentials found.\n" +
			"Run '42-cli auth login' to get started."
		b.WriteString(Panel("Credentials", noCredsBody))
		b.WriteString("\n\n")
	}

	// User token section
	if userToken != nil {
		userTokenBody := renderTokenSectionBody(userToken, true)
		b.WriteString(Panel("User Token", userTokenBody))
	} else {
		noAuthBody := fmt.Sprintf("%s Not authenticated\nRun '42-cli auth login' to authenticate",
			Icon("offline"))
		b.WriteString(Panel("User Token", noAuthBody))
	}

	// App token section
	if appToken != nil {
		b.WriteString("\n\n")
		appTokenBody := renderTokenSectionBody(appToken, false)
		b.WriteString(Panel("App Token", appTokenBody))
	}

	// Footer hint
	b.WriteString("\n\n")
	b.WriteString(Indent + RenderIf(Subtle, "Run '42-cli auth refresh' to renew tokens"))
	b.WriteString("\n")

	return b.String()
}

// renderTokenSectionBody builds the body content for a token panel
func renderTokenSectionBody(token *TokenInfo, includeRefresh bool) string {
	var b strings.Builder

	// Status with expiry info
	status := getTokenStatus(token.ExpiryUnix)
	statusLine := fmt.Sprintf("%s  %s\n",
		RenderIf(Subtle, "Status"),
		status.full,
	)
	b.WriteString(statusLine)

	// Token (masked)
	tokenLine := fmt.Sprintf("%s  %s\n",
		RenderIf(Subtle, "Token"),
		RenderIf(Code, MaskToken(token.Token)),
	)
	b.WriteString(tokenLine)

	// Refresh token (if applicable)
	if includeRefresh && token.RefreshToken != "" {
		refreshLine := fmt.Sprintf("%s  %s",
			RenderIf(Subtle, "Refresh"),
			RenderIf(Code, MaskToken(token.RefreshToken)),
		)
		b.WriteString(refreshLine)
	}

	return b.String()
}

// tokenStatus holds the rendered status representation
type tokenStatus struct {
	icon   string // The icon to display
	label  string // The status label
	detail string // Additional detail (e.g., time info)
	full   string // Combined: icon label detail
}

// getTokenStatus returns the token status with proper icons and labels
func getTokenStatus(expiryUnix int64) tokenStatus {
	if expiryUnix == 0 {
		return tokenStatus{
			icon:   Icon("offline"),
			label:  "Unknown",
			detail: "",
			full:   Icon("offline") + " Unknown",
		}
	}

	expiryTime := time.Unix(expiryUnix, 0)
	remaining := time.Until(expiryTime)

	if remaining < 0 {
		// Expired
		relTime := RelativeTimeFromTime(expiryTime)
		return tokenStatus{
			icon:   Icon("error"),
			label:  "Expired",
			detail: fmt.Sprintf("(%s)", relTime),
			full:   Icon("error") + " Expired " + RenderIf(Subtle, fmt.Sprintf("(%s)", relTime)),
		}
	}

	// Less than 1 hour - warning
	if remaining < time.Hour {
		duration := FormatDuration(remaining)
		return tokenStatus{
			icon:   Icon("warning"),
			label:  "Expiring",
			detail: fmt.Sprintf("(in %s)", duration),
			full:   Icon("warning") + " Expiring " + RenderIf(Subtle, fmt.Sprintf("(in %s)", duration)),
		}
	}

	// Valid token
	duration := FormatDuration(remaining)
	return tokenStatus{
		icon:   Icon("online"),
		label:  "Valid",
		detail: fmt.Sprintf("(expires in %s)", duration),
		full:   Icon("online") + " Valid " + RenderIf(Subtle, fmt.Sprintf("(expires in %s)", duration)),
	}
}

// RenderTokenExpiry renders token expiry information (simple version for existing code)
func RenderTokenExpiry(expiryUnix int64) string {
	status := getTokenStatus(expiryUnix)
	return status.full
}
