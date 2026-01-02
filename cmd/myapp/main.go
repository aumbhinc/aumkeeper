package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var tmplHome *template.Template

func main() {
	fmt.Println("🔄 Starting AumKeeper server...")

	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found or failed to load")
	}

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture OS signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		log.Printf("🛑 Received signal: %v", sig)
		cancel()
	}()

	// Load HTML templates
	var err error
	tmplHome, err = template.ParseFiles("api/templates/base.html", "api/templates/home.html")
	if err != nil {
		log.Fatalf("❌ Failed to load templates: %v", err)
	}
	log.Println("✅ Templates loaded")

	// Router setup
	mux := http.NewServeMux()

	// Static assets
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("api/static")))
	mux.Handle("/static/", cacheControlMiddleware(staticHandler))

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	// Home page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Home",
			"Year":  time.Now().Year(),
			// Optional: you can add description, IncludeDashboardJS, etc.
			"Description": "Anonymous LLCs and Registered Agent services for solopreneurs who disrupt the status quo.",
			"IncludeDashboardJS": false,
		}
		if err := tmplHome.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("🌐 Listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	log.Println("🔽 Shutting down server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("❌ Graceful shutdown failed: %v", err)
	}
	log.Println("✅ Server stopped cleanly")
}

// Cache control for static assets
func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable") // 7 days
		next.ServeHTTP(w, r)
	})
}
