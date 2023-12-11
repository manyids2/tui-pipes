package components

import (
	"encoding/json"
	"os"
)

type Command struct {
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
}

// Input config
type Config struct {
	Path     string
	List     Command            `json:"list,omitempty"`
	Tree     Command            `json:"tree,omitempty"`
	Preview  Command            `json:"preview"`
	Bindings map[string]Command `json:"bindings,omitempty"`
}

func ReadConfig(path string) (Config, error) {
	config := Config{}
	config.Path = path

	file, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
