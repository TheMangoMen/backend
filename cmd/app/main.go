package main

import (
	"log"
	"net/http"
	"time"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/email"
	"github.com/TheMangoMen/backend/internal/handler"
	"github.com/TheMangoMen/backend/internal/ratelimit"
	"github.com/TheMangoMen/backend/internal/store"
	"golang.org/x/time/rate"

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

	adminMiddleware := func(next http.Handler) http.Handler {
		return auther.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := auth.MustFromContext(r.Context())
			if !user.Admin {
				http.Error(w, "access forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}

	// regeneration rate of 1 token / 2 sec = 30 tokens/min
	// starting with the max of 10 tokens in reserve
	rl := ratelimit.NewMapRateLimiter[string](rate.Every(time.Second*2), 10)
	authedRatelimit := ratelimit.AuthedMiddleware(rl)

	router.Handle("POST /login/{uID}", handler.LogIn(auther, s, resendClient))

	router.Handle("GET /ranking/{jID}", auther.Middleware(handler.GetRanking(s)))
	router.Handle("POST /rankings", auther.Middleware(handler.AddRanking(s)))

	router.Handle("GET /jobs", auther.MiddlewareOptional(handler.GetJobs(s)))

	router.Handle("GET /user", auther.Middleware(handler.GetUser(s)))

	router.Handle("GET /contribution/{jID}", auther.Middleware(handler.GetContribution(s)))
	router.Handle("POST /contribution", auther.Middleware(authedRatelimit(handler.AddContribution(s))))

	router.Handle("POST /watching", auther.Middleware(authedRatelimit(handler.UpdateWatching(s))))

	router.Handle("GET /analytics/status_counts", auther.Middleware(handler.GetWatchedStatusCounts(s)))

	router.Handle("GET /stage", auther.MiddlewareOptional(handler.GetStage(s)))

	router.Handle("POST /admin/stage", adminMiddleware(handler.UpdateStage(s)))
	router.Handle("POST /admin/year", adminMiddleware(handler.UpdateYear(s)))
	router.Handle("POST /admin/season", adminMiddleware(handler.UpdateSeason(s)))
	router.Handle("POST /admin/cycle", adminMiddleware(handler.UpdateCycle(s)))
	router.Handle("GET /admin/stage", adminMiddleware(handler.GetStage(s)))
	router.Handle("GET /admin/year", adminMiddleware(handler.GetYear(s)))
	router.Handle("GET /admin/season", adminMiddleware(handler.GetSeason(s)))
	router.Handle("GET /admin/cycle", adminMiddleware(handler.GetCycle(s)))
	router.Handle("GET /admin/contributions", adminMiddleware(handler.GetContributionLogs(s)))

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
