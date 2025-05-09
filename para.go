package main

import (
	"github.com/Yuuki1578/para/lib/cmd"
	"github.com/Yuuki1578/para/lib/logger"
	"os/exec"
)

func main() {
	cmd.New().
		Append(10, "echo", "hello").
		Run(func(command *exec.Cmd) {
			output, _ := command.Output()
			logger.Print(string(output))
		})
}
