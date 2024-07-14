package backend

import (
	"log"
	"net/http"
	"os"
)

// Handle the root path to serve the /frontend/ directory

func StartWebPage() {
	// Define the path to the frontend directory
	// Relative paths are iffy and fround upon in Go but this is just a demo
	frontendDir := "./frontend"

	// Check if the directory exists
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist", frontendDir)
	}

	// Serve the frontend directory
	fs := http.FileServer(http.Dir(frontendDir))
	http.Handle("/", fs)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
