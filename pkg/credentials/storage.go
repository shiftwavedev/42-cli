package credentials

import "github.com/zalando/go-keyring"

// Storage defines the interface for credential storage
type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	DeleteAll() error
}

// KeyringStorage implements Storage using the system keyring
type KeyringStorage struct {
	service string
}

// NewKeyringStorage creates a new KeyringStorage instance
func NewKeyringStorage(service string) *KeyringStorage {
	return &KeyringStorage{
		service: service,
	}
}

// Get retrieves a value from the keyring
func (k *KeyringStorage) Get(key string) (string, error) {
	return keyring.Get(k.service, key)
}

// Set stores a value in the keyring
func (k *KeyringStorage) Set(key, value string) error {
	return keyring.Set(k.service, key, value)
}

// Delete removes a value from the keyring
func (k *KeyringStorage) Delete(key string) error {
	return keyring.Delete(k.service, key)
}

// DeleteAll removes all values from the keyring
func (k *KeyringStorage) DeleteAll() error {
	return keyring.DeleteAll(k.service)
}
