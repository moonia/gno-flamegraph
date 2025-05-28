package main

import (
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	start := time.Now()
	for time.Since(start) < 10*time.Second {
		for i := 1; i < 10000; i++ {
			_ = power(i % 20, 10)
		}
	}
}

// O(exp) & O(1)
func power(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
