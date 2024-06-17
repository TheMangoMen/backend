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
	DBConnectionAddr  string `env:"DB_CONNECTION_ADDR,required"`
	AuthPrivateKey    string `env:"AUTH_PRIVATE_KEY,required"`
	FromEmail         string `env:"FROM_EMAIL,required"`
	FromEmailPassword string `env:"FROM_EMAIL_PASSWORD,required"`
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

	outlookClient := email.NewOutlookClient(cfg.FromEmail, cfg.FromEmailPassword)

	router := http.NewServeMux()

	router.Handle("GET /login/{uID}", handler.LogIn(auther, outlookClient))
	router.Handle("GET /user", ensureAuth(handler.GetUser(s)))

	router.Handle("GET /jobs", handler.GetJobs(s))

	router.Handle("GET /rankings/{jID}", handler.GetRankings(s))
	router.Handle("POST /rankings", handler.AddRanking(s))

	router.Handle("GET /contribution", handler.GetContribution(s))
	router.Handle("POST /contribution", handler.AddContribution(s))

	allowCORS := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	}
	corsRouter := http.NewServeMux()
	corsRouter.Handle("/", allowCORS(router))

	server := http.Server{
		Addr:    ":8080",
		Handler: corsRouter,
	}
	server.ListenAndServe()
}
