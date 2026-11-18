package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestURLStoreConcurrency(t *testing.T) {
	store := NewURLStore()
	var wg sync.WaitGroup

	// Spawn 100 concurrent goroutines writing and reading simultaneously
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			target := "https://example.com" + strconv.Itoa(id)

			// Write to the thread-safe store
			code := store.Save(target)

			// Read immediately back from the thread-safe store
			got, exists := store.Get(code)
			if !exists {
				t.Errorf("Expected code %s to exist", code)
			}
			if got != target {
				t.Errorf("Expected %s, got %s", target, got)
			}
		}(i)
	}

	wg.Wait()
}
