package main

import (
	"log"
	"net/http"

	"backend/internal/db"
	"backend/internal/handlers"
)

func main() {
	// Connect to the database
	database := db.ConnectDB("mongodb://quickybooster:3yEjKoLPOUjAed5o@115.78.15.110")
	defer db.DisconnectDB()

	// Initialize handlers with the database
	handlers.InitHandlers(database)

	// Set up routes
	http.HandleFunc("/addVocabulary", handlers.AddVocabulary)
	http.HandleFunc("/getVocabularies", handlers.GetVocabularies)
	http.HandleFunc("/practiceVocabularies", handlers.PracticeVocabularies)

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}