package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	endpoint      string
	storagePath   string
	databaseDSN   string
	storeInterval int
	restore       bool
	useDB         bool
)

func parseConfiguration() {
	flag.StringVar(&endpoint, "a", ":8080", "address and port to run server")
	flag.StringVar(&storagePath, "f", "/tmp/metrics-db.json", "file to store metrics")
	flag.StringVar(&databaseDSN, "d", "", "string for database connection")
	flag.IntVar(&storeInterval, "i", 10, "interval for saving data on disk")
	flag.BoolVar(&restore, "r", true, "load metrics from file")
	flag.Parse()

	if envEndpoint := os.Getenv("ADDRESS"); envEndpoint != "" {
		endpoint = envEndpoint
	}

	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		storagePath = envStoragePath
	}

	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		databaseDSN = envDataBaseDSN
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		i, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			fmt.Println(err)
		} else {
			storeInterval = i
		}
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		r, err := strconv.ParseBool(envRestore)
		if err != nil {
			fmt.Println(err)
		} else {
			restore = r
		}
	}

	if databaseDSN != "" {
		useDB = true
	}
}
