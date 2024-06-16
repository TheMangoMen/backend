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

	router.Handle("GET /login/{uID}", handler.LogIn(auther, emailClient))
	router.Handle("GET /user", ensureAuth(handler.GetUser(s)))

	router.HandleFunc("GET /jobs", handler.GetJobs(s))

	router.HandleFunc("GET /rankings/{jid}", handler.GetRankings(s))
	router.HandleFunc("POST /rankings", handler.AddRanking(s))

	router.HandleFunc("GET /contribution", handler.GetContribution(s))
	router.HandleFunc("POST /contribution", handler.AddContribution(s))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}
