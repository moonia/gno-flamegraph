package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"
)

// Function to continuously profile and generate flamegraph while the program runs
func profiling() {
	// Create CPU profile file
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	// Prepare for streaming to the stackcollapse and flamegraph
	cmd1 := exec.Command("go", "tool", "pprof", "-raw", "./cpu.pprof")
	cmd2 := exec.Command("/home/pol/downloads/FlameGraph/stackcollapse-go.pl")
	cmd3 := exec.Command("/home/pol/downloads/FlameGraph/flamegraph.pl")

	// Pipe the output through stackcollapse and flamegraph
	pipe1, err := cmd1.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating pipe for pprof:", err)
	}
	fmt.Println(cmd1.Stderr)
	cmd2.Stdin = pipe1
	pipe2, err := cmd2.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating pipe for stackcollapse:", err)
	}
	fmt.Println(cmd2.Stderr)
	cmd3.Stdin = pipe2

	// Redirect final output to flamegraph.svg
	outfile, err := os.Create("flamegraph.svg")
	if err != nil {
		log.Fatal("Error creating flamegraph.svg:", err)
	}
	defer outfile.Close()
	cmd3.Stdout = outfile

	// Start the commands
	if err := cmd1.Start(); err != nil {
		log.Fatal("Error starting pprof:", err)
	}
	if err := cmd2.Start(); err != nil {
		log.Fatal("Error starting stackcollapse:", err)
	}
	if err := cmd3.Start(); err != nil {
		log.Fatal("Error starting flamegraph:", err)
	}

	// Run for a period and continuously generate the flamegraph
	time.Sleep(10 * time.Second)

	// Stop profiling after the program finishes or a certain condition
	pprof.StopCPUProfile()
	log.Println("CPU profiling stopped")
}

func main() {
	// Start profiling and generate the flamegraph while the program runs
	profiling()

	// Simulate your app's workload (like 1 + 1 in your example)
	for i := 0; i < 100; i++ {
		_ = fib(30) // Simulate some work
		time.Sleep(100 * time.Millisecond)
	}
}

// Simple Fibonacci function for simulation
func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
