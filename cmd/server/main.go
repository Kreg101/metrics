package main

import (
	"github.com/Kreg101/metrics/internal/handler"
	"net/http"
)

func main() {
	
	mux := handler.NewMux()
	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}

}
