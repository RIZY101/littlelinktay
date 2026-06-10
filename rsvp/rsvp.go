package main

import (
	"fmt"
	"io"
	"os"
	"log"
	"net/http"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Set allowed origins (Use specific domains instead of "*" in production for security)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Define allowable HTTP methods for cross-origin transactions
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		
		// Define acceptable request headers (essential if sending JSON or Auth tokens)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle the browser's preflight OPTIONS probe immediately
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		// Forward execution to your actual target handler
		next(w, r)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and process the payload body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	file, err := os.OpenFile("rsvp.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Write data to the file
	textToAppend := string(body) + "\n"
	if _, err := file.WriteString(textToAppend); err != nil {
		log.Fatalf("failed to write to file: %s", err)
	}
	fmt.Printf("Received payload: %s\n", string(body))
	
	// Get html response from file
	bytes, err := os.ReadFile("rsvp.html")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	
	// Set response headers+body
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func main() {

	mux := http.NewServeMux()
	
	// Bind the POST handler
	mux.HandleFunc("/", enableCORS(handlePostRequest))

	fmt.Println("Server executing smoothly on :1337...")
	if err := http.ListenAndServe(":1337", mux); err != nil {
		fmt.Printf("Server failure: %v\n", err)
	}
}
