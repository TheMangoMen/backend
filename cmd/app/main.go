package main

import (
	"log"
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/email"
	"github.com/TheMangoMen/backend/internal/handler"
	"github.com/TheMangoMen/backend/internal/store"

	"github.com/caarlos0/env/v11"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type config struct {
	DBConnectionAddr string `env:"DB_CONNECTION_ADDR,required"`
	AuthPrivateKey   string `env:"AUTH_PRIVATE_KEY,required"`
	ResendAPIKey     string `env:"RESEND_API_KEY"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panicf("error reading config\n%v\n", err)
	}
	cfg, err := env.ParseAs[config]()
	if err != nil {
		log.Panicf("error reading config\n%v\n", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DBConnectionAddr)
	if err != nil {
		log.Fatalln(err)
	}

	s := store.NewStore(db)

	auther := auth.NewAuth(cfg.AuthPrivateKey)
	ensureAuth := auther.Middleware()

	emailClient := email.NewEmailClient(cfg.ResendAPIKey, "hello@watrank.com")

	router := http.NewServeMux()

	allowCORS := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			next.ServeHTTP(w, r)
		})
	}

	router.Handle("GET /login/{uID}", allowCORS(handler.LogIn(auther, emailClient)))
	router.Handle("GET /user", ensureAuth(handler.GetUser(s)))

	router.Handle("GET /jobs", allowCORS(handler.GetJobs(s)))

	router.Handle("GET /rankings/{jid}", allowCORS(handler.GetRankings(s)))
	router.Handle("POST /rankings", allowCORS(handler.AddRanking(s)))

	router.Handle("GET /contribution", allowCORS(handler.GetContribution(s)))
	router.Handle("POST /contribution", allowCORS(handler.AddContribution(s)))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}
