package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"
)

func generateFlamegraph() {
	var out bytes.Buffer
	log.Println("Generating flamegraph...")
	a := exec.Command("go", "tool", "pprof", "-raw", "./cpu.pprof")
	b := exec.Command("/home/pol/downloads/FlameGraph/stackcollapse-go.pl")
	c := exec.Command("/home/pol/downloads/FlameGraph/flamegraph.pl")
	aOut, err := a.StdoutPipe()
	if err != nil {
		log.Fatal("can't pipe a:", err)
	}
	b.Stdin = aOut
	bOut, err := b.StdoutPipe()
	if err != nil {
		log.Fatal("can't pipe b:", err)
	}
	c.Stdin = bOut
	c.Stdout = &out
	if err := c.Start(); err != nil {
		log.Fatal("can't start c: ", err)
	}
	if err := b.Start(); err != nil {
		log.Fatal("can't start b: ", err)
	}
	if err := a.Run(); err != nil {
		log.Fatal("can't run a: ", err)
	}
	if err := b.Wait(); err != nil {
		log.Fatal("b failed: ", err)
	}
	if err := c.Wait(); err != nil {
		log.Fatal("c failed: ", err)
	}
	if err := os.WriteFile("flamegraph.svg", out.Bytes(), 0644); err != nil {
		log.Fatal("could not create flamegraph:", err)
	}
	log.Println("Flamegraph generated")
}

func profile(fn func(), t time.Duration) {
	log.Println("Starting CPU Profiling...")
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	profilingDone := make(chan struct{})
	fnDone := make(chan struct{})
	go func() {
		select {
		case <-time.After(t):
			break
		case <-fnDone:
			time.Sleep(time.Second)
			break
		}
		pprof.StopCPUProfile()
		log.Println("CPU profiling stopped")
		generateFlamegraph()
		close(profilingDone)
	}()

	fn()
	close(fnDone)
	<-profilingDone
}

func main() {
	profile(runApp, 10*time.Second)
}

func runApp() {
	fmt.Println(fib(30))
	// for i := 0; i < 30; i++ {
	// 	fmt.Println(fib(i))
	// 	time.Sleep(100 * time.Millisecond)
	// }
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
