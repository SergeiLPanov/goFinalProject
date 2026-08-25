package main

import (
	"log"
	"os"

	"goFinalProject/pkg/db"
	"goFinalProject/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbFile := defaultDBFile
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	if err := db.Init(dbFile); err != nil {
		return err
	}
	defer db.Close()

	return server.Run()
}
