package http

import (
	"net/http"
)

func StartServer() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
