package main

import (
	"log"
	"net/http"

	"github.com/TheMangoMen/backend/internal/handler"
	"github.com/TheMangoMen/backend/internal/store"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Production
	// db, err := sqlx.Connect("pgx", "postgres://admin:password@172.19.134.43:5432/Milestone1?sslmode=disable")
	// Local
	db, err := sqlx.Connect("pgx", "postgres://admin@localhost:5432/Milestone1?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	s := store.NewStore(db)

	router := http.NewServeMux()
	router.HandleFunc("GET /user/{uID}", handler.GetUser(s))
	router.HandleFunc("GET /jobs", handler.GetJobs(s))
	router.HandleFunc("GET /rankings/{jid}", handler.GetRankings(s))
	router.HandleFunc("POST /rankings", handler.AddRanking(s))
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}
