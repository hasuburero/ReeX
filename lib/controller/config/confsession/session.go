package confsession

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Node  []Node  `json:"node"`
	Group []Group `json:"group"`
}

type Group struct {
	Name     string   `json:"name"`
	NodeName []string `json:"nodename"`
}

type Node struct {
	NodeName string `json:"nodename"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Group    string `json:"group"`
}

func Read(filename string) (Config, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer fd.Close()

	buf, err := io.ReadAll(fd)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(buf, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
