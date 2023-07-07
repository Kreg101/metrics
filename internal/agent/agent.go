package agent

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Agent struct {
	updateFreq time.Duration
	sendFreq   time.Duration
	host       string
	stats      runtime.MemStats
	client     http.Client
}

func NewAgent(update time.Duration, send time.Duration, host string) *Agent {
	agent := &Agent{updateFreq: update, sendFreq: send, host: host, client: http.Client{}}
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
	var pollCount int
	go func() {
		time.Sleep(a.sendFreq)
		for k, v := range getMapOfStats(&a.stats) {
			resp, err := a.client.Post(a.host+"/update/gauge/"+k+"/"+v, "text/plain", nil)
			if err != nil {
				fmt.Println(err)
			}
			_, err = io.Copy(io.Discard, resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
		}
		_, err := a.client.Post(a.host+"/update/counter/PollCount/"+fmt.Sprintf("%d", pollCount), "text/plain", nil)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("i'm here")
	}()

	for {
		runtime.ReadMemStats(&a.stats)
		time.Sleep(a.updateFreq)
		pollCount++
	}
}
