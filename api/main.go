package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// User struct to hold user data
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

const version = "1.0.0"

// Database connection (global for simplicity)
var db *sql.DB

// Health check handler for /status route
func getStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "ok", "version": "`+version+`"}`)
}

// User handler for /user route with a query parameter ?q=
func getUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Query the database for the user
	var user User
	err := db.QueryRow("SELECT id, name, email FROM users WHERE name = ?", query).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			log.Printf("User not found: %s", query)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println("Database query error:", err)
		}
		return
	}

	// Respond with user data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	log.Printf("User retrieved: %s", user.Name)
}

// Initialize logging to file
func initLogging(logFile string) {
	// Open or create the log file in append mode
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Set output of log functions to the file
	log.SetOutput(file)

	// Optionally, set log format flags
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Include timestamp and file:line
}

func allowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request for debugging
		log.Printf("Received %s request for %s\n", r.Method, r.URL)

		w.Header().Set("Access-Control-Allow-Origin", "*")             // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Allow specific methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow specific headers

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize logging to file
	initLogging("logs.txt")

	// Overwritten by pipeline
	con := "CONNECTION_STRING"

	var err error
	db, err = sql.Open("mysql", con)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Define routes with CORS middleware applied
	mux := http.NewServeMux()
	mux.HandleFunc("/status", getStatus)
	mux.HandleFunc("/user", getUser)

	// Wrap routes with CORS middleware
	handler := allowCors(mux)

	// Start server
	port := "8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed:", err)
	}
}
