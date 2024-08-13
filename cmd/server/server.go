package main

import (
	"atmail-demo/api"
	"atmail-demo/config"
	"atmail-demo/database"
	"fmt"
	"github.com/nokusukun/faust"
	"log"
	"net/http"
)

func main() {

	app := faust.New(faust.APIInfo{
		Title:   "Users API",
		Summary: "Sample API for managing users",
		Version: "0.0.1",
	})

	db, err := database.NewDatabase(config.CONNECTION_STRING)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}

	api.UsersEndpoint(app, db)

	log.Println("Listening on", config.PORT)
	err = http.ListenAndServe(config.PORT, app)
	if err != nil {
		fmt.Println(err)
		return
	}
}
