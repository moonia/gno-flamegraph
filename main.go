package main

import (
	"bytes"
	"errors"
	// "io"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"
)

func main() {
	err := profile(runApp, 100 * time.Second)
	if err != nil {
		log.Fatal("Profiling failed:", err)
	}
}

func profile(fn func(), duration time.Duration) error {
	log.Println("Starting CPU profiling...")

	f, err := os.Create("cpu.pprof")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		return err
	}
	fnDone := make(chan struct{})
	
	go func() {
		fn()
		close(fnDone)
	}()

	select {
	case <-time.After(duration):
		log.Println("Profiling duration reached.")
	case <-fnDone:
		log.Println("Function completed before timeout.")
	}

	pprof.StopCPUProfile()
	log.Println("CPU profiling stopped.")
	return generateFlamegraph("cpu.pprof", "flamegraph.svg")
}

func generateFlamegraph(profilePath, outputPath string) error {
	log.Println("Generating flamegraph...")
	cmdA := exec.Command("go", "tool", "pprof", "-raw", profilePath) // go tool pprof -raw cpu.pprof
	cmdB := exec.Command("/Users/moonia/FlameGraph/stackcollapse-go.pl") // stackcollapse-go.pl
	cmdC := exec.Command("/Users/moonia/FlameGraph/flamegraph.pl") // flamegraph.pl
	aOut, err := cmdA.StdoutPipe()

	if err != nil {
		return errors.New("failed to get stdout from pprof: " + err.Error())
	}
	cmdB.Stdin = aOut

	bOut, err := cmdB.StdoutPipe()
	if err != nil {
		return errors.New("failed to get stdout from stackcollapse: " + err.Error())
	}
	cmdC.Stdin = bOut

	var flamegraph bytes.Buffer
	cmdC.Stdout = &flamegraph

	if err := cmdC.Start(); err != nil {
		return errors.New("failed to start flamegraph.pl: " + err.Error())
	}
	if err := cmdB.Start(); err != nil {
		return errors.New("failed to start stackcollapse-go.pl: " + err.Error())
	}
	if err := cmdA.Run(); err != nil {
		return errors.New("failed to run go tool pprof: " + err.Error())
	}
	if err := cmdB.Wait(); err != nil {
		return errors.New("stackcollapse-go.pl failed: " + err.Error())
	}
	if err := cmdC.Wait(); err != nil {
		return errors.New("flamegraph.pl failed: " + err.Error())
	}
	if err := os.WriteFile(outputPath, flamegraph.Bytes(), 0644); err != nil {
		return errors.New("failed to write flamegraph.svg: " + err.Error())
	}
	log.Println("Flamegraph generated at:", outputPath)
	return nil
}

func runApp() {
	start := time.Now()
	for time.Since(start) < 10*time.Second {
		_ = fib(10)
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
