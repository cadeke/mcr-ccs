package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

// User struct to hold user data
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Database connection (global for simplicity)
var db *sql.DB

// Health check handler for /status route
func getStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "ok"}`)
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
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println("Database query error:", err)
		}
		return
	}

	// Respond with user data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CORS Middleware
// func corsMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
//         w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allow specific methods
//         w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow specific headers
//
//         // Handle preflight requests
//         if r.Method == http.MethodOptions {
//             w.WriteHeader(http.StatusNoContent)
//             return
//         }
//
//         // Call the next handler
//         next.ServeHTTP(w, r)
//     })
// }

func main() {
	// Overwritten by pipeline
	con := "[CONNECTION STRING]"
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

	// Define routes
	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/user", getUser)

	// Start server
	port := "8080"
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
