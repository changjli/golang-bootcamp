package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const FILE_PATH = "weather_stations.csv"
const WORKERS = 3

// ReadingStats holds the aggregated statistics for a weather station.
type ReadingStats struct {
	Min, Max, Sum float64
	Count         int64
}

// processChunk reads and processes a chunk of the CSV file.
func processChunk(chunk []byte) map[string]ReadingStats {
	localStats := make(map[string]ReadingStats)
	lines := strings.Split(string(chunk), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			continue
		}
		station := parts[0]
		temp, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}

		stats, ok := localStats[station]
		if !ok {
			stats = ReadingStats{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		} else {
			if temp < stats.Min {
				stats.Min = temp
			}
			if temp > stats.Max {
				stats.Max = temp
			}
			stats.Sum += temp
			stats.Count++
		}
		localStats[station] = stats
	}
	return localStats
}

func main() {
	file, err := os.Open(FILE_PATH)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}
	fileSize := fileInfo.Size()

	numChunks := WORKERS
	chunkSize := fileSize / int64(numChunks)

	// Pipeline 1 : Assign chunks to workers 

	var wg sync.WaitGroup
	resultsChan := make(chan map[string]ReadingStats, numChunks)

	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()
			offset := int64(chunkIndex) * chunkSize
			// For the last chunk, read until the end of the file
			size := chunkSize
			if chunkIndex == numChunks-1 {
				size = fileSize - offset
			}

			buffer := make([]byte, size)
			_, err := file.ReadAt(buffer, offset)
			if err != nil && err != io.EOF {
				log.Printf("Error reading chunk %d: %v", chunkIndex, err)
				return
			}
			resultsChan <- processChunk(buffer)
		}(i)
	}

	wg.Wait()
	close(resultsChan)

	// Merge the results from all goroutines
	globalStats := make(map[string]ReadingStats)
	for localStats := range resultsChan {
		for station, stats := range localStats {
			global, ok := globalStats[station]
			if !ok {
				globalStats[station] = stats
			} else {
				if stats.Min < global.Min {
					global.Min = stats.Min
				}
				if stats.Max > global.Max {
					global.Max = stats.Max
				}
				global.Sum += stats.Sum
				global.Count += stats.Count
				globalStats[station] = global
			}
		}
	}

	// Sort and print the final results
	stations := make([]string, 0, len(globalStats))
	for station := range globalStats {
		stations = append(stations, station)
	}
	sort.Strings(stations)

	fmt.Println("Aggregated Weather Readings:")
	for _, station := range stations {
		stats := globalStats[station]
		avg := stats.Sum / float64(stats.Count)
		fmt.Printf("%s: Min=%.2f, Max=%.2f, Avg=%.2f\n", station, stats.Min, stats.Max, avg)
	}
}
