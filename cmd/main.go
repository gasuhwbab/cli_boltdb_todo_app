package main

import (
	"log"

	"github.com/gasuhwbab/cli_todo_app/internal/cli"
	"github.com/gasuhwbab/cli_todo_app/internal/db"
)

func main() {
	if err := db.Db.StartDb(); err != nil {
		log.Fatal(err)
	}
	defer db.Db.Close()
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
