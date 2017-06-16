package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/victorspringer/trapAdvisor/database"
	"github.com/victorspringer/trapAdvisor/routing"
)

func main() {
	var err error
	database.DB, err = database.Init("prod")
	if err != nil {
		panic(err)
	}
	defer database.DB.Close()

	router := routing.Router()

	log.Fatal(http.ListenAndServe(getPort(), router))
}

func getPort() string {
	if os.Getenv("PORT") != "" {
		return os.Getenv("PORT")
	}
	return ":8080"
}
