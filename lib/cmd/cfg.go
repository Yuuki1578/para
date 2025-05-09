package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
)

// Default file config name
const PARA_CFG string = "para.json"

// Indicating error while parsing config file
var ErrParsingFailed error = errors.New("Parsing json failed")

// Section for configuring command to run and how many times it run
//
//	{
//	  "command": ["man", "gcc"],
//	  "count": 1
//	}
type JsonSection struct {
	Command []string `json:"command"`
	Count   uint64   `json:"count"`
}

// Main section for config file
//
//	{
//	  "session": [
//	    {
//	      "command": [],
//	      "count": 0
//	    }
//	  ]
//	}
type JsonConfig struct {
	Session []JsonSection `json:"session"`
}

// Iterate over a slice of JsonSection and run a closure
func (config *JsonConfig) ForEach(fn func(command *JsonSection)) {
	for _, section := range config.Session {
		fn(&section)
	}
}

// Parsing config file into struct *JsonConfig
func Parse(reader io.Reader) (*JsonConfig, error) {
	cfg := new(JsonConfig)
	switch content, err := io.ReadAll(reader); err {
	case nil:
		if err := json.Unmarshal(content, cfg); err == nil {
			return cfg, nil
		}
	}

	return nil, ErrParsingFailed
}

// Open a config file on read-only mode (0444)
func OpenConfig(locate string) (*os.File, error) {
	currentPath := path.Join(locate)
	if config, err := os.OpenFile(currentPath, os.O_RDONLY, 0444); err != nil {
		return nil, err
	} else {
		return config, nil
	}
}

// Open the default config file path with read-only mode
func OpenDefault() (*os.File, error) {
	return OpenConfig(PARA_CFG)
}
