package agent

import (
	"fmt"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"runtime"
	"time"
)

type Agent struct {
	updateFreq time.Duration
	sendFreq   time.Duration
	host       string
	stats      runtime.MemStats
}

func NewAgent(update int, send int, host string) *Agent {
	agent := &Agent{
		updateFreq: time.Duration(update) * time.Second,
		sendFreq:   time.Duration(send) * time.Second,
		host:       host,
	}
	return agent
}

func getMapOfStats(stats runtime.MemStats) map[string]float64 {
	res := make(map[string]float64, 0)
	res["Alloc"] = float64(stats.Alloc)
	res["BuckHashSys"] = float64(stats.BuckHashSys)
	res["Frees"] = float64(stats.Frees)
	res["GCCPUFraction"] = float64(stats.GCCPUFraction)
	res["GCSys"] = float64(stats.GCSys)
	res["HeapAlloc"] = float64(stats.HeapAlloc)
	res["HeapIdle"] = float64(stats.HeapIdle)
	res["HeapInuse"] = float64(stats.HeapInuse)
	res["HeapObjects"] = float64(stats.HeapObjects)
	res["HeapReleased"] = float64(stats.HeapReleased)
	res["HeapSys"] = float64(stats.HeapSys)
	res["LastGC"] = float64(stats.LastGC)
	res["Lookups"] = float64(stats.Lookups)
	res["MCacheInuse"] = float64(stats.MCacheInuse)
	res["MCacheSys"] = float64(stats.MCacheSys)
	res["MSpanInuse"] = float64(stats.MSpanInuse)
	res["MSpanSys"] = float64(stats.MSpanSys)
	res["Mallocs"] = float64(stats.Mallocs)
	res["NextGC"] = float64(stats.NextGC)
	res["NumForcedGC"] = float64(stats.NumForcedGC)
	res["NumGC"] = float64(stats.NumGC)
	res["OtherSys"] = float64(stats.OtherSys)
	res["PauseTotalNs"] = float64(stats.PauseTotalNs)
	res["StackInuse"] = float64(stats.StackInuse)
	res["StackSys"] = float64(stats.StackSys)
	res["Sys"] = float64(stats.Sys)
	res["TotalAlloc"] = float64(stats.TotalAlloc)
	res["RandomValue"] = rand.Float64()
	return res
}

func (a *Agent) Start() {
	client := resty.New()
	var pollCount int64

	go func() {
		for range time.Tick(a.sendFreq) {
			url := a.host + "/update/"
			mp := getMapOfStats(a.stats)
			fmt.Println(mp)
			for k, v := range mp {
				//url := a.host + "/update/gauge/" + k + "/" + fmt.Sprintf("%f", v)
				//go func(url string, client *resty.Client) {
				//	_, err := client.R().Post(url)
				//	if err != nil {
				//		fmt.Println(err)
				//	}
				//}(url, client)

				m := metric.Metric{ID: k, MType: "gauge", Delta: nil, Value: &v}
				go func(url string, m metric.Metric) {
					resp, err := client.R().SetBody(m).SetHeader("Content-Type", "application/json").Post(url)

					if err != nil {
						fmt.Println(err)
					}

					//		//req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))
					//		//if err != nil {
					//		//	fmt.Println(err)
					//		//}
					//		//req.Header.Set()
					//		//c := http.Client{}
					//		//resp, err := c.Do(req)
					//		//if err != nil {
					//		//	fmt.Println(err)
					//		//}

					fmt.Println(resp.StatusCode())
				}(url, m)

			}

			url = a.host + "/update/counter/PollCount/" + fmt.Sprintf("%d", pollCount)
			go func(url string, client *resty.Client) {
				_, err := client.R().Post(url)
				if err != nil {
					fmt.Println(err)
				}
			}(url, client)
		}
	}()

	for range time.Tick(a.updateFreq) {
		runtime.ReadMemStats(&a.stats)
		pollCount++
	}
}
