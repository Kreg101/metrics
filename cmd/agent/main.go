package main

import "github.com/Kreg101/metrics/internal/agent"

func main() {

	parseFlags()

	a := agent.NewAgent(flagPollInterval, flagReportInterval, "http://"+flagEndpoint)
	a.Start()

	//x := 1.23
	//
	//m := metric.Metrics{
	//	ID:    "sys",
	//	MType: "gauge",
	//	Value: &x,
	//}
	//
	//res, err := json.Marshal(m)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(string(res))
	//
	//var s metric.Metrics
	//
	//err = json.Unmarshal(res, &s)
	//
	//fmt.Println(*s.Value)

}
