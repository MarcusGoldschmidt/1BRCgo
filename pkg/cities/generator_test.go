package cities

import (
	"log"
	"os"
	"testing"
)

func BenchmarkGenerateFromCities(b *testing.B) {
	cities := GetCities()

	for i := 0; i < b.N; i++ {
		file, err := os.CreateTemp("", "1brcgo_cities")
		if err != nil {
			log.Fatal(err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				b.Error(err)
			}
		}(file.Name())

		err = GenerateFromCities(cities, file, 10_000_000)
		if err != nil {
			log.Fatal(err)
		}
	}
}
