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
		_ = fib(30)
	}
}

// O(n) for iterative version
func fib(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
