package main

import (
	"flag"
	"os"
)

var (
	flagEndpoint string
)

func parseFlags() {
	flag.StringVar(&flagEndpoint, "a", ":8080", "address and port to run server")
	flag.Parse()
	if envEndpoint := os.Getenv("ADDRESS"); envEndpoint != "" {
		flagEndpoint = envEndpoint
	}
}
