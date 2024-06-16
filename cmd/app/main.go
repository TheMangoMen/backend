package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TheMangoMen/backend/internal/handler"
	"github.com/TheMangoMen/backend/internal/store"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	DB_CONNECTION := os.Getenv("DB_CONNECTION")
	db, err := sqlx.Connect("pgx", DB_CONNECTION)
	if err != nil {
		log.Fatalln(err)
	}

	s := store.NewStore(db)

	router := http.NewServeMux()
	router.HandleFunc("GET /user/{uID}", handler.GetUser(s))
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
