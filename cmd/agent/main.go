package main

import (
	"github.com/Kreg101/metrics/internal/agent"
	"time"
)

func main() {

	a := agent.NewAgent(2*time.Second, 10*time.Second, `http://localhost:8080`)
	a.Start()

}
