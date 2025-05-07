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

/**** EXPERIMENTAL ****/
type JsonSection struct {
	Command []string `json:"command"`
	Count   uint64   `json:"count"`
}

type JsonConfig struct {
	Session []JsonSection `json:"session"`
}

func (cfg *JsonConfig) ForEach(fn func(command *JsonSection)) {
	if cfg == nil || fn == nil {
		return
	}

	for _, section := range cfg.Session {
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

/**** EXPERIMENTAL ****/

type Config struct {
	CmdList [][]string `json:"commands"`
}

func (cfg *Config) Map(fn func(command []string)) {
	if cfg == nil || fn == nil {
		return
	}

	for _, command := range cfg.CmdList {
		fn(command)
	}
}

func OpenConfig() (*os.File, error) {
	currentPath := path.Join(PARA_CFG)
	if config, err := os.OpenFile(currentPath, os.O_RDONLY, 0444); err != nil {
		return nil, err
	} else {
		return config, nil
	}
}

func UnpackJson(file io.Reader) (*Config, error) {
	config := new(Config)
	content, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}
