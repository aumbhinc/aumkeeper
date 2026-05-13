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
	"aumkeeper/internal/repository"
	"aumkeeper/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("🔄 Starting AumKeeper server...")

	// Load environment
	_ = godotenv.Load()

	// Root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	// -----------------------------
	// Load templates
	// -----------------------------
	templates := template.New("").Funcs(api.FuncMap())

	if _, err := templates.ParseGlob("api/templates/*.html"); err != nil {
		log.Fatalf("❌ template parse error: %v", err)
	}

	if _, err := templates.ParseGlob("api/templates/*/*.html"); err != nil {
		log.Fatalf("❌ template parse error: %v", err)
	}

	log.Println("✅ Templates loaded successfully")

	// -----------------------------
	// Database
	// -----------------------------
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
		log.Println("✅ Postgres connected")
	} else {
		log.Println("⚠️ DATABASE_URL missing (local mode)")
	}

	// -----------------------------
	// Dependency wiring
	// -----------------------------
	repos := repository.New(pool)

	staffRepo := repository.NewStaffRepository(repos)

	staffService := services.NewStaffService(
		staffRepo,
	)

	addStaffHandler := handlers.NewAddStaffHandler(
		templates,
		staffService,
	)

	// -----------------------------
	// Router
	// -----------------------------
	mux := http.NewServeMux()

	// Static assets
	static := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("api/static")),
	)

	mux.Handle(
		"/static/",
		cacheControlMiddleware(static),
	)

	// Health endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Routes
	mux.Handle("/", handlers.HomeHandler(templates))
	mux.Handle("/dashboard", handlers.DashboardHandler(templates))
	mux.Handle("/staffs", handlers.StaffsHandler(templates))
	mux.Handle("/addstaff", addStaffHandler)

	// -----------------------------
	// HTTP Server
	// -----------------------------
	port := getPort()

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("🌐 Server listening on :%s", port)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	gracefulShutdown(server, 10*time.Second)
}

// ------------------------------------------------
// Helpers
// ------------------------------------------------

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)

	signal.Notify(
		c,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-c
	cancel()
}

func gracefulShutdown(
	server *http.Server,
	timeout time.Duration,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ graceful shutdown error: %v", err)
		return
	}

	log.Println("✅ Server gracefully stopped")
}

func cacheControlMiddleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Header().Set(
			"Cache-Control",
			"public, max-age=604800, immutable",
		)

		next.ServeHTTP(w, r)
	})
}