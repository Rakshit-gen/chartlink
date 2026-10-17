package main

import (
	"os/exec"
	"runtime"
	"strings"
)

// cpuName is the processor model, for the benchmark table.
func cpuName() string {
	if runtime.GOOS == "darwin" {
		if out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
			return strings.TrimSpace(string(out))
		}
	}
	return runtime.GOOS + "/" + runtime.GOARCH
}
