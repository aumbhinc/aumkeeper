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

	"github.com/aumbhinc/aumkeeper/api"
	"github.com/joho/godotenv"
)

var templates *template.Template

func main() {
	log.Println("🔄 Starting AumKeeper server...")

	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	loadTemplates()

	mux := http.NewServeMux()
	setupRoutes(mux)

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

func loadTemplates() {
	var err error
	templates = template.New("").Funcs(api.FuncMap())

	templates, err = templates.ParseGlob("api/templates/*.html")
	if err != nil {
		log.Fatalf("❌ Template parse error: %v", err)
	}

	log.Println("✅ Templates loaded")
}

func setupRoutes(mux *http.ServeMux) {
	static := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("api/static")),
	)
	mux.Handle("/static/", cacheControlMiddleware(static))

	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)
	mux.HandleFunc("/staffs", staffsHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", map[string]any{
		"Title": "Home",
		"Year":  time.Now().Year(),
	})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "dashboard", map[string]any{
		"Title":              "Dashboard",
		"Year":               time.Now().Year(),
		"IncludeDashboardJS": true,
	})
}

func staffsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "staffs", map[string]any{
		"Title": "Staff Manager",
		"Year":  time.Now().Year(),
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func renderTemplate(w http.ResponseWriter, name string, data map[string]any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		log.Println("❌ Render error:", err)
		http.Error(w, "Internal Server Error", 500)
	}
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
}

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		next.ServeHTTP(w, r)
	})
}
