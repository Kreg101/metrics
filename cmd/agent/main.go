package main

import (
	"github.com/Kreg101/metrics/internal/agent"
)

func main() {

	parseFlags()
	a := agent.NewAgent(flagPollInterval, flagReportInterval, "http://"+flagEndpoint)
	a.Start()

}
