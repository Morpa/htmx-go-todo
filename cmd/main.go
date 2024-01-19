package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Morpa/htmx-go-todo/internal/repository"
	"github.com/Morpa/htmx-go-todo/internal/repository/dbrepo"
)

const port = 3000

type application struct {
	DB repository.DatabaseRepo
}

func main() {
	var app application

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.SqliteDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	err = parseTemplates()
	if err != nil {
		log.Panic(err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
