package agent

import (
	"fmt"
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

func getMapOfStats(stats *runtime.MemStats) map[string]string {
	res := make(map[string]string)
	res["Alloc"] = fmt.Sprintf("%f", float64(stats.Alloc))
	res["BuckHashSys"] = fmt.Sprintf("%f", float64(stats.BuckHashSys))
	res["Frees"] = fmt.Sprintf("%f", float64(stats.Frees))
	res["GCCPUFraction"] = fmt.Sprintf("%f", float64(stats.GCCPUFraction))
	res["GCSys"] = fmt.Sprintf("%f", float64(stats.GCSys))
	res["HeapAlloc"] = fmt.Sprintf("%f", float64(stats.HeapAlloc))
	res["HeapIdle"] = fmt.Sprintf("%f", float64(stats.HeapIdle))
	res["HeapInuse"] = fmt.Sprintf("%f", float64(stats.HeapInuse))
	res["HeapObject"] = fmt.Sprintf("%f", float64(stats.HeapObjects))
	res["HeapReleased"] = fmt.Sprintf("%f", float64(stats.HeapReleased))
	res["HeapSys"] = fmt.Sprintf("%f", float64(stats.HeapSys))
	res["LastGC"] = fmt.Sprintf("%f", float64(stats.LastGC))
	res["Lookups"] = fmt.Sprintf("%f", float64(stats.Lookups))
	res["MCacheInuse"] = fmt.Sprintf("%f", float64(stats.MCacheInuse))
	res["MCacheSys"] = fmt.Sprintf("%f", float64(stats.MCacheSys))
	res["MSpanInuse"] = fmt.Sprintf("%f", float64(stats.MSpanInuse))
	res["MSpanSys"] = fmt.Sprintf("%f", float64(stats.MSpanSys))
	res["Mallocs"] = fmt.Sprintf("%f", float64(stats.Mallocs))
	res["NextGC"] = fmt.Sprintf("%f", float64(stats.NextGC))
	res["NumForcedGC"] = fmt.Sprintf("%f", float64(stats.NumForcedGC))
	res["NumGC"] = fmt.Sprintf("%f", float64(stats.NumGC))
	res["OtherSys"] = fmt.Sprintf("%f", float64(stats.OtherSys))
	res["PauseTotalNs"] = fmt.Sprintf("%f", float64(stats.PauseTotalNs))
	res["StackInuse"] = fmt.Sprintf("%f", float64(stats.StackInuse))
	res["StackSys"] = fmt.Sprintf("%f", float64(stats.StackSys))
	res["Sys"] = fmt.Sprintf("%f", float64(stats.Sys))
	res["TotalAlloc"] = fmt.Sprintf("%f", float64(stats.TotalAlloc))
	res["RandomValue"] = fmt.Sprintf("%f", rand.Float64())
	return res
}

func (a *Agent) Start() {
	client := resty.New()
	var pollCount int

	go func() {
		for range time.Tick(a.sendFreq) {
			for k, v := range getMapOfStats(&a.stats) {
				go func(host, k, v string, client *resty.Client) {
					_, err := client.R().Post(host + "/update/gauge/" + k + "/" + v)
					if err != nil {
						fmt.Println(err)
					}
				}(a.host, k, v, client)
			}
			go func(host string, client *resty.Client) {
				_, err := client.R().Post(host + "/update/counter/PollCount/" + fmt.Sprintf("%d", pollCount))
				if err != nil {
					fmt.Println(err)
				}
			}(a.host, client)
		}
	}()

	for range time.Tick(a.updateFreq) {
		runtime.ReadMemStats(&a.stats)
		pollCount++
	}
}
