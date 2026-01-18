package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aumkeeper/api"
	"aumkeeper/api/handlers"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("🔄 Starting AumKeeper server...")

	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	// Load templates globally
	templates := template.New("").Funcs(api.FuncMap())
	if _, err := templates.ParseGlob("api/templates/*.html"); err != nil {
		log.Fatalf("❌ Template parse error: %v", err)
	}
	log.Println("✅ Templates loaded")

	mux := http.NewServeMux()

	// Static assets
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("api/static")))
	mux.Handle("/static/", cacheControlMiddleware(static))

	// Health
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// CENTRALIZED PAGE HANDLERS
	mux.Handle("/", handlers.HomeHandler(templates))
	mux.Handle("/dashboard", handlers.DashboardHandler(templates))
	mux.Handle("/staffs", handlers.StaffsHandler(templates))

	server := &http.Server{
		Addr:              ":" + getPort(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("🌐 Listening on :%s", getPort())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	gracefulShutdown(server, 10*time.Second)
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	cancel()
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("✅ Server gracefully stopped")
}

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		next.ServeHTTP(w, r)
	})
}
