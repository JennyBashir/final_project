package http

import (
	"fmt"
	"net/http"
)

func Handlers(res http.ResponseWriter, req *http.Request) {
	fmt.Println("получен запрос")
}
