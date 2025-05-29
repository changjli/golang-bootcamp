package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const FILE_PATH = "weather_stations.csv"
const WORKERS = 3

// ReadingStats holds the aggregated statistics for a weather station.
type ReadingStats struct {
	Min, Max, Sum float64
	Count         int64
}

// Split the file to chunks
func makeChunk(workers int) <-chan []byte {
	chanOut := make(chan []byte, workers)

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

	chunkSize := fileSize / int64(workers)

	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()
			offset := int64(chunkIndex) * chunkSize
			// For the last chunk, read until the end of the file
			size := chunkSize
			if chunkIndex == workers-1 {
				size = fileSize - offset
			}

			buffer := make([]byte, size)
			_, err := file.ReadAt(buffer, offset)
			if err != nil && err != io.EOF {
				log.Printf("Error reading chunk %d: %v", chunkIndex, err)
				return
			}
			chanOut <- buffer
		}(i)
	}

	wg.Wait()
	close(chanOut)

	return chanOut
}

// Process seperated chunks
func processChunk(chanIn <-chan []byte, workers int) <-chan map[string]ReadingStats {
	chanOut := make(chan map[string]ReadingStats, workers)

	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()

			for chunk := range chanIn {
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
				chanOut <- localStats
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}

// Write result to csv non concurrent
func writeResultsToCSV(filePath string, stations []string, stats map[string]ReadingStats) {
	// Create the output file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create output CSV file: %v", err)
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{"Station", "Min", "Max", "Avg"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header to CSV: %v", err)
	}

	// Write the data rows
	for _, station := range stations {
		s := stats[station]
		avg := s.Sum / float64(s.Count)

		// Format data into a slice of strings
		record := []string{
			station,
			fmt.Sprintf("%.2f", s.Min),
			fmt.Sprintf("%.2f", s.Max),
			fmt.Sprintf("%.2f", avg),
		}

		if err := writer.Write(record); err != nil {
			log.Fatalf("Failed to write record to CSV: %v", err)
		}
	}
	log.Printf("Successfully wrote results to %s\n", filePath)
}

func main() {

	startTime := time.Now()
	log.Println("Processing started...")

	// Pipeline 1 : Assign chunks to workers
	chanMakeChunk := makeChunk(WORKERS)

	// Pipeline 2 : Proccess local chunk
	chanProccessChunk := processChunk(chanMakeChunk, WORKERS)

	// Pipeline 3 : Merge the results from all goroutines
	globalStats := make(map[string]ReadingStats)
	for localStats := range chanProccessChunk {
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

	duration := time.Since(startTime)
	log.Printf("Processing finished. Total processing time: %s\n", duration)

	// writeResultsToCSV("results.csv", stations, globalStats)
}
