package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
)

const PARA_CFG string = "para.json"

var ErrParsingFailed error = errors.New("Parsing json failed")

type JsonSection struct {
	Command []string `json:"command"`
	Count   uint64   `json:"count"`
}

type JsonConfig struct {
	Session []JsonSection `json:"session"`
}

func (config *JsonConfig) ForEach(fn func(command *JsonSection)) {
	for _, section := range config.Session {
		fn(&section)
	}
}

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

func OpenConfig(locate string) (*os.File, error) {
	currentPath := path.Join(locate)
	if config, err := os.OpenFile(currentPath, os.O_RDONLY, 0444); err != nil {
		return nil, err
	} else {
		return config, nil
	}
}

func OpenDefault() (*os.File, error) {
	return OpenConfig(PARA_CFG)
}
