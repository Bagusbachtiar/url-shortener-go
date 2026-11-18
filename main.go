package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

// loggingMiddleware wraps any http.Handler to log method, path, and latency
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // Pass the request to the actual route handler
		log.Printf("⚡ [%s] %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	store := NewURLStore()
	mux := http.NewServeMux()

	// Route 1: POST /shorten -> Accepts a long URL and returns a short code
	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "Please provide a valid JSON body with 'url'", http.StatusBadRequest)
			return
		}

		code := store.Save(req.URL)
		resp := ShortenResponse{
			ShortCode: code,
			ShortURL:  fmt.Sprintf("http://localhost:8080/r/%s", code),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	})

	// Route 2: GET /r/{code} -> Looks up the short code and redirects the visitor
	mux.HandleFunc("GET /r/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code") // Extracts {code} directly from the URL path!

		originalURL, exists := store.Get(code)
		if !exists {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		// Redirect visitor with HTTP 302 Found
		http.Redirect(w, r, originalURL, http.StatusFound)
	})

	port := ":8080"
	log.Printf("🚀 URL Shortener API running on http://localhost%s", port)
	// Wrap the router inside our custom loggingMiddleware
	if err := http.ListenAndServe(port, loggingMiddleware(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
