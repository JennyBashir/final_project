package main

import (
	"log"

	"final_project/pkg/db"
	serv "final_project/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("database opening error", err)
	}
	//defer db.Close()

	err = serv.StartServer()
	if err != nil {
		log.Fatal("couldn't start the server", err)
	}
}
