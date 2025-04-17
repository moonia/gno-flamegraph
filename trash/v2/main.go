package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func profiling() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	
	go func() {
		time.Sleep(10 * time.Second)
		pprof.StopCPUProfile()
		log.Println("CPU profiling stopped")
	}()
}

func main() {
	go profiling()

	runApp()
}

func runApp() {
	for {
		_ = fib(30)
		time.Sleep(100 * time.Millisecond)
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
