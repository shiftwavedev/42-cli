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
	b.WriteString(Divider(65))
	b.WriteString("\n\n")

	// Credentials section
	if creds != nil {
		b.WriteString(SectionHeader("CREDENTIALS"))
		b.WriteString("\n")
		b.WriteString(KeyValue("Login", creds.Login, 12))
		b.WriteString("\n")
		b.WriteString(KeyValue("Client ID", creds.ClientID, 12))
		b.WriteString("\n")
		b.WriteString(KeyValue("Secret", MaskToken(creds.ClientSecret), 12))
		b.WriteString("\n")
	}

	// User token section
	if userToken != nil {
		b.WriteString("\n")
		b.WriteString(SectionHeader("USER TOKEN"))
		b.WriteString("\n")
		b.WriteString(renderTokenSection(userToken, true))
	} else {
		b.WriteString("\n")
		b.WriteString(SectionHeader("USER TOKEN"))
		b.WriteString("\n")
		b.WriteString(KeyValue("Status", Badge("Not authenticated", Muted), 12))
		b.WriteString("\n")
		b.WriteString(Indent + RenderIf(Subtle, "Run '42-cli auth login' to authenticate"))
		b.WriteString("\n")
	}

	// App token section
	if appToken != nil {
		b.WriteString("\n")
		b.WriteString(SectionHeader("APP TOKEN"))
		b.WriteString("\n")
		b.WriteString(renderTokenSection(appToken, false))
	}

	// Footer hint
	b.WriteString("\n")
	b.WriteString(Indent + RenderIf(Subtle, "Use '42-cli auth refresh' to renew tokens"))
	b.WriteString("\n")

	return b.String()
}

func renderTokenSection(token *TokenInfo, includeRefresh bool) string {
	var b strings.Builder

	// Status with expiry info
	status := getTokenStatus(token.ExpiryUnix)
	b.WriteString(KeyValue("Status", status.badge+" "+status.detail, 12))
	b.WriteString("\n")

	// Token (masked)
	b.WriteString(KeyValue("Token", RenderIf(Code, MaskToken(token.Token)), 12))
	b.WriteString("\n")

	// Refresh token (if applicable)
	if includeRefresh && token.RefreshToken != "" {
		b.WriteString(KeyValue("Refresh", RenderIf(Code, MaskToken(token.RefreshToken)), 12))
		b.WriteString("\n")
	}

	return b.String()
}

type tokenStatus struct {
	badge  string
	detail string
}

func getTokenStatus(expiryUnix int64) tokenStatus {
	if expiryUnix == 0 {
		return tokenStatus{
			badge:  Badge("Unknown", Muted),
			detail: "",
		}
	}

	expiryTime := time.Unix(expiryUnix, 0)
	remaining := time.Until(expiryTime)

	if remaining < 0 {
		return tokenStatus{
			badge:  Badge("Expired", Error),
			detail: RenderIf(Subtle, fmt.Sprintf("(%s)", RelativeTimeFromTime(expiryTime))),
		}
	}

	// Less than 1 hour - warning
	if remaining < time.Hour {
		return tokenStatus{
			badge:  Badge("Expiring", Warning),
			detail: RenderIf(Subtle, fmt.Sprintf("(in %s)", FormatDuration(remaining))),
		}
	}

	// Less than 24 hours - still valid but show time
	if remaining < 24*time.Hour {
		return tokenStatus{
			badge:  Badge("Valid", Success),
			detail: RenderIf(Subtle, fmt.Sprintf("(expires in %s)", FormatDuration(remaining))),
		}
	}

	// More than 24 hours
	return tokenStatus{
		badge:  Badge("Valid", Success),
		detail: RenderIf(Subtle, fmt.Sprintf("(expires in %s)", FormatDuration(remaining))),
	}
}

// RenderTokenExpiry renders token expiry information (simple version for existing code)
func RenderTokenExpiry(expiryUnix int64) string {
	status := getTokenStatus(expiryUnix)
	return status.badge + " " + status.detail
}
