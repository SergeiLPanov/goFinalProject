package main

import (
	"log"
	"os"

	"goFinalProject/pkg/db"
	"goFinalProject/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main() {
	dbFile := defaultDBFile
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
