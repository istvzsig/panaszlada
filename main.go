package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/istvzsig/panaszlada/handlers"
	"github.com/istvzsig/panaszlada/internal/geo"
	"github.com/istvzsig/panaszlada/storage"

	"github.com/istvzsig/retryx"
	"github.com/istvzsig/retryx/breaker"
	"github.com/istvzsig/retryx/retry"

	"github.com/istvzsig/ratelx"
)

func main() {

	// =====================
	// ENV
	// =====================
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("no .env.local file (using system env)")
	}

	port := os.Getenv("PORT")
	apiKey := os.Getenv("API_KEY")

	if port == "" || apiKey == "" {
		log.Fatal("missing env vars")
	}

	// =====================
	// DB CONFIG
	// =====================
	var dsn string

	if os.Getenv("LOCAL_MODE") == "true" {
		dsn = "postgres://postgres:postgres@localhost:5432/panaszlada?sslmode=disable"
		log.Println("using LOCAL database")
	} else {
		dsn = os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL missing")
		}
		log.Println("using REMOTE database")
	}

	// =====================
	// DB CONNECT
	// =====================
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("db ping error:", err)
	}

	log.Println("database connected")

	// =====================
	// STORE
	// =====================

	pgStore := storage.NewPostgresStore(db)

	policy := retry.DefaultPolicy()

	wrapper := retryx.Wrapper{
		Retry: policy,
		Breaker: breaker.New(breaker.BreakerConfig{
			FailureThreshold: 5,
			SuccessThreshold: 2,
			OpenTimeout:      10 * time.Second,
		}),
	}

	store := storage.NewRetryStore(pgStore, wrapper)

	h := handlers.NewReportHandler(store, apiKey)

	// =====================
	// ROUTER
	// =====================
	r := mux.NewRouter()

	// =====================
	// RATE LIMIT
	// =====================
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip := r.RemoteAddr
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				ip = strings.Split(xff, ",")[0]
			}

			if !getLimiter(ip).Allow() {
				http.Error(w, "rate limited", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// =====================
	// API
	// =====================
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/reports", h.CreateReport).Methods("POST")
	api.HandleFunc("/reports", h.ListReports).Methods("GET")
	api.HandleFunc("/reports/{tracking_code}", h.GetReport).Methods("GET")

	// =====================
	// STATIC
	// =====================
	fs := http.FileServer(http.Dir("./frontend"))

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", fs),
	)

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		path := "./frontend" + r.URL.Path

		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			http.ServeFile(w, r, path)
			return
		}

		if strings.Contains(r.URL.Path, ".") {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "./frontend/index.html")
	})

	// =====================
	// CORS
	// =====================
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	// =====================
	// SERVER
	// =====================
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// =====================
	// GEO
	// =====================
	if err := geo.LoadHungary("./internal/assets/hungary.geojson"); err != nil {
		log.Fatal(err)
	}

	log.Printf("Panaszláda running on port=%s", port)
	log.Fatal(srv.ListenAndServe())
}

var (
	limiters   = make(map[string]*ratelx.Limiter)
	limitersMu sync.Mutex
)

func getLimiter(ip string) *ratelx.Limiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()

	l, ok := limiters[ip]
	if !ok {
		l = ratelx.New(5, 5.0, false)
		limiters[ip] = l
	}

	return l
}
