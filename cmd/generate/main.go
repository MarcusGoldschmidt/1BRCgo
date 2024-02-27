package main

import (
	"bufio"
	_ "embed"
	"flag"
	"github.com/MarcusGoldschmidt/1brcgo/pkg/cities"
	"github.com/MarcusGoldschmidt/1brcgo/pkg/unit"
	"log"
	"os"
)

var BufferSize = unit.MB * 10

var size = flag.Int("size", cities.BILLION, "Number of lines to generate")
var file = flag.String("file", "./generated_file.csv", "File to write to")

func main() {
	flag.Parse()

	availableCities := cities.GetCities()

	fileP, err := os.OpenFile(*file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	write := bufio.NewWriterSize(fileP, int(BufferSize))

	err = cities.GenerateFromCities(availableCities, write, *size)
	if err != nil {
		log.Fatal(err)
	}
}
