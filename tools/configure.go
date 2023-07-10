package tools

import (
	"bufio"
	"encoding/json"
	"os"
)

type Configure struct {
	Name     string            `json:"name"`
	Mode     string            `json:"mode"`
	Host     string            `json:"host"`
	Port     string            `json:"port"`
	Database DatabaseConfigure `json:"database"`
}

type DatabaseConfigure struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	Name string `json:"name"`
}

var _cfg *Configure

func ParseConfigure(path string) (*Configure, error) {
	file, err := os.Open((path))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&_cfg); err != nil {
		return nil, err
	}
	return _cfg, nil
}
