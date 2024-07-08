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
	DBConnectionAddr string `env:"DB_CONNECTION_ADDR,required"`
	AuthPrivateKey   string `env:"AUTH_PRIVATE_KEY,required"`

	FromEmail    string `env:"FROM_EMAIL,required"`
	ResendAPIKey string `env:"RESEND_API_KEY,required"`
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
	resendClient := email.NewResendClient(cfg.FromEmail, cfg.ResendAPIKey)

	router := http.NewServeMux()

	// HI NORMAN! THIS IS FOR YOU! UNCOMMENT IT AND WRAP A ROUTE LIKE USUAL, NO NEED TO USE auther.Middleware THIS WILL HANDLE THAT ALREADY :)
	// adminMiddleware := func(next http.Handler) http.Handler {
	// 	return auther.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		user := auth.MustFromContext(r.Context())
	// 		if !user.Admin {
	// 			http.Error(w, "access forbidden", http.StatusForbidden)
	// 			return
	// 		}
	// 		next.ServeHTTP(w, r)
	// 	}))
	// }

	router.Handle("POST /login/{uID}", handler.LogIn(auther, s, resendClient))

	router.Handle("GET /rankings/{jID}", handler.GetRankings(s))
	router.Handle("POST /rankings", auther.Middleware(handler.AddRanking(s)))

	router.Handle("GET /jobs", auther.MiddlewareOptional(handler.GetJobs(s)))

	router.Handle("GET /user", auther.Middleware(handler.GetUser(s)))

	router.Handle("GET /contribution/{jID}", auther.Middleware(handler.GetContribution(s)))
	router.Handle("POST /contribution", auther.Middleware(handler.AddContribution(s)))

	router.Handle("POST /watching", auther.Middleware(handler.UpdateWatching(s)))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	corsRouter := c.Handler(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: corsRouter,
	}
	log.Printf("Listening on %s", server.Addr)
	server.ListenAndServe()
}
