package pkg

import (
	"bufio"
	"bytes"
	"github.com/MarcusGoldschmidt/1brcgo/pkg/unit"
	"log"
	"os"
	"sync"
	"testing"
)

func BenchmarkGenerateFromCities(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file, err := os.OpenFile("/Users/marcus/projetos/1brcgo/test2.csv", os.O_RDONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(file)

		_, err = Parse(reader)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestReadChunk(t *testing.T) {
	type args struct {
		data  string
		chunk unit.Size
	}
	tests := []struct {
		name     string
		args     args
		expected []string
	}{
		{
			name: "TestReadChunk",
			args: args{
				data:  "Aa;50.13\nBB;22.40\nCC;1.40",
				chunk: unit.B * 4,
			},
			expected: []string{
				"Aa;50.13\n",
				"BB;22.40\n",
				"CC;1.40",
			},
		},
		{
			name: "TestBigReadChunk",
			args: args{
				data:  "Aa;50.13\nBB;22.40\nCC;1.40\nDD;55.55",
				chunk: unit.B * 20,
			},
			expected: []string{
				"Aa;50.13\nBB;22.40\n",
				"CC;1.40\n",
				"DD;55.55",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunkChan := make(chan []byte, 1)
			wg := sync.WaitGroup{}

			wg.Add(1)

			byteReader := bytes.NewReader([]byte(tt.args.data))

			reader := bufio.NewReader(byteReader)

			go func() {
				defer close(chunkChan)

				err := readChunk(reader, chunkChan, tt.args.chunk)
				if err != nil {
					t.Error(err)
				}
			}()

			count := 0

			for chunk := range chunkChan {
				if count >= len(tt.expected) {
					t.Fatalf("Expected %d chunks, but got %d", len(tt.expected), count)
				}

				if tt.expected[count] != string(chunk) {
					t.Fatalf("Expected %s, but got %s count %d", tt.expected[count], string(chunk), count)
				}
				count++
			}

		})
	}
}

func TestParseChunk(t *testing.T) {
	type args struct {
		chunk    []byte
		sendChan chan *valueMessage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestParseChunk",
			args: args{
				chunk:    []byte("Aaa;50.13\nBB;10.40"),
				sendChan: make(chan *valueMessage, 2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseChunk(tt.args.chunk, tt.args.sendChan)

			close(tt.args.sendChan)

			first := <-tt.args.sendChan
			if first.station != "Aaa" {
				t.Errorf("First station should be Aaa, but got %s", first.station)
			}
			if first.value != 50.13 {
				t.Errorf("First value should be 50.13, but got %f", first.value)
			}

			second := <-tt.args.sendChan
			if second.station != "BB" {
				t.Errorf("Second station should be BB, but got %s", second.station)
			}
			if second.value != 10.40 {
				t.Errorf("Second value should be 10.40, but got %f", second.value)
			}
		})
	}
}
