package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"time"
)

// go build test.go
// go run main.go ./test
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <executable_path>")
	}
	executablePath := os.Args[1]
	duration := 30 * time.Second

	err := profile(executablePath, duration)
	if err != nil {
		log.Fatal("Profiling failed:", err)
	}
}

func profile(executablePath string, duration time.Duration) error {
	log.Println("Starting profiling of:", executablePath)

	cmd := exec.Command(executablePath)
	cmd.Env = append(os.Environ(), "PPROF_PROFILE=1")
	cmd.Start()

	log.Printf("Started process: %s", executablePath)
	time.Sleep(duration)

	if err := cmd.Process.Kill(); err != nil {
		return errors.New("failed to kill process: " + err.Error())
	}

	profilePath := "cpu.pprof"
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return errors.New("profile file not found: make sure the executable generates cpu.pprof")
	}

	return generateFlamegraph(profilePath, "flamegraph.svg")
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
