package main

import (
	"log"
	"net/http"

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

	log.Fatal(http.ListenAndServe(":8080", router))
}
