package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	flagEndpoint       string
	key                string
	flagReportInterval int
	flagPollInterval   int
)

func parseFlags() {
	flag.StringVar(&flagEndpoint, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&key, "k", "", "key for hash")
	flag.IntVar(&flagReportInterval, "r", 10, "frequency of sending metrics on server in seconds")
	flag.IntVar(&flagPollInterval, "p", 2, "frequency of updating metrics in seconds")
	flag.Parse()

	if envEndpoint := os.Getenv("ADDRESS"); envEndpoint != "" {
		flagEndpoint = envEndpoint
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		key = envKey
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		res, err := strconv.Atoi(envReportInterval)
		if err != nil {
			fmt.Println(err)
		} else {
			flagReportInterval = res
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		res, err := strconv.Atoi(envPollInterval)
		if err != nil {
			fmt.Println(err)
		} else {
			flagPollInterval = res
		}
	}
}
