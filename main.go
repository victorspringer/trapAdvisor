package main

import (
	"log"
	"net/http"
	"os"

	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/victorspringer/trapAdvisor/database"
	"github.com/victorspringer/trapAdvisor/routing"
)

func main() {
	var err error
	if os.Getenv("DOMAIN") == "" || os.Getenv("PORT") == "" ||
		os.Getenv("DB_USER") == "" || os.Getenv("DB_ADDR") == "" || os.Getenv("DB_PORT") == "" ||
		os.Getenv("FB_CLIENT_ID") == "" || os.Getenv("FB_CLIENT_SECRET") == "" {
		err = errors.New("empty environment variable(s)")
	}
	if err != nil {
		log.Fatal(err)
		return
	}

	database.DB, err = database.Init("prod")
	if err != nil {
		log.Fatal(err)
		return
	}
	database.DB.SetConnMaxLifetime(0)
	defer database.DB.Close()

	router := routing.Router()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
