package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

//URLStore holds our in-memmory map protected by a Read/Write Mutex
type URLStore struct {
	mu	sync.RWMutex
	urls map[string]string
}

//NewURLStore initialize  an empty thread-safe store 
func NewURLStore() *URLStore{
	return &URLStore{
		urls: make(map[string]string),
	}
}

//generatedCode creates a random 6-character hex string (e.g "a1b2c3")
func generatedCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

//Save locks the map for writing, generates a short code, and stores the origin URL
 func (s *URLStore) save(originalURL string) string {
	s.mu.Lock()	 // Lock the door so no one else can read/write during update
	defer s.mu.Unlock() // Unlock automatically when the function finishes

	code := generatedCode()
	s.urls[code] = originalURL
	return code
 }

 //Get Locks the map for reading and returns the original URL if found
 func