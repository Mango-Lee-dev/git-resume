package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrInvalidAPIKey   = errors.New("invalid API key")
)

// Session represents a user session with their API key
type Session struct {
	ID        string    `json:"id"`
	APIKey    string    `json:"-"` // Never expose in JSON
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	LastUsed  time.Time `json:"last_used"`
}

// Manager handles session lifecycle
type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	ttl      time.Duration
}

// NewManager creates a new session manager
func NewManager(ttl time.Duration) *Manager {
	m := &Manager{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}

	// Start cleanup goroutine
	go m.cleanup()

	return m
}

// CreateSession validates the API key and creates a new session
func (m *Manager) CreateSession(apiKey string) (*Session, error) {
	// Validate API key format
	if len(apiKey) < 20 || apiKey[:7] != "sk-ant-" {
		return nil, ErrInvalidAPIKey
	}

	// Validate API key by making a test request to Anthropic
	if err := m.validateAPIKey(apiKey); err != nil {
		return nil, err
	}

	session := &Session{
		ID:        uuid.New().String(),
		APIKey:    apiKey,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.ttl),
		LastUsed:  time.Now(),
	}

	m.mu.Lock()
	m.sessions[session.ID] = session
	m.mu.Unlock()

	return session, nil
}

// GetSession retrieves a session by ID
func (m *Manager) GetSession(id string) (*Session, error) {
	m.mu.RLock()
	session, exists := m.sessions[id]
	m.mu.RUnlock()

	if !exists {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		m.DeleteSession(id)
		return nil, ErrSessionExpired
	}

	// Update last used
	m.mu.Lock()
	session.LastUsed = time.Now()
	m.mu.Unlock()

	return session, nil
}

// GetAPIKey retrieves the API key for a session
func (m *Manager) GetAPIKey(sessionID string) (string, error) {
	session, err := m.GetSession(sessionID)
	if err != nil {
		return "", err
	}
	return session.APIKey, nil
}

// DeleteSession removes a session
func (m *Manager) DeleteSession(id string) {
	m.mu.Lock()
	delete(m.sessions, id)
	m.mu.Unlock()
}

// ExtendSession extends the session expiration
func (m *Manager) ExtendSession(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[id]
	if !exists {
		return ErrSessionNotFound
	}

	session.ExpiresAt = time.Now().Add(m.ttl)
	session.LastUsed = time.Now()

	return nil
}

// validateAPIKey makes a minimal API call to verify the key is valid
func (m *Manager) validateAPIKey(apiKey string) error {
	reqBody := map[string]interface{}{
		"model":      "claude-3-haiku-20240307",
		"max_tokens": 1,
		"messages": []map[string]string{
			{"role": "user", "content": "hi"},
		},
	}

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate API key: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		// Success
		return nil
	case 401:
		return ErrInvalidAPIKey
	case 403:
		return ErrInvalidAPIKey
	case 429:
		// Rate limited but key is valid
		return nil
	case 529:
		// API overloaded but key is valid
		return nil
	default:
		// For other errors, assume key might be valid
		// Real validation happens when making actual API calls
		return nil
	}
}

// cleanup periodically removes expired sessions
func (m *Manager) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for id, session := range m.sessions {
			if now.After(session.ExpiresAt) {
				delete(m.sessions, id)
			}
		}
		m.mu.Unlock()
	}
}
