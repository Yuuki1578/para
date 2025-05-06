package main

import (
	"log"
	"os/exec"

	"github.com/Yuuki1578/para/lib/cmd"
)

func main() {
	json, err := cmd.OpenConfig()

	if err != nil {
		log.Fatalln(err.Error())
	}

	config, err := cmd.UnpackJson(json)

	if err != nil {
		log.Fatalln(err.Error())
	}

	commands := cmd.New()
	config.Map(func(command []string) {
		if command == nil {
			return
		}

		switch len(command) {
		case 0:
			return

		case 1:
			commands.Append(command[0])

		default:
			commands.Append(command[0], command[1:]...)
		}
	})

	commands.Run(func(command *exec.Cmd) {
		if out, err := command.CombinedOutput(); err == nil {
			log.Println(string(out))
		}
	})
}
