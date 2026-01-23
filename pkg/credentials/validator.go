package credentials

import "regexp"

// Validator handles credential validation
type Validator struct {
	login42Pattern      *regexp.Regexp
	clientIDPattern     *regexp.Regexp
	clientSecretPattern *regexp.Regexp
}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	return &Validator{
		login42Pattern:      regexp.MustCompile("^[a-z-]+$"),
		clientIDPattern:     regexp.MustCompile("^u-s4t2ud-[a-z0-9]+$"),
		clientSecretPattern: regexp.MustCompile("^s-s4t2ud-[a-z0-9]+$"),
	}
}

// Validate checks if credentials are in the correct format
func (v *Validator) Validate(creds *Credentials) bool {
	return v.login42Pattern.MatchString(creds.Login42) &&
		v.clientIDPattern.MatchString(creds.ClientID) &&
		v.clientSecretPattern.MatchString(creds.ClientSecret)
}

// ValidateLogin42 checks if login42 is in the correct format
func (v *Validator) ValidateLogin42(login string) bool {
	return v.login42Pattern.MatchString(login)
}

// ValidateClientID checks if client ID is in the correct format
func (v *Validator) ValidateClientID(clientID string) bool {
	return v.clientIDPattern.MatchString(clientID)
}

// ValidateClientSecret checks if client secret is in the correct format
func (v *Validator) ValidateClientSecret(secret string) bool {
	return v.clientSecretPattern.MatchString(secret)
}
