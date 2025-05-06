package cmd

import (
	"encoding/json"
)

type ParaCommands struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type ParaConfig struct {
	Commands []ParaCommands `json:"commands"`
}
