package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func StartServer() error {
	log.Println("starting the server")

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return fmt.Errorf("invalid port: %s", port)
	}

	webDir := filepath.Join(".", "web")
	_, err = os.Stat(webDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("the directory %s was not found", webDir)
	}

	http.Handle(`/`, http.FileServer(http.Dir(webDir)))

	addr := ":" + port
	log.Printf("the server is running on port %s", port)

	return http.ListenAndServe(addr, nil)
}
