package cities

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

var BILLION = 1_000_000_000

//go:embed weather_stations.csv
var weatherStations string

var citiesCache = make(map[string]struct{})

func parseCachedCities(cities map[string]struct{}) []string {
	result := make([]string, len(cities))

	i := 0
	for k := range cities {
		result[i] = k
		i++
	}

	return result
}

func GetCities() []string {
	if len(citiesCache) > 0 {
		return parseCachedCities(citiesCache)
	}

	lines := strings.Split(weatherStations, "\n")

	for _, line := range lines {
		if strings.Contains(line, "#") {
			continue
		}

		city := strings.Split(line, ";")[0]

		if _, ok := citiesCache[city]; ok == false {
			citiesCache[city] = struct{}{}
		}
	}

	return parseCachedCities(citiesCache)
}

func getRandomList(array []string, size int) []string {
	resultArray := make([]string, size)

	for i := 0; i < size; i++ {
		randomIndex := rand.Intn(len(array))
		resultArray[i] = array[randomIndex]
	}

	return resultArray
}

func GenerateFromCities(cities []string, writeTo *bufio.Writer, totalLines ...int) error {
	rowsToGenerate := BILLION

	if len(totalLines) > 0 {
		rowsToGenerate = totalLines[0]
	}

	startTime := time.Now()

	batchSizeToWrite := 1_000_000

	if rowsToGenerate < batchSizeToWrite {
		batchSizeToWrite = rowsToGenerate / 2
	}

	totalProcessed := 0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	replyTo := make(chan *bytes.Buffer, 256)

	for i := 0; i < runtime.NumCPU(); i++ {
		go generateDataWorker(ctx, cities, batchSizeToWrite, replyTo)
	}

	progress := float64(0)

	for {
		select {
		case buffer := <-replyTo:
			if totalProcessed >= rowsToGenerate {
				fmt.Printf("Total Time: %s\n", time.Now().Sub(startTime))
				fmt.Printf("Total Processed: %d %d\n", totalProcessed, rowsToGenerate)
				err := writeTo.Flush()
				if err != nil {
					return err
				}
				return nil
			}

			_, err := writeTo.Write(buffer.Bytes())
			if err != nil {
				return err
			}

			totalProcessed = totalProcessed + batchSizeToWrite

			newProgress := float64(totalProcessed) / float64(rowsToGenerate) * 100

			if progress != math.Floor(newProgress) {
				progress = newProgress
				fmt.Printf("Total Processed: %d %.0f%%\n", totalProcessed, newProgress)
			}
		}
	}
}

func generateDataWorker(ctx context.Context, cities []string, batch int, replyTo chan *bytes.Buffer) {
	maxValue := 99.9
	minValue := -99.9

	for {
		select {
		case <-ctx.Done():
			return
		default:
			randomList := getRandomList(cities, batch)

			buffer := bytes.NewBuffer(make([]byte, 0, batch*100))

			for _, s := range randomList {
				value := minValue + rand.Float64()*(maxValue-minValue)
				buffer.WriteString(createRow(s, value))
			}

			replyTo <- buffer
		}
	}
}

func createRow(city string, value float64) string {
	return fmt.Sprintf("%s;%.2f\n", city, value)
}
