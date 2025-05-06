package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/Yuuki1578/para/lib/cmd"
)

func main() {
	err := cmd.New().
		Append("find", ".").
		Append("ls", "-A").
		Append("echo", "Hi!").
		Append("curl", "https://www.google.com").
		Run(func(command *exec.Cmd) {
			output, _ := command.Output()
			os.Stdout.Write(output)
		})

	if err != nil {
		log.Fatalln(err.Error())
	}
}
