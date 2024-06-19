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
	"github.com/rs/cors"
)

type config struct {
	DBConnectionAddr  string `env:"DB_CONNECTION_ADDR,required"`
	AuthPrivateKey    string `env:"AUTH_PRIVATE_KEY,required"`
	FromEmail         string `env:"FROM_EMAIL,required"`
	FromEmailPassword string `env:"FROM_EMAIL_PASSWORD,required"`
}

func main() {
	// .env parsing is optional if we can parse env variables from the system
	if err := godotenv.Load(); err != nil {
		log.Printf("error reading config file\n%v\n", err)
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

	router.Handle("POST /login/{uID}", handler.LogIn(auther, outlookClient))

	router.Handle("GET /rankings/{jID}", handler.GetRankings(s))
	router.Handle("POST /rankings", handler.AddRanking(s))

	router.Handle("GET /jobs", auther.MiddlewareOptional(handler.GetJobs(s))(handler.GetJobs(s)))

	router.Handle("GET /user", ensureAuth(handler.GetUser(s)))

	router.Handle("GET /contribution/{jID}", ensureAuth(handler.GetContribution(s)))
	router.Handle("POST /contribution", ensureAuth(handler.AddContribution(s)))

	router.Handle("POST /watching", ensureAuth(handler.UpdateWatching(s)))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
	})
	corsRouter := c.Handler(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: corsRouter,
	}
	log.Printf("Listening on %s", server.Addr)
	server.ListenAndServe()
}
