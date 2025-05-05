package main

import (
	// "github.com/Yuuki1578/para/lib/cmd"
	"os"
	"os/exec"
)

func main() {
	for range 10 {
		out, _ := exec.Command("ls", "-A").Output()
		os.Stdout.Write(out)
	}
}
