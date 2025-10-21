package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	serverURL    = "http://srv.msk01.gigacorp.local/_stats"
	pollInterval = 10 * time.Second
	maxErrors    = 3
)

func main() {
	errorCount := 0

	for {
		stats, err := fetchStats()
		if err != nil {
			errorCount++
			if errorCount >= maxErrors {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(pollInterval)
			continue
		}

		errorCount = 0
		checkMetrics(stats)
		time.Sleep(pollInterval)
	}
}

func fetchStats() ([]int, error) {
	resp, err := http.Get(serverURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error")
	}

	scanner := bufio.NewScanner(resp.Body)
	if !scanner.Scan() {
		return nil, fmt.Errorf("error")
	}

	line := strings.TrimSpace(scanner.Text())
	values := strings.Split(line, ",")

	if len(values) != 7 {
		return nil, fmt.Errorf("error")
	}

	stats := make([]int, 7)
	for i, val := range values {
		num, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			return nil, fmt.Errorf("error")
		}
		stats[i] = num
	}

	return stats, nil
}

func checkMetrics(stats []int) {
	loadAvg := stats[0]
	totalMemory := stats[1]
	usedMemory := stats[2]
	totalDisk := stats[3]
	usedDisk := stats[4]
	networkBandwidth := stats[5]
	networkUsage := stats[6]

	if loadAvg > 30 {
		fmt.Printf("Load Average is too high: %d\n", loadAvg)
	}

	if totalMemory > 0 {
		memoryUsagePercent := (usedMemory * 100) / totalMemory
		if memoryUsagePercent > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", memoryUsagePercent)
		}
	}

	if totalDisk > 0 {
		diskUsagePercent := (usedDisk * 100) / totalDisk
		if diskUsagePercent > 90 {
			freeSpaceMB := (totalDisk - usedDisk) / (1024 * 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeSpaceMB)
		}
	}

	if networkBandwidth > 0 {
		networkUsagePercent := (networkUsage * 100) / networkBandwidth
		if networkUsagePercent > 90 {
			freeBandwidthMbit := (networkBandwidth - networkUsage) * 8 / (1000 * 1000)
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeBandwidthMbit)
		}
	}
}
