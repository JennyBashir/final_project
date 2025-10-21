package main

import (
	"log"

	serv "final_project/http"
)

func main() {
	err := serv.StartServer()
	if err != nil {
		log.Fatal("couldn't start the server", err)
	}
}
