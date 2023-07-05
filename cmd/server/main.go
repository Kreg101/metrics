package main

import "github.com/Kreg101/metrics/internal/server"

func main() {

	s := server.CreateNewServer()
	err := s.ListenAndServe(`:8080`)
	if err != nil {
		panic(err)
	}

}
