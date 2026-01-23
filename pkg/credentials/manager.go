package credentials

import "fmt"

// Manager handles credential operations
type Manager struct {
	storage   Storage
	validator *Validator
}

// NewManager creates a new credentials Manager
func NewManager(storage Storage) *Manager {
	return &Manager{
		storage:   storage,
		validator: NewValidator(),
	}
}

// Store saves credentials to storage
func (m *Manager) Store(creds *Credentials) error {
	if !m.validator.Validate(creds) {
		return fmt.Errorf("invalid credentials format")
	}

	if err := m.storage.Set("login42", creds.Login42); err != nil {
		return fmt.Errorf("failed to store login42: %w", err)
	}

	if err := m.storage.Set("client_uid", creds.ClientID); err != nil {
		return fmt.Errorf("failed to store client_uid: %w", err)
	}

	if err := m.storage.Set("client_secret", creds.ClientSecret); err != nil {
		return fmt.Errorf("failed to store client_secret: %w", err)
	}

	return nil
}

// GetLogin retrieves the stored login42
func (m *Manager) GetLogin() (string, error) {
	return m.storage.Get("login42")
}

// GetClientID retrieves the stored client ID
func (m *Manager) GetClientID() (string, error) {
	return m.storage.Get("client_uid")
}

// GetClientSecret retrieves the stored client secret
func (m *Manager) GetClientSecret() (string, error) {
	return m.storage.Get("client_secret")
}

// GetAll retrieves all credentials
func (m *Manager) GetAll() (*Credentials, error) {
	login, err := m.GetLogin()
	if err != nil {
		return nil, fmt.Errorf("failed to get login42: %w", err)
	}

	clientID, err := m.GetClientID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client_uid: %w", err)
	}

	clientSecret, err := m.GetClientSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to get client_secret: %w", err)
	}

	return &Credentials{
		Login42:      login,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

// Clear removes all credentials from storage
func (m *Manager) Clear() error {
	if err := m.storage.DeleteAll(); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}
	return nil
}

// Validate checks if credentials are valid
func (m *Manager) Validate(creds *Credentials) bool {
	return m.validator.Validate(creds)
}
