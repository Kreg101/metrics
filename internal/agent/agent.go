package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kreg101/metrics/internal/communication"
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
}

func NewAgent(update int, send int, host string) *Agent {
	agent := &Agent{
		updateFreq: time.Duration(update) * time.Second,
		sendFreq:   time.Duration(send) * time.Second,
		host:       host,
	}
	return agent
}

func getMapOfStats(stats *runtime.MemStats) map[string]float64 {
	res := make(map[string]float64)
	res["Alloc"] = float64(stats.Alloc)
	res["BuckHashSys"] = float64(stats.BuckHashSys)
	res["Frees"] = float64(stats.Frees)
	res["GCCPUFraction"] = float64(stats.GCCPUFraction)
	res["GCSys"] = float64(stats.GCSys)
	res["HeapAlloc"] = float64(stats.HeapAlloc)
	res["HeapIdle"] = float64(stats.HeapIdle)
	res["HeapInuse"] = float64(stats.HeapInuse)
	res["HeapObject"] = float64(stats.HeapObjects)
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
	var pollCount int64

	go func() {
		for range time.Tick(a.sendFreq) {
			for k, v := range getMapOfStats(&a.stats) {
				url := a.host + "/update/gauge/" + k + "/" + fmt.Sprintf("%f", v)
				go func(url string) {
					resp, err := http.Post(url, "text/plain", nil)
					if err != nil {
						fmt.Println(err)
					}
					defer resp.Body.Close()
				}(url)

				m := communication.Metrics{
					ID:    k,
					MType: "gauge",
					Value: &v,
				}
				res, err := json.Marshal(m)
				if err != nil {
					fmt.Println(err)
				}

				go func(url string, js []byte) {
					resp, e := http.Post(url, "application/json", bytes.NewBuffer(js))
					if e != nil {
						fmt.Println(e)
					}
					fmt.Println(resp)
					defer resp.Body.Close()
				}(a.host+"/update/", res)

				m.Value = nil
				res, err = json.Marshal(m)
				if err != nil {
					fmt.Println(err)
				}
				go func(url string, js []byte) {
					resp, e := http.Post(url, "application/json", bytes.NewBuffer(js))
					if e != nil {
						fmt.Println(e)
					}
					var m communication.Metrics
					_ = json.NewDecoder(resp.Body).Decode(&m)
					fmt.Println(m.ID, m.Value)
					defer resp.Body.Close()
				}(a.host+"/value/", res)

			}

			url := a.host + "/update/counter/PollCount/" + fmt.Sprintf("%d", pollCount)
			go func(url string) {
				resp, err := http.Post(url, "text/plain", nil)
				if err != nil {
					fmt.Println(err)
				}
				defer resp.Body.Close()
			}(url)

			m := communication.Metrics{
				ID:    "PollCount",
				MType: "counter",
				Delta: &pollCount,
			}
			res, err := json.Marshal(m)
			if err != nil {
				fmt.Println(err)
			}

			go func(url string, js []byte) {
				resp, e := http.Post(url, "application/json", bytes.NewBuffer(js))
				if e != nil {
					fmt.Println(e)
				}
				fmt.Println(resp)
				defer resp.Body.Close()
			}(a.host+"/update/", res)

			m.Delta = nil
			res, err = json.Marshal(m)
			if err != nil {
				fmt.Println(err)
			}

			go func(url string, js []byte) {
				resp, e := http.Post(url, "application/json", bytes.NewBuffer(js))
				if e != nil {
					fmt.Println(e)
				}
				var m communication.Metrics
				_ = json.NewDecoder(resp.Body).Decode(&m)
				fmt.Println(m.ID, m.Delta)
				defer resp.Body.Close()
			}(a.host+"/value/", res)
		}
	}()

	for range time.Tick(a.updateFreq) {
		runtime.ReadMemStats(&a.stats)
		pollCount++
	}
}
