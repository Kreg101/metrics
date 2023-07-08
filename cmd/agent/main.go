package main

import (
	"github.com/Kreg101/metrics/internal/agent"
)

func main() {

	parseFlags()
	a := agent.NewAgent(pollInterval, reportInterval, endpoint)
	a.Start()

}
