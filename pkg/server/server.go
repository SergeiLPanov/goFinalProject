package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"goFinalProject/pkg/api"
)

const (
	webDir      = "./web"
	defaultPort = 7540
)

func Run() error {
	port := defaultPort
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}
