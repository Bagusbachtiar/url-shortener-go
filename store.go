package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

// URLStore holds our in-memory map protected by a Read/Write Mutex
type URLStore struct {
	mu   sync.RWMutex
	urls map[string]string
}

// NewURLStore initializes an empty thread-safe store
func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]string),
	}
}

// generateCode creates a random 6-character hex string (e.g., "a1b2c3")
func generateCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Save locks the map for writing, generates a short code, and stores the original URL
func (s *URLStore) Save(originalURL string) string {
	s.mu.Lock()         // Lock the door so no one else can read/write during update
	defer s.mu.Unlock() // Unlock automatically when the function finishes

	code := generateCode()
	s.urls[code] = originalURL
	return code
}

// Get locks the map for reading and returns the original URL if found
func (s *URLStore) Get(code string) (string, bool) {
	s.mu.RLock()         // Read-Lock allows multiple users to read at the same time safely
	defer s.mu.RUnlock() // Release Read-Lock when done

	originalURL, exists := s.urls[code]
	return originalURL, exists
}
