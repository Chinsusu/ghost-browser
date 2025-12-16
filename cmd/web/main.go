package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/user/ghost-browser/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Initialize the app
	application := app.NewApp()
	application.Startup(nil)

	// Serve static files
	distFS, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		log.Fatal("Failed to get dist filesystem:", err)
	}

	// Setup routes
	http.Handle("/", http.FileServer(http.FS(distFS)))
	
	// API routes (simple REST API)
	http.HandleFunc("/api/profiles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			profiles, _ := application.GetProfiles()
			w.Write([]byte(`{"profiles": []}`)) // Simplified for now
		}
	})

	log.Println("Ghost Browser Web Server starting on http://localhost:8080")
	log.Println("Open your browser and go to: http://localhost:8080")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}