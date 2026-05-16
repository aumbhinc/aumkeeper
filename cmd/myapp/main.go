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
	"aumkeeper/internal/render"
	"aumkeeper/internal/repository"
	"aumkeeper/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	log.Println("🔄 Starting AumKeeper server...")

	_ = godotenv.Load()

	// --------------------------------------------------
	// Root context + graceful shutdown
	// --------------------------------------------------

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleSignals(cancel)

	// --------------------------------------------------
	// Templates
	// --------------------------------------------------

	templates := template.New("").Funcs(api.FuncMap())

	if _, err := templates.ParseGlob("api/templates/*.html"); err != nil {
		log.Fatalf("❌ template parse error: %v", err)
	}

	if _, err := templates.ParseGlob("api/templates/partials/*.html"); err != nil {
		log.Fatalf("❌ partial template parse error: %v", err)
	}

	log.Println("✅ Templates loaded successfully")

	// --------------------------------------------------
	// Renderer
	// --------------------------------------------------

	renderer := render.NewRenderer(templates)

	// --------------------------------------------------
	// Database
	// --------------------------------------------------

	var pool *pgxpool.Pool

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {

		dbPool, err := pgxpool.New(ctx, dbURL)
		if err != nil {
			log.Fatalf("❌ database connection failed: %v", err)
		}

		if err := dbPool.Ping(ctx); err != nil {
			log.Fatalf("❌ database ping failed: %v", err)
		}

		pool = dbPool
		defer pool.Close()

		log.Println("✅ Postgres connected")

	} else {
		log.Println("⚠️ DATABASE_URL missing (running in local mode)")
	}

	// --------------------------------------------------
	// Repository + Service layer
	// --------------------------------------------------

	var staffService *services.StaffService

	if pool != nil {
		repos := repository.New(pool)
		staffRepo := repository.NewStaffRepository(repos)
		staffService = services.NewStaffService(staffRepo)

		log.Println("✅ Services initialized (DB mode)")
	} else {
		log.Println("⚠️ Services initialized in MOCK mode (no DB)")
	}

	// --------------------------------------------------
	// Handler layer (FIXED: ALL USE *render.Renderer)
	// --------------------------------------------------

	homeHandler := handlers.NewHomeHandler(renderer)
	dashboardHandler := handlers.NewDashboardHandler(renderer)
	staffsHandler := handlers.NewStaffsHandler(renderer)

	var addStaffHandler *handlers.AddStaffHandler
	if staffService != nil {
		addStaffHandler = handlers.NewAddStaffHandler(renderer, staffService)
	}

	// --------------------------------------------------
	// Router
	// --------------------------------------------------

	mux := http.NewServeMux()

	// Static files
	static := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("api/static")),
	)

	mux.Handle("/static/", cacheControlMiddleware(static))

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Routes
	mux.Handle("/", homeHandler)
	mux.Handle("/dashboard", dashboardHandler)
	mux.Handle("/staffs", staffsHandler)

	if addStaffHandler != nil {
		mux.Handle("/addstaff", addStaffHandler)
	} else {
		mux.HandleFunc("/addstaff", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "DB not configured", http.StatusServiceUnavailable)
		})
	}

	// --------------------------------------------------
	// Server
	// --------------------------------------------------

	port := getPort()

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("🌐 Server running on :%s", port)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Printf("❌ server error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()

	gracefulShutdown(server, 10*time.Second)
}

// --------------------------------------------------
// Helpers
// --------------------------------------------------

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

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ shutdown error: %v", err)
		return
	}

	log.Println("✅ Server gracefully stopped")
}

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		next.ServeHTTP(w, r)
	})
}