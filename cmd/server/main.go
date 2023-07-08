package main

import (
	"github.com/Kreg101/metrics/internal/server"
)

func main() {

	parseFlags()

	s := server.CreateNewServer()
	err := s.ListenAndServe(endpoint)
	if err != nil {
		panic(err)
	}

}

