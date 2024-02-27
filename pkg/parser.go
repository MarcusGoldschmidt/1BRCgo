package pkg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/MarcusGoldschmidt/1brcgo/pkg/unit"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type StationMetrics struct {
	Station string
	Min     float32
	Mean    float32
	Max     float32
}

type valueMessage struct {
	station string
	value   float32
}

func Parse(reader *bufio.Reader) ([]*StationMetrics, error) {
	bucket := NewMetricsBucket()
	valueChannel := make(chan []byte, 10)
	aggChannel := make(chan map[string]*metricsAggregate, 24)

	wg := &sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			for msg := range valueChannel {
				values := parseChunk(msg)
				aggChannel <- aggregateFromMsg(values)
			}
			wg.Done()
		}()
	}

	go func() {
		readChunk(reader, valueChannel)
		close(valueChannel)
		wg.Wait()
		close(aggChannel)
	}()

	for msg := range aggChannel {
		bucket.Aggregate(msg)
	}

	return bucket.GetMetrics(), nil
}

func parseChunk(chunk []byte) []*valueMessage {
	lines := bytes.Split(chunk, []byte{'\n'})

	toSend := make([]*valueMessage, 0, len(lines))

	for _, lineByte := range lines {
		line := string(lineByte)

		splitValues := strings.Split(line, ";")

		if len(splitValues) != 2 {
			continue
		}

		value, err := strconv.ParseFloat(splitValues[1], 32)
		if err != nil {
			log.Fatal(err)
		}

		toSend = append(toSend, &valueMessage{
			station: splitValues[0],
			value:   float32(value),
		})
	}

	return toSend
}

func readChunk(reader *bufio.Reader, chunkChan chan []byte, bufferSize ...unit.Size) {
	size := unit.MB * 100
	if len(bufferSize) > 0 {
		size = bufferSize[0]
	}

	var progress unit.Size = 0
	progressString := progress.ToString()

	buf := make([]byte, size)
	leftover := make([]byte, 0, size/2)

	for {
		readTotal, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}
		// If we have leftover from the previous read, we need to prepend it to the current buffer
		buf = buf[:readTotal]

		lastNewLineIndex := bytes.LastIndex(buf, []byte{'\n'})

		if lastNewLineIndex == -1 {
			leftover = append(leftover, buf...)
			continue
		}

		newLeft := buf[lastNewLineIndex+1:]

		toSend := make([]byte, 0, readTotal)
		toSend = append(leftover, buf[:lastNewLineIndex+1]...)

		leftover = make([]byte, len(buf[lastNewLineIndex+1:]))
		copy(leftover, newLeft)

		chunkChan <- toSend

		progress += unit.Size(len(toSend))

		if progressString != progress.ToString() {
			progressString = progress.ToString()
			fmt.Println("Progress: ", progressString)
		}
	}
}
