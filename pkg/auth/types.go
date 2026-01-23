package auth

// TokenScope represents the scope of a token
type TokenScope int

const (
	// ScopeApplication represents client credentials token (public API access)
	ScopeApplication TokenScope = iota
	// ScopeUser represents OAuth user token (authenticated endpoints)
	ScopeUser
)

// String returns the string representation of the token scope
func (s TokenScope) String() string {
	switch s {
	case ScopeApplication:
		return "application"
	case ScopeUser:
		return "user"
	default:
		return "unknown"
	}
}

// Token represents a stored token with its metadata
type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}
