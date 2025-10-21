package main

import (
	"net/http"

	serv "github.com/JennyBashir/final_project/http"
)

func main() {
	serv.StartServer()
	http.HandleFunc(`\`, serv.Handlers)

}
