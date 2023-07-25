package main

import (
	"github.com/Kreg101/metrics/internal/agent"
	"time"
)

func main() {

	parseFlags()
	for i := 0; i < 10; i++ {
		go func() {
			a := agent.NewAgent(flagPollInterval, flagReportInterval, "http://"+flagEndpoint)
			a.Start()
		}()
	}

	time.Sleep(1000 * time.Second)

}
