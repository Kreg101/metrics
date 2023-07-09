package main

import (
	"github.com/Kreg101/metrics/internal/server"
	"github.com/Kreg101/metrics/internal/server/storage"
)

func main() {

	parseFlags()

	s := server.CreateNewServer(storage.NewStorage())
	err := s.ListenAndServe(flagEndpoint)
	if err != nil {
		panic(err)
	}

}
