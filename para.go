package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/Yuuki1578/para/lib/cmd"
)

func main() {
	reader := strings.NewReader(`
		{
			"session": [
				{
					"command": ["curl", "https://www.google.com"],
					"count": 10
				}
			]
		}
	`)

	conf, err := cmd.Parse(reader)
	if err != nil {
		log.Fatalln(err.Error())
	}

	group := cmd.Init()
	conf.ForEach(func(command *cmd.JsonSection) {
		count, args := command.Count, command.Command

		switch len(args) {
		case 0:
			log.Println("Empty command")

		case 1:
			group.Append(count, args[0])

		default:
			group.Append(count, args[0], args...)
		}
	})

	group.Run(func(command *exec.Cmd) {
		out, _ := command.CombinedOutput()
		log.Println(string(out))
	})
}
