package main

import (
	"log"

	"github.com/gasuhwbab/cli_todo_app/internal/cli"
	"github.com/gasuhwbab/cli_todo_app/internal/db"
)

const dbPath = "/Users/ruslanmuradov/github.com:gasuhwbab/cli_todo_app/data/todo_app.db"

func main() {
	storage := db.NewStorage(dbPath)
	if err := storage.StartDb(); err != nil {
		log.Fatal(err)
	}
	app := cli.NewCLI(storage)
	app.Run()
}
