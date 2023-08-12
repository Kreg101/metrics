package agent

import (
	"encoding/json"
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

func sendGaugeJSON(client *resty.Client, url string, k string, v float64) {
	m := metric.Metric{ID: k, MType: "gauge", Delta: nil, Value: &v}
	resp, err := client.R().SetBody(m).SetHeader("Content-Type", "application/json").Post(url)

	if err != nil {
		fmt.Printf("can't send gauge metric to server: %e\n", err)
	}

	m1 := metric.Metric{}
	err = json.Unmarshal(resp.Body(), &m1)
	if err != nil {
		fmt.Println("hah")
		return
	}

	fmt.Println(m1.ID, m1.MType, m1.Delta, *m1.Value)
}

func sendCounterJSON(client *resty.Client, url string, k string, v int64) {
	m := metric.Metric{ID: k, MType: "counter", Delta: &v, Value: nil}
	resp, err := client.R().SetBody(m).SetHeader("Content-Type", "application/json").Post(url)

	if err != nil {
		fmt.Printf("can't send metric to server: %e\n", err)
		return
	}

	m1 := metric.Metric{}
	err = json.Unmarshal(resp.Body(), &m1)
	if err != nil {
		fmt.Println("hah")
		return
	}

	fmt.Println(m1.ID, m1.MType, *m1.Delta, m1.Value)
}

func (a *Agent) Start() {
	client := resty.New()
	var pollCount int64
	url := a.host + "/updates/"

	go func() {
		for range time.Tick(a.sendFreq) {
			var metrics []metric.Metric

			for k, v := range getMapOfStats(a.stats) {
				m := metric.Metric{ID: k, MType: "gauge", Delta: nil, Value: new(float64)}
				*m.Value = v
				metrics = append(metrics, m)
			}

			m := metric.Metric{ID: "PollCount", MType: "counter", Delta: &pollCount, Value: nil}
			metrics = append(metrics, m)

			for _, m := range metrics {
				fmt.Println(m.ID, m.MType, m.Delta, m.Value)
			}

			_, err := client.R().SetBody(metrics).SetHeader("Content-Type", "application/json").Post(url)
			if err != nil {
				fmt.Printf("can't get correct response from server: %e", err)
			}

		}
	}()

	for range time.Tick(a.updateFreq) {
		runtime.ReadMemStats(&a.stats)
		pollCount++
	}
}
