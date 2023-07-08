package main

import (
	"flag"
	"github.com/Kreg101/metrics/internal/agent"
)

var (
	endpoint       string
	reportInterval int
	pollInterval   int
)

func parseFlags() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "frequency of sending metrics on server in seconds")
	flag.IntVar(&pollInterval, "p", 2, "frequency of updating metrics in seconds")
	flag.Parse()
}

func main() {

	parseFlags()
	a := agent.NewAgent(pollInterval, reportInterval, "http://"+endpoint)
	a.Start()

}
