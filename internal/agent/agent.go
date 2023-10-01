package agent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	key        string
	stats      runtime.MemStats
}

func NewAgent(update int, send int, host string, key string) *Agent {
	return &Agent{
		updateFreq: time.Duration(update) * time.Second,
		sendFreq:   time.Duration(send) * time.Second,
		host:       host,
		key:        key,
	}
}

func getMapOfStats(stats runtime.MemStats) map[string]float64 {
	res := make(map[string]float64)
	res["Alloc"] = float64(stats.Alloc)
	res["BuckHashSys"] = float64(stats.BuckHashSys)
	res["Frees"] = float64(stats.Frees)
	res["GCCPUFraction"] = stats.GCCPUFraction
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

func (a *Agent) hash(m []Metric) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(a.key))
	h.Write(b)
	src := h.Sum(nil)

	dst := make([]byte, hex.EncodedLen(len(src)))

	hex.Encode(dst, src)
	return string(dst), nil
}

func (a *Agent) Start() {
	client := resty.New()
	var pollCount int64

	go func() {
		for range time.Tick(a.sendFreq) {
			var metrics []Metric

			for k, v := range getMapOfStats(a.stats) {
				m := Metric{ID: k, MType: "gauge", Value: new(float64)}
				*m.Value = v
				metrics = append(metrics, m)
			}

			m := Metric{ID: "PollCount", MType: "counter", Delta: &pollCount}
			metrics = append(metrics, m)

			r := client.R().SetBody(metrics).
				SetHeader("Content-Type", "application/json").
				SetHeader("Accept-Encoding", "gzip")

			if a.key != "" {
				hash, err := a.hash(metrics)
				if err != nil {
					fmt.Printf("can't get hash %v\n", err)
				} else {
					r.SetHeader("HashSHA256", hash)
				}
			}

			_, err := r.Post(a.host + "/updates/")

			if err != nil {
				fmt.Printf("can't get correct response from server: %v\n", err)
			}
		}
	}()

	for range time.Tick(a.updateFreq) {
		runtime.ReadMemStats(&a.stats)
		pollCount++
	}
}
