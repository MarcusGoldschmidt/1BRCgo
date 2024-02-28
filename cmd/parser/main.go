package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/MarcusGoldschmidt/1brcgo/pkg"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strings"
	"time"
)

func main() {
	f, err := os.Create("./trace.out")
	if err != nil {
		log.Fatal("could not create trace execution profile: ", err)
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()

	f, err = os.Create("./cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// get args
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Please provide a file name")
	}

	file, err := os.OpenFile(args[1], os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	start := time.Now()

	result, err := pkg.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Station < result[j].Station
	})

	var stringsBuilder strings.Builder
	for _, i := range result {
		stringsBuilder.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", i.Station, i.Min, i.Mean, i.Max))
	}
	fmt.Println(stringsBuilder.String()[:stringsBuilder.Len()-2])

	fmt.Printf("elapsed: %s\n", time.Since(start).String())

	f, err = os.Create("./memory.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
