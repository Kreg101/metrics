package main

import (
	"flag"
	"github.com/Kreg101/metrics/internal/server"
)

var (
	endpoint string
)

func parseFlags() {
	flag.StringVar(&endpoint, "a", ":8080", "address and port to run server")
	flag.Parse()
}

func main() {

	parseFlags()

	s := server.CreateNewServer()
	err := s.ListenAndServe(endpoint)
	if err != nil {
		panic(err)
	}

}
