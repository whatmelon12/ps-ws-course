package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/whatmelon12/ps-ws-course/database"
	"github.com/whatmelon12/ps-ws-course/handlers"
)

const basePath = "/api"

func main() {
	database.SetupDatabase()
	handlers.SetupProductRoutes(basePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
