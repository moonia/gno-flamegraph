package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/pprof"
)

func Record(f func(), out string) {
    prof, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(prof)
    f()
    pprof.StopCPUProfile()
    prof.Close()

    err := exec.Command("go", "tool", "pprof", "-svg", os.Args[0], "cpu.prof").Run()
    if err != nil {
        fmt.Println(err)
    }
    os.Rename("profile001.svg", out)
}