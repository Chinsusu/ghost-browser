package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/user/ghost-browser/internal/app"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	// Initialize the app
	application := app.NewApp()
	application.Startup(nil)

	// API routes
	http.HandleFunc("/api/profiles", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == "OPTIONS" {
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			profiles, err := application.GetProfiles()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"profiles": profiles,
			})
		} else if r.Method == "POST" {
			// Create new random profile
			profile, err := application.GenerateRandomProfile()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(profile)
		}
	})

	http.HandleFunc("/api/proxies", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == "OPTIONS" {
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			proxies, err := application.GetProxies()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"proxies": proxies,
			})
		}
	})

	// Health check
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"message": "Ghost Browser API is running",
		})
	})

	// Serve frontend files from disk
	http.Handle("/", http.FileServer(http.Dir("./frontend/dist/")))

	log.Println("🚀 Ghost Browser API Server starting on http://localhost:8080")
	log.Println("📱 Frontend: http://localhost:8080")
	log.Println("🔧 API: http://localhost:8080/api/health")
	log.Println("👤 Profiles: http://localhost:8080/api/profiles")
	log.Println("🌐 Proxies: http://localhost:8080/api/proxies")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}