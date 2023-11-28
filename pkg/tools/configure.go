package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Configure struct {
	Name     string            `json:"name"`
	Mode     string            `json:"mode"`
	Host     string            `json:"host"`
	Port     string            `json:"port"`
	Database DatabaseConfigure `json:"database"`
	Redis    RedisConfigure    `json:"redis"`
}

type DatabaseConfigure struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	Name string `json:"name"`
}

type RedisConfigure struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Pwd  string `json:"pwd"`
}

var _cfg *Configure

func ParseConfigure() (*Configure, error) {
	env := os.Getenv("CONFIG_PATH")
	fmt.Println("CONFIG_PATH", env)
	var path string
	if customPath := os.Getenv("CONFIG_PATH"); customPath != "" {
		path = customPath
	}
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
